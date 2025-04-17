package fileIO

import (
	"os"
	"time"
)

// fakeFileInfo is a mock implementation of os.FileInfo used to simulate test paths
// without relying on the actual filesystem structure.
type FakeFileInfo struct {
	NameVal  string
	IsDirVal bool
}

// Name returns the file name (mocked)
func (f *FakeFileInfo) Name() string {
	return f.NameVal
}

// Size returns the file size (mocked)
func (f *FakeFileInfo) Size() int64 {
	return 0
}

// Mode returns file permissions (mocked)
func (f *FakeFileInfo) Mode() os.FileMode {
	return 0644
}

// ModTime returns a zero timestamp (mocked)
func (f *FakeFileInfo) ModTime() time.Time {
	return time.Time{}
}

// IsDir reports whether the file is a directory (mocked)
func (f *FakeFileInfo) IsDir() bool {
	return f.IsDirVal
}

// Sys is unused and returns nil (mocked)
func (f *FakeFileInfo) Sys() interface{} {
	return nil
}
