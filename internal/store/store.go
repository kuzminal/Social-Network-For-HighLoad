package store

import (
	"context"
	"io"
)

type Store interface {
	io.Closer
	Ping(ctx context.Context) error
}
