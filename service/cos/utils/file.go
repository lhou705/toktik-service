package utils

// 实现一个虚拟文件对象用于文件上传
import (
	"bytes"
	"os"
	"time"
)

type FileInfo struct {
	name string
	Data []byte
}

func (fi FileInfo) Name() string { return fi.name }

func (fi FileInfo) Size() int64 { return int64(len(fi.Data)) }

func (fi FileInfo) Mode() os.FileMode { return 0444 }

func (fi FileInfo) ModTime() time.Time { return time.Time{} }

func (fi FileInfo) IsDir() bool { return false }

func (fi FileInfo) Sys() any { return nil }

type File struct {
	*bytes.Reader
	FileInfo FileInfo
}

func (f *File) Close() error { return nil } // Noop, nothing to do
func (f *File) Readdir(int) ([]os.FileInfo, error) {
	return nil, nil // We are not a directory but a single file
}
func (f *File) Stat() (os.FileInfo, error) {
	return f.FileInfo, nil
}
