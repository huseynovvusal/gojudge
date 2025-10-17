package judge

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func RunCode(language string, code string) (string, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	tmpFile, err := os.CreateTemp("", "submission-*.py")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(code); err != nil {
		tmpFile.Close()
		return "", err
	}
	if err := tmpFile.Close(); err != nil {
		return "", err
	}

	res, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "python:3.12-slim",
		Cmd:   []string{"python", "/app/submission.py"},
		Tty:   false,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: tmpFile.Name(),
				Target: "/app/submission.py",
			},
		},
		Resources: container.Resources{
			Memory:   512 * 1024 * 1024, // 512 MB
			NanoCPUs: 50000000,          // 0.5 CPU
		},
	}, nil, nil, "")
	if err != nil {
		return "", err
	}

	defer func() {
		_ = cli.ContainerRemove(ctx, res.ID, container.RemoveOptions{Force: true})
	}()

	// start timer and run the container
	start := time.Now()
	if err := cli.ContainerStart(ctx, res.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	timeout := time.After(1 * time.Second)
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
		return logs.String(), fmt.Errorf("execution timeout")
	}

	statusCh, errCh := cli.ContainerWait(ctx, res.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return logs.String(), err
		}
	case <-statusCh:
	}

	elapsedMs := time.Since(start).Milliseconds()
	fmt.Printf("execution took %d ms\n", elapsedMs)

	return logs.String() + fmt.Sprintf("\nExecution time: %d ms", elapsedMs), nil

}
