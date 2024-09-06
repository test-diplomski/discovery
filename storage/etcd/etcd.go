package etcd

import (
	"context"
	"fmt"
	"github.com/c12s/discovery/strategy"
	"github.com/coreos/etcd/clientv3"
	"time"
)

type ETCD struct {
	kv       clientv3.KV
	client   *clientv3.Client
	services map[string][]string
	s        strategy.Strategy
}

func New(endpoints []string, timeout time.Duration, s strategy.Strategy) (*ETCD, error) {
	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: timeout,
		Endpoints:   endpoints,
	})

	if err != nil {
		return nil, err
	}

	return &ETCD{
		kv:       clientv3.NewKV(cli),
		client:   cli,
		services: make(map[string][]string),
		s:        s,
	}, nil
}

func (e *ETCD) Store(ctx context.Context, name string) (bool, error) {
	// minimum lease TTL is 10-second
	resp, err := e.client.Grant(ctx, 10)
	if err != nil {
		return false, err
	}

	// after 10 seconds, the key will be removed
	_, err = e.client.Put(ctx, name, "", clientv3.WithLease(resp.ID))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (e *ETCD) Get(ctx context.Context, name string) (string, error) {
	elems := e.services[join(name)]
	size := len(elems)
	i, err := e.s.Next(ctx, size)
	if err != nil {
		return "", err
	}
	return elems[i], nil
}

func (e *ETCD) Watcher(ctx context.Context) {
	go func() {
		rch := e.client.Watch(ctx, "/heartbeat/", clientv3.WithPrefix())
		for {
			select {
			case <-ctx.Done():
				fmt.Println(ctx.Err())
				return
			case result := <-rch:
				for _, ev := range result.Events {
					switch ev.Type {
					case clientv3.EventTypePut:
						e.put(string(ev.Kv.Key))
						// log.Println("PUT: ", string(ev.Kv.Key))
					case clientv3.EventTypeDelete:
						e.del(string(ev.Kv.Key))
						// log.Println("DELETE: ", string(ev.Kv.Key))
					}
				}
			}
		}
	}()
}
