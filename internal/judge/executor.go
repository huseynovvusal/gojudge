package judge

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type ExecutionResult struct {
	Output      string
	ExecutionMs int64
}

func RunCode(language string, code string, timeLimit int16, memoryLimit int16, cpuLimit int16) (ExecutionResult, error) {
	switch language {
	case "python":
		return RunPythonWithNsjail(code, timeLimit, memoryLimit, cpuLimit)
	case "c":
		return RunCWithNsjail(code, timeLimit, memoryLimit, cpuLimit)
	case "cpp":
		return RunCppWithNsjail(code, timeLimit, memoryLimit, cpuLimit)
	default:
		return ExecutionResult{}, fmt.Errorf("unsupported language: %s", language)
	}
}

// RunPythonWithNsjail runs Python code using nsjail with specified resource limits.
// timeLimit, memoryLimit, and cpuLimit are in seconds and megabytes.
// Returns the output and execution time in milliseconds.
func RunPythonWithNsjail(code string, timeLimit int16, memoryLimit int16, cpuLimit int16) (ExecutionResult, error) {
	tmpFile, _ := os.CreateTemp("", "submission-*.py")
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString(code)
	tmpFile.Close()

	memStr := fmt.Sprintf("%d", memoryLimit)
	cpuStr := fmt.Sprintf("%d", cpuLimit)
	timeStr := fmt.Sprintf("%d", timeLimit)

	start := time.Now()
	cmd := exec.Command(
		"nsjail",
		"--bindmount_ro", "/usr/bin/python3",
		"--bindmount_ro", "/usr/lib",
		"--bindmount_ro", "/lib",
		"--bindmount_ro", "/lib64",
		"--bindmount_ro", tmpFile.Name(),
		"--rlimit_as", memStr,
		"--rlimit_cpu", cpuStr,
		"--time_limit", timeStr,
		"--log", "error", // Log only errors
		"--",
		"/usr/bin/python3", tmpFile.Name(),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return ExecutionResult{}, err
	}

	elapsedMs := time.Since(start).Milliseconds()

	return ExecutionResult{
		Output:      string(output),
		ExecutionMs: elapsedMs,
	}, nil

}

func RunCWithNsjail(code string, timeLimit int16, memoryLimit int16, cpuLimit int16) (ExecutionResult, error) {
	srcFile, _ := os.CreateTemp("", "submission-*.c")
	defer os.Remove(srcFile.Name())
	srcFile.WriteString(code)
	srcFile.Close()

	binPath := srcFile.Name() + ".out"
	exec.Command("gcc", srcFile.Name(), "-o", binPath).Run()
	defer os.Remove(binPath)

	memStr := fmt.Sprintf("%d", memoryLimit)
	cpuStr := fmt.Sprintf("%d", cpuLimit)
	timeStr := fmt.Sprintf("%d", timeLimit)

	start := time.Now()
	cmd := exec.Command(
		"nsjail",
		"--bindmount_ro", binPath,
		"--bindmount_ro", "/usr/lib",
		"--bindmount_ro", "/lib",
		"--bindmount_ro", "/lib64",
		"--rlimit_as", memStr,
		"--rlimit_cpu", cpuStr,
		"--time_limit", timeStr,
		"--",
		srcFile.Name(),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return ExecutionResult{}, err
	}

	elapsedMs := time.Since(start).Milliseconds()

	return ExecutionResult{
		Output:      string(output),
		ExecutionMs: elapsedMs,
	}, nil

}

func RunCppWithNsjail(code string, timeLimit int16, memoryLimit int16, cpuLimit int16) (ExecutionResult, error) {
	srcFile, _ := os.CreateTemp("", "submission-*.cpp")
	defer os.Remove(srcFile.Name())
	srcFile.WriteString(code)
	srcFile.Close()

	binPath := srcFile.Name() + ".out"
	exec.Command("g++", srcFile.Name(), "-o", binPath).Run()
	defer os.Remove(binPath)

	memStr := fmt.Sprintf("%d", memoryLimit)
	cpuStr := fmt.Sprintf("%d", cpuLimit)
	timeStr := fmt.Sprintf("%d", timeLimit)

	start := time.Now()
	cmd := exec.Command(
		"nsjail",
		"--bindmount_ro", binPath,
		"--bindmount_ro", "/usr/lib",
		"--bindmount_ro", "/lib",
		"--bindmount_ro", "/lib64",
		"--rlimit_as", memStr,
		"--rlimit_cpu", cpuStr,
		"--time_limit", timeStr,
		"--",
		srcFile.Name(),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return ExecutionResult{}, err
	}

	elapsedMs := time.Since(start).Milliseconds()

	return ExecutionResult{
		Output:      string(output),
		ExecutionMs: elapsedMs,
	}, nil

}
