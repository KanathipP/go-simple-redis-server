package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tidwall/resp"
)

func TestOfficialRedisClient(t *testing.T) {
	listenAddr := ":5002"
	server := NewServer(Config{
		ListenAddr: listenAddr,
	})

	go func() {
		log.Fatal(server.Start())
	}()

	time.Sleep(time.Second)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost%s", listenAddr),
		Password: "",
		DB:       0,
	})

	testCases := map[string]string{
		"foo":  "bar",
		"your": "mom",
		"step": "dad",
		"im":   "stuck",
	}

	for key, val := range testCases {
		err := rdb.Set(context.Background(), key, val, 0).Err()
		if err != nil {
			t.Fatal(err)
		}

		newVal, err := rdb.Get(context.Background(), key).Result()
		if err != nil {
			t.Fatal(err)
		}

		if val != newVal {
			t.Fatalf("expected %s but got %s", val, newVal)
		}
	}
}

func TestFooBar(t *testing.T) {
	buf := &bytes.Buffer{}
	rw := resp.NewWriter(buf)
	rw.WriteString("OK")
	fmt.Println(buf.String())
	in := map[string]string{
		"first":  "1",
		"second": "2",
	}
	out := respWriteMap(in)
	fmt.Println(string(out))
}
