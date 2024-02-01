package main

import (
	"connection-accepter/config"
	"connection-accepter/pool"
	"context"
	"log"
)

func main() {
	log.Println("Started client...")
	ctx, _ := context.WithDeadline(context.Background(), config.Timeout)

	connPool, err := pool.NewConnPool(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = connPool.Dial(config.Network, config.Address, 1000)
	if err != nil {
		log.Fatal(err)
	}

	err = connPool.Wait()
	if err != nil {
		log.Println(err)
	}
}
