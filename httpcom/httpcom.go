package httpcom

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cocobao/log"
	"golang.org/x/net/context/ctxhttp"
)

func Get(ctx context.Context, url string, headEx map[string]string) ([]byte, error) {
	var req *http.Request
	var res *http.Response
	var err error
	client := &http.Client{}
	tryTime := 0

tryAgain:
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if headEx != nil {
		for k, v := range headEx {
			req.Header.Set(k, v)
		}
	}

	ctxto, cancel := context.WithTimeout(ctx, 3*time.Second)
	res, err = ctxhttp.Do(ctxto, client, req)
	cancel()
	if err != nil {
		log.Warn("http get err:", err, tryTime)
		select {
		case <-ctx.Done():
			return nil, err
		default:
		}

		tryTime++
		if tryTime < 3 {
			goto tryAgain
		}
		return nil, err
	}
	if res.Body == nil {
		return nil, errors.New("get response is nil")
	}
	defer res.Body.Close()

	var result []byte
	result, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return result, fmt.Errorf("StatusCode:%d", res.StatusCode)
	}

	return result, nil
}

func Post(ctx context.Context, url string, bd []byte) ([]byte, error) {
	var req *http.Request
	var res *http.Response
	var err error
	client := &http.Client{}
	tryTime := 0

tryAgain:
	req, err = http.NewRequest("POST", url, bytes.NewReader(bd))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctxto, cancel := context.WithTimeout(ctx, 3*time.Second)
	res, err = ctxhttp.Do(ctxto, client, req)
	cancel()
	if err != nil {
		log.Warn("push post err:", err, tryTime)
		select {
		case <-ctx.Done():
			return nil, err
		default:
		}

		tryTime++
		if tryTime < 3 {
			goto tryAgain
		}
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("post result:%d", res.StatusCode)
	}
	if res.Body == nil {
		return nil, errors.New("post response is nil")
	}

	var result []byte
	result, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}
