# configfile

[![Build Status](https://travis-ci.org/acoshift/configfile.svg?branch=master)](https://travis-ci.org/acoshift/configfile)
[![Coverage Status](https://coveralls.io/repos/github/acoshift/configfile/badge.svg?branch=master)](https://coveralls.io/github/acoshift/configfile?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/acoshift/configfile)](https://goreportcard.com/report/github.com/acoshift/configfile)
[![GoDoc](https://godoc.org/github.com/acoshift/configfile?status.svg)](https://godoc.org/github.com/acoshift/configfile)
[![Sourcegraph](https://sourcegraph.com/github.com/acoshift/configfile/-/badge.svg)](https://sourcegraph.com/github.com/acoshift/configfile?badge)

Read config from file, useful when read data from kubernetes configmaps, and secret.

## Example

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/acoshift/configfile"
	"github.com/garyburd/redigo/redis"
)

var config = configfile.NewReader("config")

var (
	addr      = config.StringDefault("addr", ":8080")
	redisAddr = config.MustString("redis_addr")
	redisPass = config.String("redis_pass")
	redisDB   = config.Int("redis_db")
)

func main() {
	pool := redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				redisAddr,
				redis.DialPassword(redisPass),
				redis.DialDatabase(redisDB),
			)
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := pool.Get()
		defer c.Close()
		cnt, err := redis.Int64(c.Do("INCR", "cnt"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "count: %d", cnt)
	})

	http.ListenAndServe(addr, nil)
}
```
