package logger

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"
)

var TempDir string = ""

func init() {
	TempDir = filepath.Join(os.TempDir(), "LoggerOption-Non-existent-dir")
}

func TestLogger(t *testing.T) {
	var err error
	testFile := []string{
		"",
		filepath.Join(TempDir, "test", "test.log"),
		os.TempDir(),
	}

	if runtime.GOOS != "windows" {
		noAccess := "/etc/passwd"
		_, etcerr := os.Stat(noAccess)
		if etcerr == nil && syscallAccess(noAccess, syscall.O_RDWR) != nil {
			testFile = append(testFile, noAccess)
		}
	}

	for _, v := range testFile {
		_, err = Init(&LoggerOption{
			LogFilePath: v,
		})
		assert.NotEqual(t, nil, err)
	}

	fname := filepath.Join(os.TempDir(), "test.log")
	lg, err := Init(&LoggerOption{
		LogFilePath: fname,
		CatTime:     CAT_DAY,
		FileExpire:  -100,
	})
	require.Equal(t, nil, err)
	assert.Equal(t, 0, lg.jack.MaxAge)

	lg, err = Init(&LoggerOption{
		LogFilePath: fname,
		CatTime:     CAT_DAY,
		FileExpire:  90,
	})
	require.Equal(t, nil, err)

	assert.NotEqual(t, nil, SugaredLogger)
	assert.NotEqual(t, nil, lg.jack)
	assert.Equal(t, fname, lg.jack.Filename)
	assert.Equal(t, 90, lg.jack.MaxAge)
}

func rmTempDir() error {
	_, err := os.Stat(TempDir)
	if err == nil {
		return os.RemoveAll(TempDir)
	}
	return nil
}

func createTempDir() error {
	if err := rmTempDir(); err != nil {
		return err
	}
	return os.MkdirAll(TempDir, os.ModePerm)
}
