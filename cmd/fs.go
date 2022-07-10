package cmd

import "os"

type fileSystem interface {
	TempDir() string
	MkdirAll(name string, perm os.FileMode) error
	WriteFile(name string, data []byte, perm os.FileMode) error
}

type osFs struct{}

func (osFs) MkdirAll(name string, perm os.FileMode) error {
	return os.MkdirAll(name, perm)
}

func (osFs) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (osFs) TempDir() string {
	return os.TempDir()
}

var fs fileSystem = osFs{}
