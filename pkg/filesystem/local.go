package filesystem

import (
	"context"
	"io"
	"os"
)

type localFileSystem struct {
}

func NewLocalFileSystem() Filesystem {
	return &localFileSystem{}
}

var _ Filesystem = (*localFileSystem)(nil)

func (l *localFileSystem) CreateDirectory(ctx context.Context, name string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if _, err := os.Stat(name); os.IsNotExist(err) {
		if err := os.MkdirAll(name, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (l *localFileSystem) SaveFile(ctx context.Context, reader io.Reader, path string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	destination, err := os.Create(path)

	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, reader)

	return err
}
