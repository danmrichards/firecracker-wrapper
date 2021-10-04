package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) error {
	sf, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer sf.Close()

	si, err := sf.Stat()
	if err != nil {
		return fmt.Errorf("state src file: %w", err)
	} else if !si.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	df, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create destination file: %w", err)
	}
	defer df.Close()

	if _, err = io.Copy(df, sf); err != nil {
		return fmt.Errorf("copy file contents: %w", err)
	}
	if err = df.Sync(); err != nil {
		return fmt.Errorf("sync destination file: %w", err)
	}

	// Preserve original file permissions.
	return os.Chmod(dst, si.Mode())
}

// CopyDir copies the contents of the source directory to the destination.
func CopyDir(src, dst string) error {
	if src == dst {
		return fmt.Errorf("source and destination directory cannot be the same")
	}

	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	file, err := f.Stat()
	if err != nil {
		return err
	}
	if !file.IsDir() {
		return fmt.Errorf("Source " + file.Name() + " is not a directory!")
	}

	if err = os.Mkdir(dst, 0755); err != nil {
		return err
	}

	dirContents, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, df := range dirContents {
		if df.IsDir() {
			if err = CopyDir(
				filepath.Join(src, df.Name()), filepath.Join(dst, df.Name()),
			); err != nil {
				return err
			}
		} else {
			if err = CopyFile(
				filepath.Join(src, df.Name()), filepath.Join(dst, df.Name()),
			); err != nil {
				return err
			}
		}
	}

	return nil
}

// DirectoryExists returns true if path exists and is a directory.
func DirectoryExists(path string) bool {
	if path == "" {
		return false
	}
	if info, err := os.Stat(path); err == nil {
		return info.IsDir()
	}
	return false
}

// ExecCommand executes the given command and returns the combined output.
func ExecCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command %q exited with %q: %v", cmd.Args, out, err)
	}

	return string(bytes.TrimSpace(out)), nil
}
