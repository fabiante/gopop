package pdftoppm_test

import (
	"context"
	"github.com/fabiante/gopop/pdftoppm"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
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

	t.Run("limit to first page", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "gopop-pdftoppm-test-*")
		require.NoError(t, err)

		cmd, err := pdftoppm.NewCommand(
			"test.pdf", dir+"/img",
			pdftoppm.Last(1),
		)
		require.NoError(t, err)
		require.NotNil(t, cmd)

		err = cmd.Run(context.Background())
		require.NoError(t, err)

		files, err := os.ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, 1, len(files), "expected 1 output file")
	})

	t.Run("can define png format", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "gopop-pdftoppm-test-*")
		require.NoError(t, err)

		cmd, err := pdftoppm.NewCommand(
			"test.pdf", dir+"/img",
			pdftoppm.PNG(),
		)
		require.NoError(t, err)
		require.NotNil(t, cmd)

		err = cmd.Run(context.Background())
		require.NoError(t, err)

		files, err := os.ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, 2, len(files), "expected 1 output file")

		for _, file := range files {
			require.True(t, strings.HasSuffix(file.Name(), ".png"), "expected output file %q to be in PNG format", file.Name())

			// Get file info
			info, err := os.Stat(dir + "/" + file.Name())
			require.NoError(t, err)
			size := info.Size()

			// We know that without this limit, pages will be larger that 25K
			// This is an assertion ment to ensure the test using scale works
			require.GreaterOrEqual(t, size, int64(25000), "expected output file %q to be small (scale was set to 100)", file.Name())
		}
	})

	t.Run("can define max width/height using scale", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "gopop-pdftoppm-test-*")
		require.NoError(t, err)

		cmd, err := pdftoppm.NewCommand(
			"test.pdf", dir+"/img",
			pdftoppm.ScaleTo(100), // low scale will result in small image size
		)
		require.NoError(t, err)
		require.NotNil(t, cmd)

		err = cmd.Run(context.Background())
		require.NoError(t, err)

		files, err := os.ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, 2, len(files), "expected 1 output file")

		for _, file := range files {
			// Get file info
			info, err := os.Stat(dir + "/" + file.Name())
			require.NoError(t, err)
			size := info.Size()
			// We know that without this limit, pages will be larger that 25K
			require.Less(t, size, int64(25000), "expected output file %q to be small (scale was set to 100)", file.Name())
		}
	})

	t.Run("can define resolution", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "gopop-pdftoppm-test-*")
		require.NoError(t, err)

		cmd, err := pdftoppm.NewCommand(
			"test.pdf", dir+"/img",
			pdftoppm.Resolution(300), // high resolution will result in larger image size
		)
		require.NoError(t, err)
		require.NotNil(t, cmd)

		err = cmd.Run(context.Background())
		require.NoError(t, err)

		files, err := os.ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, 2, len(files), "expected 1 output file")

		for _, file := range files {
			// Get file info
			info, err := os.Stat(dir + "/" + file.Name())
			require.NoError(t, err)
			size := info.Size()
			// We know that with the higher DPI, images will be larger that 20M
			require.Greater(t, size, int64(25000000), "expected output file %q to be large (due to high dpi)", file.Name())
		}
	})
}
