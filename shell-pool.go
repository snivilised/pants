package pants

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// ShellSession defines the interface for a persistent shell session.
type ShellSession interface {
	Execute(ctx context.Context, command string) (string, error)
	Close() error
}

// InteractiveShellSession implements ShellSession using a persistent
// background shell process.
type InteractiveShellSession struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	reader *bufio.Reader
	marker string
}

// NewInteractiveShellSession creates a new persistent shell session.
func NewInteractiveShellSession(shellPath string) (*InteractiveShellSession, error) {
	cmd := exec.Command(shellPath, "+m", "-i")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	marker := "__PANTS_CMD_DONE__"
	return &InteractiveShellSession{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		reader: bufio.NewReader(stdout),
		marker: marker,
	}, nil
}

// Execute sends a command to the shell and waits for it to complete.
func (s *InteractiveShellSession) Execute(ctx context.Context, command string) (string, error) {
	// We append a marker to know when the command is done.
	// This is a simple implementation and might need refinement for production use.
	fullCmd := fmt.Sprintf("%s; echo %s\n", command, s.marker)
	if _, err := io.WriteString(s.stdin, fullCmd); err != nil {
		return "", err
	}

	var output strings.Builder
	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			return output.String(), err
		}
		if strings.Contains(line, s.marker) {
			break
		}
		output.WriteString(line)
	}

	return strings.TrimSpace(output.String()), nil
}

// Close terminates the shell session.
func (s *InteractiveShellSession) Close() error {
	_ = s.stdin.Close()
	return s.cmd.Wait()
}

// ShellPool is a specialized pool for shell command execution.
type ShellPool struct {
	*ManifoldStatePool[string, string, ShellSession]
}

// NewShellPool creates a new pool specialized for shell commands.
func NewShellPool(ctx context.Context,
	shellPath string,
	wg WaitGroup,
	options ...Option,
) (*ShellPool, error) {
	// Define the manifold function that uses the shell session.
	mf := func(command string, session ShellSession) (string, error) {
		return session.Execute(ctx, command)
	}

	// Add state initializer and finalizer to the options.
	initializer := func(id RoutineID) interface{} {
		session, err := NewInteractiveShellSession(shellPath)
		if err != nil {
			return nil
		}
		return session
	}

	finalizer := func(state interface{}) {
		if session, ok := state.(ShellSession); ok {
			_ = session.Close()
		}
	}

	opts := append(options,
		WithStateInitializer(initializer),
		WithStateFinalizer(finalizer),
	)

	base, err := NewManifoldStatePool(ctx, mf, wg, opts...)
	if err != nil {
		return nil, err
	}

	return &ShellPool{
		ManifoldStatePool: base,
	}, nil
}
