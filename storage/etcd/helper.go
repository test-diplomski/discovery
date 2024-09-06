package etcd

import (
	"strings"
)

func (e *ETCD) in(key string) (bool, string, string) {
	parts := strings.Split(key, "|")
	mkey := parts[0]
	skey := parts[1]
	if _, ok := e.services[mkey]; ok {
		return true, mkey, skey
	}
	return false, mkey, skey
}

func (e *ETCD) cache(mkey, skey string) {
	for _, item := range e.services[mkey] {
		if item == skey {
			return
		}
	}
	e.services[mkey] = append(e.services[mkey], skey)
}

func (e *ETCD) put(key string) {
	if ok, mkey, skey := e.in(key); ok {
		e.cache(mkey, skey)
	} else {
		e.services[mkey] = append(e.services[mkey], skey)
	}
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func (e *ETCD) del(key string) {
	if ok, mkey, skey := e.in(key); ok {
		e.services[mkey] = remove(e.services[mkey], skey)
	}
}

func join(name string) string {
	return strings.Join([]string{"/heartbeat", name}, "/")
}
