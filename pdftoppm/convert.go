package pdftoppm

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

type Command struct {
	args []string
}

// NewCommand creates a pdftoppm command with the given PDF path and options.
//
// This does not support writing the PDF or reading the output from the stdin/stdout streams.
// If you need that feature, consider adding it to this project.
func NewCommand(pdfPath, imagePathPrefix string, opts ...ConvertOption) (*Command, error) {
	cmd := &Command{
		args: make([]string, 0),
	}

	for _, opt := range opts {
		opt(cmd)
	}

	cmd.args = append(cmd.args, pdfPath)
	cmd.args = append(cmd.args, imagePathPrefix)

	return cmd, nil
}

// Run executes the pdftoppm command and waits for it to finish.
//
// If this does not return an error, images are rendered into the configured output path (see NewCommand).
func (cmd *Command) Run(ctx context.Context) error {
	c := exec.CommandContext(ctx, "pdftoppm", cmd.args...)

	var buf bytes.Buffer
	c.Stderr = &buf

	err := c.Run()
	if err != nil {
		return fmt.Errorf("pdftoppm failed with status %v and stderr: %s", c.ProcessState.ExitCode(), buf.String())
	}
	return nil
}
