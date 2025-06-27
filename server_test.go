package main

import (
	"context"
	"fmt"
	"goredis/client"
	"log"
	"sync"
	"testing"
	"time"
)

func TestServerWithMultiClients(t *testing.T) {
	server := NewServer(Config{})
	go func() {
		log.Fatal(server.Start())
	}()

	time.Sleep(time.Second)

	nClients := 10
	var wg sync.WaitGroup
	wg.Add(nClients)

	for i := 0; i < nClients; i++ {
		go func(it int) {
			c, err := client.New("localhost:5001")
			if err != nil {
				log.Fatal(err)
			}

			defer c.Close()

			key := fmt.Sprintf("client_foo_%d", i)
			set_val := fmt.Sprintf("client_bar_%d", i)

			if err := c.Set(context.Background(), key, set_val); err != nil {
				log.Fatal(err)
			}

			get_val, err := c.Get(context.Background(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("client %d got this value back%s\n", i, get_val)
			wg.Done()
		}(i)
	}

	wg.Wait()

	time.Sleep(time.Second)

	if len(server.peers) != 0 {
		t.Fatalf("expected 0 peers but got %d", len(server.peers))
	}
}
