package strategy

import (
	"context"
)

type Strategy interface {
	Next(ctx context.Context, size int) (int, error)
}
