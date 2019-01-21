package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cocobao/log"
	"github.com/coreos/etcd/clientv3"
)

var (
	etcdService *EtcdService
)

type EtcdService struct {
	etcdClient *clientv3.Client
	timeout    time.Duration
	watchKeys  []string
	mux        sync.Mutex
}

func GetEtcdService() *EtcdService {
	return etcdService
}

func Setup(usname, pwd string, endPoints []string, dialTimeout time.Duration) error {
	etcdService = &EtcdService{}
	etcdService.timeout = dialTimeout
	acfg := &AuthCfg{
		Username: usname,
		Password: pwd,
	}
	if usname == "" && pwd == "" {
		acfg = nil
	}

	var err error
	etcdService.etcdClient, err = NewClient(endPoints, dialTimeout, nil, acfg)
	if err != nil {
		log.Error("connect etcd server failed!", err)
		return err
	}
	return nil
}

func (p *EtcdService) Set(key string, value interface{}) (int64, error) {
	j, err := json.Marshal(value)
	if err != nil {
		return -1, fmt.Errorf("json marshal value fail,%v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout*time.Second)
	// lp, err := p.etcdClient.Lease.Grant(ctx, 30*60)
	// cancel()
	// if err != nil {
	// 	log.Errorf("etcd lease grant fail, key:%s, val:%v, err:%v", key, value, err)
	// 	return 0, err
	// }
	// p.etcdClient.Lease.KeepAlive(context.Background(), lp.ID)

	ctx, cancel = context.WithTimeout(context.Background(), p.timeout*time.Second)
	_, err = p.etcdClient.Put(ctx, key, string(j) /*clientv3.WithLease(lp.ID)*/)
	cancel()
	if err != nil {
		// p.etcdClient.Lease.Revoke(context.Background(), lp.ID)
		return -1, err
	}
	// return int64(lp.ID), nil
	return 0, nil
}

func (p *EtcdService) Get(key string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout*time.Second)
	defer cancel()
	grep, err := p.etcdClient.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.Errorf("etcdClient.Get service fail, err:%v, keyname:%s", err, key)
		return nil, err
	}

	result := make(map[string]interface{}, len(grep.Kvs))
	for _, kv := range grep.Kvs {
		result[string(kv.Key)] = string(kv.Value)
	}
	return result, nil
}

func (p *EtcdService) Del(key string, leaseId int64) error {
	if leaseId > 0 {
		p.etcdClient.Lease.Revoke(context.Background(), clientv3.LeaseID(leaseId))
	}
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout*time.Second)
	defer cancel()
	_, err := p.etcdClient.Delete(ctx, key)
	return err
}

func (p *EtcdService) Watch(key string, cb func(t int, k string, v string) bool) {
	p.mux.Lock()
	defer p.mux.Unlock()
	for _, k := range p.watchKeys {
		if strings.Compare(k, key) == 0 {
			return
		}
	}
	p.watchKeys = append(p.watchKeys, key)

	go func() {
		defer func() {
			p.mux.Lock()
			defer p.mux.Unlock()
			for index := 0; index < len(p.watchKeys); index++ {
				if strings.Compare(p.watchKeys[index], key) == 0 {
					p.watchKeys = append(p.watchKeys[:index], p.watchKeys[index+1:]...)
					log.Debugf("watch:%s thread break", key)
					return
				}
			}

			log.Warnf("break and not found watch %s record", key)
		}()
		ctxwatch := context.Background()
		wc := p.etcdClient.Watch(ctxwatch, key, clientv3.WithPrefix())
		defer ctxwatch.Done()
		for {
			da := <-wc
			for _, change := range da.Events {
				kvk := string(change.Kv.Key)
				if !cb(int(change.Type), kvk, string(change.Kv.Value)) {
					return
				}
			}
		}
	}()
}

func (p *EtcdService) Stop() {
	p.etcdClient.Close()
}
