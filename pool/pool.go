package pool

import (
	"context"
	"errors"
	"log"
	"net"
)

type ConnPool struct {
	conns    []net.Conn
	listener net.Listener
	connsCh  chan net.Conn
}

func (c *ConnPool) AcceptConns(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Println("Server is shutting down...")
			return ctx.Err()
		default:
			accept, err := c.listener.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			c.connsCh <- accept
		}
	}
}

func (c *ConnPool) Dial(network, address string, amount int64) error {
	for i := 0; int64(i) < amount; i++ {
		conn, err := net.Dial(network, address)
		if err != nil {
			return err
		}
		log.Printf("#%d Succesful dial", i)
		c.connsCh <- conn
	}
	return nil
}

func (c *ConnPool) Close() error {
	for i, conn := range c.conns {
		err := conn.Close()
		if err != nil {
			log.Printf("#%d Conn close error: %s", i, err.Error())
		}
		log.Printf("#%d Conn closed", i)
	}
	return nil
}

func (c *ConnPool) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return errors.Join(ctx.Err())
	}
}

func NewConnPool(opts ...Option) (*ConnPool, error) {
	var connPool = &ConnPool{
		conns:   make([]net.Conn, 0),
		connsCh: make(chan net.Conn),
	}

	for _, option := range opts {
		if err := option(connPool); err != nil {
			return nil, err
		}
	}

	go connPool.addConns()

	return connPool, nil
}

func (c *ConnPool) Clear() error {
	// no need to call c.listner.Close,
	// because it closes by itself after ctx is done or canceled

	err := c.Close()
	if err != nil {
		return err
	}

	close(c.connsCh)
	return nil
}

func (c *ConnPool) addConns() {
	var count int
	for conn := range c.connsCh {
		log.Printf("#%d Connection accepted:", count)
		count++
		c.conns = append(c.conns, conn)
	}
}

type Option func(pool *ConnPool) error

func WithListner(l net.Listener) Option {
	return func(pool *ConnPool) error {
		pool.listener = l
		return nil
	}
}

func WithListnerConfig(ctx context.Context, network, address string) Option {
	return func(pool *ConnPool) error {
		var lc net.ListenConfig
		listener, err := lc.Listen(ctx, network, address)
		if err != nil {
			return err
		}
		pool.listener = listener
		return nil
	}
}
