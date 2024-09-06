package storage

import (
	"context"
)

type DB interface {
	Store(ctx context.Context, name string) (bool, error)
	Get(ctx context.Context, name string) (string, error)
	Watcher(ctx context.Context)
}
