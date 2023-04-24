package gopop

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPDFInfoParser_Parse(t *testing.T) {
	in := `Producer:        PDFKit.NET 23.1.101.39834 DMV9
CreationDate:    Wed Apr 19 12:58:37 2023 CEST
ModDate:         Wed Apr 19 12:58:37 2023 CEST
Custom Metadata: no
Metadata Stream: yes
Tagged:          no
UserProperties:  no
Suspects:        no
Form:            AcroForm
JavaScript:      no
Pages:           2
Encrypted:       no
Page    1 size:  595.32 x 841.92 pts (A4)
Page    1 rot:   0
Page    2 size:  595.32 x 841.92 pts (A4)
Page    2 rot:   0
File size:       230296 bytes
Optimized:       no
PDF version:     1.7
`
	parser := NewPDFInfoParser(bytes.NewBufferString(in))
	info, err := parser.Parse()
	require.NoError(t, err)
	require.NotNil(t, info)

	assert.Len(t, info.Pages, 2, "wrong amount of pages")
	assert.Equal(t, "AcroForm", info.Properties["Form"])
	assert.Equal(t, "2", info.Properties["Pages"])
	assert.Equal(t, "230296 bytes", info.Properties["File size"])

	t.Run("json representation", func(t *testing.T) {

		j, err := json.Marshal(info)
		require.NoError(t, err)

		assert.JSONEq(t, `{
			"Properties": {
				"CreationDate": "Wed Apr 19 12:58:37 2023 CEST",
				"Custom Metadata": "no",
				"Encrypted": "no",
				"File size": "230296 bytes",
				"Form": "AcroForm",
				"JavaScript": "no",
				"Metadata Stream": "yes",
				"ModDate": "Wed Apr 19 12:58:37 2023 CEST",
				"Optimized": "no",
				"PDF version": "1.7",
				"Pages": "2",	
				"Producer": "PDFKit.NET 23.1.101.39834 DMV9",
				"Suspects": "no",
				"Tagged": "no",
				"UserProperties": "no"
			},
			"Pages": {
				"1": {
					"rot": "0",
					"size": "595.32 x 841.92 pts (A4)"
				},
				"2": {
					"rot": "0",
					"size": "595.32 x 841.92 pts (A4)"
				}
			}
		}`, string(j))
	})

	t.Run("info methods", func(t *testing.T) {
		assert.Equal(t, "PDFKit.NET 23.1.101.39834 DMV9", info.Producer())
	})

	t.Run("page methods", func(t *testing.T) {
		page1 := info.Page(1)
		require.NotNil(t, page1)

		t.Run("raw values", func(t *testing.T) {
			assert.Equal(t, "595.32 x 841.92 pts (A4)", page1.SizeRaw())
			assert.Equal(t, "0", page1.RotRaw())
		})

		t.Run("parsed values", func(t *testing.T) {
			size := page1.Size()
			assert.Equal(t, 595.32, size.Width)
			assert.Equal(t, 841.92, size.Height)
			assert.Equal(t, "pts", size.Unit)
			assert.Equal(t, "(A4)", size.Note)
		})
	})
}
