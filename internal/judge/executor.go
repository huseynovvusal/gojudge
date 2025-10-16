package judge

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
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
	tmpFile.WriteString(code)
	tmpFile.Close()

	res, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "python:3.12-alpine",
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
			Memory:   128 * 1024 * 1024,
			NanoCPUs: 500000000, // 0.5 CPU
		},
	}, nil, nil, "")
	if err != nil {
		return "", err
	}

	defer func() {
		_ = cli.ContainerRemove(ctx, res.ID, container.RemoveOptions{Force: true})
	}()

	if err := cli.ContainerStart(ctx, res.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	timeout := time.After(3 * time.Second)
	done := make(chan struct{})
	var logs bytes.Buffer

	go func() {
		out, _ := cli.ContainerLogs(ctx, res.ID, types.ContainerLogsOptions{
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

	return logs.String(), nil

}
