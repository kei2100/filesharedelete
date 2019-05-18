package filesharedelete

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestOpenFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "fsdtest")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	of, err := openFileForTest(filepath.Join(dir, "foo"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	f := &onceCloseFile{File: of}
	defer f.Close()

	if _, err := f.WriteString("foo"); err != nil {
		t.Fatalf("failed to write foo: %+v", err)
	}
	renamed := f.Name() + ".bk"
	if err := os.Rename(f.Name(), renamed); err != nil {
		t.Fatalf("failed to rename: %+v", err)
	}
	if _, err := f.WriteString("bar"); err != nil {
		t.Fatalf("faield to write bar: %+v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("faield to close: %+v", err)
	}

	b, err := ioutil.ReadFile(renamed)
	if err != nil {
		t.Fatalf("failed to read %s: %+v", renamed, err)
	}
	if g, w := string(b), "foobar"; g != w {
		t.Errorf("content got %v, want %v", g, w)
	}
}

type onceCloseFile struct {
	once sync.Once
	*os.File
}

func (ocf *onceCloseFile) Close() error {
	var err error
	ocf.once.Do(func() {
		err = ocf.File.Close()
	})
	return err
}
