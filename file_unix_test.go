// +build linux freebsd darwin

package filesharedelete

import (
	"os"
)

func openFileForTest(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
}
