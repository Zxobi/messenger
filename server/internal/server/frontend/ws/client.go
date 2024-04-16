package ws

import (
	"github.com/dvid-messanger/internal/lib/logger"
	"github.com/gorilla/websocket"
	"log/slog"
	"net"
	"sync"
	"time"
)

type Client struct {
	log        *slog.Logger
	conn       *websocket.Conn
	cfg        *ClientConfig
	send       chan []byte
	registrar  Registrar
	msgHandler MsgHandler
	props      map[string]string
	mu         sync.RWMutex
}

type ClientConfig struct {
	sendMsgBuff  int
	readMsgLimit int64
	writeWait    time.Duration
	pongWait     time.Duration
	pingPeriod   time.Duration
}

type Registrar interface {
	Register(c *Client)
	Unregister(c *Client)
}

type MsgHandler interface {
	HandleMsg(c *Client, msg []byte) error
}

func NewClient(log *slog.Logger, conn *websocket.Conn, config *ClientConfig, registrar Registrar, handler MsgHandler) *Client {
	return &Client{
		log:        log,
		conn:       conn,
		cfg:        config,
		send:       make(chan []byte, config.sendMsgBuff),
		registrar:  registrar,
		msgHandler: handler,
		props:      make(map[string]string),
	}
}

func (c *Client) Serve() {
	c.registrar.Register(c)

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

func (c *Client) Send(msg []byte) {
	c.send <- msg
}

func (c *Client) SetValue(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.props[key] = value
}

func (c *Client) GetValue(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.props[key]
}

func (c *Client) readPump() {
	const op = "frontend.readPump"
	log := c.log.With(slog.String("op", op))

	defer func() {
		c.registrar.Unregister(c)
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

		if err = c.msgHandler.HandleMsg(c, message); err != nil {
			log.Error("failed to handle client msg", logger.Err(err))
			break
		}
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
