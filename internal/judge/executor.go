package judge

import (
	"context"
	"os"
)

func RunCode(language string, code string) (string, error) {
	ctx := context.Background()

	tmpDir, err := os.MkdirTemp("", "submission-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	filePath := tmpDir + "/main.py"

	if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
		return "", err
	}

}
