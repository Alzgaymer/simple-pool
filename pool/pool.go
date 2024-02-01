package pool

import (
	"context"
	"errors"
	"log"
	"net"
)

type ConnPool struct {
	conns    []net.Conn
	ctx      context.Context
	listener net.Listener
	connsCh  chan net.Conn
}

func (c *ConnPool) HandleConns() error {
	for {
		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		case conn := <-c.connsCh:
			c.conns = append(c.conns, conn)
			log.Printf("#%d Connection accepted", len(c.conns))
		}
	}
}

func (c *ConnPool) AcceptConns() error {
	for {
		select {
		case <-c.ctx.Done():
			log.Println("Server is shutting down...")
			err := c.listener.Close()
			if err != nil {
				return err
			}

			return c.ctx.Err()
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
		c.conns = append(c.conns, conn)
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

func (c *ConnPool) Wait() error {
	select {
	case <-c.ctx.Done():
		return errors.Join(c.ctx.Err(), c.Close())
	}
}

func NewConnPool(ctx context.Context, opts ...Option) (*ConnPool, error) {
	var connPool = &ConnPool{
		conns:   make([]net.Conn, 0),
		ctx:     ctx,
		connsCh: make(chan net.Conn),
	}

	for _, option := range opts {
		if err := option(connPool); err != nil {
			return nil, err
		}
	}

	return connPool, nil
}

type Option func(pool *ConnPool) error

func WithListner(l net.Listener) Option {
	return func(pool *ConnPool) error {
		pool.listener = l
		return nil
	}
}

func WithListnerConfig(network, address string) Option {
	return func(pool *ConnPool) error {
		var lc net.ListenConfig
		listener, err := lc.Listen(pool.ctx, network, address)
		if err != nil {
			return err
		}
		pool.listener = listener
		return nil
	}
}
