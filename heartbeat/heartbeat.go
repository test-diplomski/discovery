package heartbeat

import (
	"context"
)

type Heartbeat interface {
	Watch(ctx context.Context, f func(data string))
}
