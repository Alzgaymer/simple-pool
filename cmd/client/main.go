package main

import (
	"context"
	"log"
	"simple-pool/config"
	"simple-pool/pool"
)

func main() {
	log.Println("Started client...")
	ctx, _ := context.WithDeadline(context.Background(), config.Timeout)

	connPool, err := pool.NewConnPool()
	if err != nil {
		log.Fatal(err)
	}

	err = connPool.Dial(config.Network, config.Address, 10)
	if err != nil {
		log.Fatal(err)
	}

	err = connPool.Wait(ctx)
	if err != nil {
		log.Println(err)
	}

	err = connPool.Clear()
	if err != nil {
		log.Println(err)
		return
	}
}
