package judge

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type ExecutionResult struct {
	Output      string
	ExecutionMs int64
}

func RunCode(language string, code string) (ExecutionResult, error) {
	var (
		sourceFile string
		image      string
		cmd        []string
	)

	switch language {
	case "python":
		sourceFile = "submission.py"
		image = "python:3.12-alpine"
		cmd = []string{"python", "/app/" + sourceFile}
	case "cpp":
		sourceFile = "submission.cpp"
		image = "gcc:13.2.0"
		cmd = []string{"/bin/sh", "-c", "g++ /app/" + sourceFile + " -o /app/a.out && /app/a.out"}
	case "c":
		sourceFile = "submission.c"
		image = "gcc:13.2.0"
		cmd = []string{"/bin/sh", "-c", "gcc /app/" + sourceFile + " -o /app/a.out && /app/a.out"}
	default:
		return ExecutionResult{}, fmt.Errorf("unsupported language: %s", language)
	}

	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return ExecutionResult{}, err
	}

	tmpFile, err := os.CreateTemp("", "submission-*."+filepath.Ext(sourceFile)[1:])
	if err != nil {
		return ExecutionResult{}, err
	}

	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(code); err != nil {
		tmpFile.Close()
		return ExecutionResult{}, err
	}

	if err := tmpFile.Close(); err != nil {
		return ExecutionResult{}, err
	}

	res, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   cmd,
		Tty:   false,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: tmpFile.Name(),
				Target: "/app/" + sourceFile,
			},
		},
		Resources: container.Resources{
			Memory:   256 * 1024 * 1024, // 256 MB
			NanoCPUs: 500000000,         // 0.5 CPU
		},
	}, nil, nil, "")
	if err != nil {
		return ExecutionResult{}, err
	}

	defer func() {
		_ = cli.ContainerRemove(ctx, res.ID, container.RemoveOptions{Force: true})
	}()

	// start timer and run the container
	start := time.Now()
	if err := cli.ContainerStart(ctx, res.ID, container.StartOptions{}); err != nil {
		return ExecutionResult{}, err
	}

	timeout := time.After(10 * time.Second)
	done := make(chan struct{})
	var logs bytes.Buffer

	go func() {
		out, _ := cli.ContainerLogs(ctx, res.ID, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
		})
		defer out.Close()

		_, _ = io.Copy(&logs, out)

		close(done)
	}()

	select {
	case <-done:
	case <-timeout:
		cli.ContainerKill(ctx, res.ID, "SIGKILL")
		return ExecutionResult{}, fmt.Errorf("execution timeout")
	}

	statusCh, errCh := cli.ContainerWait(ctx, res.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return ExecutionResult{}, err
		}
	case <-statusCh:
	}

	elapsedMs := time.Since(start).Milliseconds()

	return ExecutionResult{
		Output:      logs.String(),
		ExecutionMs: elapsedMs,
	}, nil

}
