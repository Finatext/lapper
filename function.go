package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

const (
	initialBufSize = 65536
	maxBufSize     = 6553600
)

type Function struct {
	Command string
	Args    []string
	Payload []byte
	Stdout  io.Writer
	Stderr  io.Writer
}

func NewFunction(command string, args []string, payload []byte) *Function {
	f := &Function{
		Command: command,
		Args:    args,
		Payload: payload,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}

	return f
}

func (fn *Function) SetStdout(w io.Writer) {
	fn.Stdout = w
}

func (fn *Function) SetStderr(w io.Writer) {
	fn.Stderr = w
}

func (fn *Function) Run() (string, string, error) {
	cmd := exec.Command(fn.Command, fn.Args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", fmt.Errorf("Failed to open stdout pipe: %s", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", "", fmt.Errorf("Failed to open stderr pipe: %s", err)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", "", fmt.Errorf("Failed to open stdin pipe: %s", err)
	}

	var stdout1, stderr1 bytes.Buffer
	stdout2 := io.TeeReader(stdout, &stdout1)
	stderr2 := io.TeeReader(stderr, &stderr1)

	err = cmd.Start()
	if err != nil {
		return "", "", fmt.Errorf("Failed to start command: %s", err)
	}

	if _, err := stdin.Write(fn.Payload); err != nil {
		return "", "", fmt.Errorf("failed to write to stdin: %w", err)
	}
	stdin.Close()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		scanner := bufio.NewScanner(stdout2)
		buf := make([]byte, initialBufSize)
		scanner.Buffer(buf, maxBufSize)

		for scanner.Scan() {
			fmt.Fprintln(fn.Stdout, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(fn.Stderr, "Scan error (stdout):", err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		scanner := bufio.NewScanner(stderr2)
		buf := make([]byte, initialBufSize)
		scanner.Buffer(buf, maxBufSize)

		for scanner.Scan() {
			fmt.Fprintln(fn.Stderr, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(fn.Stderr, "Scan error (stderr):", err)
		}
		wg.Done()
	}()

	wg.Wait()
	err = cmd.Wait()
	// Before checking command success, we need to collect all the output
	stdoutBytes := stdout1.String()
	stderrBytes := stderr1.String()

	if err != nil {
		var exitError *exec.ExitError
		if ok := errors.As(err, &exitError); ok {
			return stdoutBytes, stderrBytes, fmt.Errorf("command failed with exit code %d", exitError.ExitCode())
		} else {
			// If this is not an ExitError, it's a unexpected situation
			return stderrBytes, stderrBytes, fmt.Errorf("failed to execute command: %w", err)
		}
	}

	return stdoutBytes, stderrBytes, nil
}
