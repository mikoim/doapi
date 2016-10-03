package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	BaseURL = "http://duet.doshisha.ac.jp/info/KK1000.jsp?katei=%d&kouchi=%d&youbi=%d"
)

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

var (
	pool *redis.Pool
)

func init() {
	pool = newPool()
}

func getCache(url string) (*[]byte, error) {
	conn := pool.Get()
	defer conn.Close()

	v, err := redis.Bytes(conn.Do("GET", url))
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func setCache(url string, v *[]byte) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", url, *v)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", url, "1800")
	if err != nil {
		return err
	}

	return nil
}

func Get(katei int, kouchi int, youbi int) (io.Reader, error) {
	url := fmt.Sprintf(BaseURL, katei, kouchi, youbi)

	a, err := getCache(url)
	if err == nil {
		return bytes.NewReader(*a), nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body := bytes.Buffer{}
	body.ReadFrom(resp.Body)

	tmp := body.Bytes()
	err = setCache(url, &tmp)
	if err != nil {
		log.Println(err)
	}

	return bytes.NewReader(body.Bytes()), nil
}
