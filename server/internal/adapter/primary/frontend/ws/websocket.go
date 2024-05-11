package ws

import (
	"github.com/dvid-messanger/internal/adapter/primary"
	"github.com/dvid-messanger/internal/pkg/logger"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	log       *slog.Logger
	upgrader  websocket.Upgrader
	clientCfg ClientConfig
	registry  primary.ClientRegistry
	handler   MsgHandler
}

func NewWsServer(
	log *slog.Logger,
	registry primary.ClientRegistry,
	handler MsgHandler,
	sendBuffSize int,
	rBuffSize int,
	wBuffSize int,
	hsTimeout time.Duration,
	msgLimit int64,
	writeWait time.Duration,
	pongWait time.Duration,
) *Server {
	return &Server{
		log: log,
		upgrader: websocket.Upgrader{
			HandshakeTimeout: hsTimeout,
			ReadBufferSize:   rBuffSize,
			WriteBufferSize:  wBuffSize,
		},
		registry: registry,
		handler:  handler,
		clientCfg: ClientConfig{
			sendMsgBuff:  sendBuffSize,
			readMsgLimit: msgLimit,
			writeWait:    writeWait,
			pongWait:     pongWait,
			pingPeriod:   (pongWait * 9) / 10,
		},
	}
}

func (s *Server) Handle(w http.ResponseWriter, r *http.Request) {
	const op = "websocket.Handle"
	log := s.log.With(slog.String("op", op))

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("connection upgrade failed", logger.Err(err))
		return
	}

	id := uuid.New()
	NewClient(id[:], s.log, s.registry, s.handler, conn, &s.clientCfg).Serve()
}
