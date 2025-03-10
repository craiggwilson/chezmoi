package chezmoi

import (
	"io"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/twpayne/chezmoi/v2/pkg/archivetest"
)

func TestWalkArchive(t *testing.T) {
	for _, tc := range []struct {
		name          string
		dataFunc      func(map[string]interface{}) ([]byte, error)
		archiveFormat ArchiveFormat
	}{
		{
			name:          "tar",
			dataFunc:      archivetest.NewTar,
			archiveFormat: ArchiveFormatTar,
		},
		{
			name:          "zip",
			dataFunc:      archivetest.NewZip,
			archiveFormat: ArchiveFormatZip,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := map[string]interface{}{
				"dir1": map[string]interface{}{
					"subdir1": map[string]interface{}{
						"file1": "",
						"file2": "",
					},
					"subdir2": map[string]interface{}{
						"file1": "",
						"file2": "",
					},
				},
				"dir2": map[string]interface{}{
					"subdir1": map[string]interface{}{
						"file1": "",
						"file2": "",
					},
					"subdir2": map[string]interface{}{
						"file1": "",
						"file2": "",
					},
				},
				"file1":    "",
				"file2":    "",
				"symlink1": &archivetest.Symlink{Target: "file1"},
				"symlink2": &archivetest.Symlink{Target: "file2"},
			}
			data, err := tc.dataFunc(root)
			require.NoError(t, err)

			expectedNames := []string{
				"dir1",
				"dir1/subdir1",
				"dir1/subdir1/file1",
				"dir1/subdir1/file2",
				"dir1/subdir2",
				"dir2",
				"file1",
				"file2",
				"symlink1",
			}

			var actualNames []string
			walkArchiveFunc := func(name string, info fs.FileInfo, r io.Reader, linkname string) error {
				actualNames = append(actualNames, name)
				switch name {
				case "dir1/subdir2":
					return fs.SkipDir
				case "dir2":
					return fs.SkipDir
				case "symlink1":
					return Break
				default:
					return nil
				}
			}
			require.NoError(t, WalkArchive(data, tc.archiveFormat, walkArchiveFunc))
			assert.Equal(t, expectedNames, actualNames)
		})
	}
}
