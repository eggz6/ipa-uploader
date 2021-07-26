package uploader

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type fs struct {
	root string
}

func NewFS() *fs {
	return &fs{root: "./output"}
}

func (f *fs) Upload(ctx context.Context, r io.Reader, path string) (string, error) {
	path = f.root + path
	ps, _ := filepath.Split(path)
	err := os.MkdirAll(ps, os.ModePerm)
	if err != nil {
		return "", err
	}

	dest, err := os.Create(path)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(dest, r)
	if err != nil {
		return "", err
	}

	return "http://localhost:8080/" + path, nil
}
