package ws

import (
	"github.com/dvid-messanger/internal/adapter/primary"
	"github.com/dvid-messanger/internal/pkg/logger"
	"github.com/gorilla/websocket"
	"log/slog"
	"net"
	"sync"
	"time"
)

type Client struct {
	id       []byte
	log      *slog.Logger
	conn     *websocket.Conn
	registry primary.ClientRegistry
	handler  MsgHandler
	cfg      *ClientConfig
	send     chan []byte
}

type ClientConfig struct {
	sendMsgBuff  int
	readMsgLimit int64
	writeWait    time.Duration
	pongWait     time.Duration
	pingPeriod   time.Duration
}

type MsgHandler interface {
	Handle(c *Client, msg []byte)
}

func NewClient(
	id []byte,
	log *slog.Logger,
	registry primary.ClientRegistry,
	handler MsgHandler,
	conn *websocket.Conn,
	config *ClientConfig,
) *Client {
	return &Client{
		id:       id,
		log:      log,
		registry: registry,
		handler:  handler,
		conn:     conn,
		cfg:      config,
		send:     make(chan []byte, config.sendMsgBuff),
	}
}

func (c *Client) Serve() {
	const op = "client.Serve"
	log := c.log.With(slog.String("op", op))

	if err := c.registry.Register(c); err != nil {
		log.Error("failed to register client", logger.Err(err))
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		c.readPump()
	}()
	go func() {
		defer wg.Done()
		c.writePump()
	}()

	wg.Wait()
}

func (c *Client) GetAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Client) GetId() []byte {
	return c.id
}

func (c *Client) Send(msg []byte) error {
	c.send <- msg
	return nil
}

func (c *Client) readPump() {
	const op = "frontend.readPump"
	log := c.log.With(slog.String("op", op))

	defer func() {
		if err := c.registry.Unregister(c.id); err != nil {
			log.Error("failed to unregister client", logger.Err(err))
		}
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(c.cfg.readMsgLimit)
	_ = c.conn.SetReadDeadline(time.Now().Add(c.cfg.pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(c.cfg.pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Error("failed to read client msg", logger.Err(err))
			break
		}

		c.handler.Handle(c, message)
	}
}

func (c *Client) writePump() {
	const op = "frontend.writePump"
	log := c.log.With(slog.String("op", op))

	ticker := time.NewTicker(c.cfg.pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(c.cfg.writeWait))
			if !ok {
				log.Debug("send channel closed")
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
				log.Error("failed to write message", logger.Err(err))
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(c.cfg.writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error("failed to write ping message", logger.Err(err))
				return
			}
		}
	}
}
