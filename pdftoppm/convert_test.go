package pdftoppm_test

import (
	"context"
	"github.com/fabiante/gopop/pdftoppm"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestCommand_Convert(t *testing.T) {
	t.Run("using minimal parameters", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "gopop-pdftoppm-test-*")
		require.NoError(t, err)

		cmd, err := pdftoppm.NewCommand("test.pdf", dir+"/img")
		require.NoError(t, err)
		require.NotNil(t, cmd)

		err = cmd.Run(context.Background())
		require.NoError(t, err)

		files, err := os.ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, 2, len(files), "expected 2 output files, one for each page of the PDF")
	})
}
