package main

import (
	"github.com/c12s/discovery/heartbeat/nats"
	"github.com/c12s/discovery/model/config"
	"github.com/c12s/discovery/service"
	"github.com/c12s/discovery/storage/etcd"
	"github.com/c12s/discovery/strategy/basic"
	"log"
	"time"
)

func main() {
	conf, err := config.ConfigFile()
	if err != nil {
		log.Fatal(err)
		return
	}

	s, err := basic.NewStrategy()
	if err != nil {
		log.Fatal(err)
		return
	}

	db, err := etcd.New(conf.Db, 10*time.Second, s)
	if err != nil {
		log.Fatal(err)
		return
	}

	w, err := nats.New(conf.Heartbeat, conf.HeartbeatTopic)
	if err != nil {
		log.Fatal(err)
		return
	}

	service.Run(
		conf.ConfVersion,
		conf.Address,
		db,
		w,
	)
}
