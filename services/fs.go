package services

import "os"

type FileSystem interface {
	TempDir() string
	MkdirAll(name string, perm os.FileMode) error
	WriteFile(name string, data []byte, perm os.FileMode) error
}

type OsFs struct{}

func (OsFs) MkdirAll(name string, perm os.FileMode) error {
	return os.MkdirAll(name, perm)
}

func (OsFs) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (OsFs) TempDir() string {
	return os.TempDir()
}
