package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cocobao/comm_go/etcd"
	"github.com/cocobao/log"
	"golang.org/x/net/context/ctxhttp"
)

var (
	index = 0

	httpServs map[string]*HttpServiceInfo
)

func init() {
	httpServs = make(map[string]*HttpServiceInfo, 3)
}

func HttpService(serviceName string) *HttpServiceInfo {
	if v, ok := httpServs[serviceName]; ok {
		return v
	}

	s := &HttpServiceInfo{
		ServiceName: serviceName,
	}

	s.LoadServiceAddrs()
	etcd.GetEtcdService().Watch(serviceName, s.watchCall)
	httpServs[serviceName] = s
	return s
}

type AddrInfo struct {
	Name string
	IP   string
	Port int
}

type HttpServiceInfo struct {
	ServiceName  string
	ServiceAddrs []AddrInfo
}

func (s *HttpServiceInfo) watchCall(t int, k string, v string) bool {
	s.LoadServiceAddrs()
	return true
}

func (s *HttpServiceInfo) LoadServiceAddrs() error {
	m, err := etcd.GetEtcdService().Get(s.ServiceName)
	if err != nil {
		log.Errorf("get etcd service:%s fail", s.ServiceName)
		return err
	}
	s.ServiceAddrs = make([]AddrInfo, len(m))
	i := 0
	for k, val := range m {
		var adval map[string]interface{}
		str := val.(string)
		if err = json.Unmarshal([]byte(str), &adval); err != nil {
			return err
		}

		var t AddrInfo
		t.Name = k
		if vv, ok := adval["addr"].(string); ok {
			t.IP = vv
		}

		if vv, ok := adval["port"].(float64); ok {
			t.Port = int(vv)
		}

		if len(t.IP) > 0 && t.Port > 0 {
			s.ServiceAddrs[i] = t
		}
		i++
	}
	return nil
}

func (s *HttpServiceInfo) Get(url string, result interface{}) (int, error) {
	var req *http.Request
	var res *http.Response
	var err error

	ctx := context.Background()
	client := &http.Client{}
	tryTime := 0

	index++
	if index >= len(s.ServiceAddrs) {
		index = 0
	}

	url = fmt.Sprintf("http://%s:%d%s", s.ServiceAddrs[index].IP, s.ServiceAddrs[index].Port, url)
tryAgain:
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	ctxto, cancel := context.WithTimeout(ctx, 3*time.Second)
	res, err = ctxhttp.Do(ctxto, client, req)
	cancel()
	if err != nil {
		log.Warn("push post err:", err, tryTime)
		select {
		case <-ctx.Done():
			return -1, err
		default:
		}

		tryTime++
		if tryTime < 3 {
			goto tryAgain
		}
		return -1, err
	}
	if res.Body == nil {
		return -1, fmt.Errorf("post response is nil")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return res.StatusCode, fmt.Errorf("post result:%d", res.StatusCode)
	}
	var rdata []byte
	rdata, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, err
	}

	if result == nil {
		return res.StatusCode, nil
	}
	return res.StatusCode, json.Unmarshal(rdata, result)
}

func (s *HttpServiceInfo) Post(url string, bd interface{}, result interface{}) (int, error) {
	var req *http.Request
	var res *http.Response
	var err error

	ctx := context.Background()
	client := &http.Client{}
	tryTime := 0

	var data []byte
	if bd != nil {
		data, err = json.Marshal(bd)
	} else {
		data, err = json.Marshal(map[string]interface{}{})
	}

	if err != nil {
		return -1, err
	}

	index++
	if index >= len(s.ServiceAddrs) {
		index = 0
	}
	url = fmt.Sprintf("http://%s:%d%s", s.ServiceAddrs[index].IP, s.ServiceAddrs[index].Port, url)
tryAgain:
	req, err = http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	ctxto, cancel := context.WithTimeout(ctx, 3*time.Second)
	res, err = ctxhttp.Do(ctxto, client, req)
	cancel()
	if err != nil {
		log.Warn("push post err:", err, tryTime)
		select {
		case <-ctx.Done():
			return -1, err
		default:
		}

		tryTime++
		if tryTime < 3 {
			goto tryAgain
		}
		return -1, err
	}
	if res.Body == nil {
		return -1, fmt.Errorf("post response is nil")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return res.StatusCode, fmt.Errorf("post result:%d", res.StatusCode)
	}
	var rdata []byte
	rdata, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, err
	}

	if result == nil {
		return res.StatusCode, nil
	}
	return res.StatusCode, json.Unmarshal(rdata, result)
}
