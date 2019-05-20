package filesharedelete

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestRenameFileWhileOpening(t *testing.T) {
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

func TestDeleteFileWhileOpening(t *testing.T) {
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
	if err := os.Remove(f.Name()); err != nil {
		t.Fatalf("failed to remove foo: %+v", err)
	}
	if _, err := f.WriteString("bar"); err != nil {
		t.Fatalf("faield to write bar: %+v", err)
	}

	if _, err := f.Seek(0, 0); err != nil {
		t.Fatalf("failed to seek: %+v", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read: %+v", err)
	}
	if g, w := string(b), "foobar"; g != w {
		t.Errorf("content got %v, want %v", g, w)
	}

	if err := f.Close(); err != nil {
		t.Fatalf("faield to close: %+v", err)
	}
	if fi, err := os.Stat(f.Name()); err == nil || !os.IsNotExist(err) {
		t.Fatalf("unexpected file stat fi=%+v err=%+v", fi, err)
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
