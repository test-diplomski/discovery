package basic

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

type BasicStrategy struct {
}

func (s *BasicStrategy) Next(ctx context.Context, size int) (int, error) {
	rand.Seed(time.Now().Unix())
	if size > 0 {
		return rand.Int() % size, nil
	}
	return -1, errors.New("No registered services")
}

func NewStrategy() (*BasicStrategy, error) {
	return &BasicStrategy{}, nil
}
