package main

import (
	"context"
	"log"
	"simple-pool/config"
	"simple-pool/pool"
)

func main() {
	log.Println("Started server...")
	ctx, _ := context.WithDeadline(context.Background(), config.Timeout)

	connPool, err := pool.NewConnPool(pool.WithListnerConfig(ctx, config.Network, config.Address))
	if err != nil {
		log.Fatal(err)
	}

	go connPool.AcceptConns(ctx)

	if err = connPool.Wait(ctx); err != nil {
		log.Println(err)
	}

	err = connPool.Clear()
	if err != nil {
		log.Println(err)
		return
	}
}
