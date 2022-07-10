package services

import "os"

type IFileService interface {
	TempDir() string
	MkdirAll(path string, perm os.FileMode) error
	WriteFile(path string, data []byte, perm os.FileMode) error
	ReadFile(path string) (content []byte, err error)
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

func (OsFs) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
