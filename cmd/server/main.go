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

	connPool, err := pool.NewConnPool(ctx, pool.WithListnerConfig(config.Network, config.Address))
	if err != nil {
		log.Fatal(err)
	}

	go connPool.AcceptConns()

	go connPool.HandleConns()

	if err = connPool.Wait(); err != nil {
		log.Println(err)
	}

}
