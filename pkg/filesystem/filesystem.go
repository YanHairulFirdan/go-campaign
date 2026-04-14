package filesystem

import (
	"context"
	"io"
)

type Filesystem interface {
	// CheckDirectoryExistance(ctx context.Context, path string) (bool, error)
	CreateDirectory(ctx context.Context, name string) error
	SaveFile(ctx context.Context, file io.Reader, path string) error
}
