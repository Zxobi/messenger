package ws

import (
	"github.com/dvid-messanger/internal/pkg/logger"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"time"
)

type WsServer struct {
	log       *slog.Logger
	upgrader  websocket.Upgrader
	clientCfg ClientConfig
	registrar Registrar
	handler   MsgHandler
}

func NewWsServer(
	log *slog.Logger,
	registrar Registrar,
	handler MsgHandler,
	sendBuffSize int,
	rBuffSize int,
	wBuffSize int,
	hsTimeout time.Duration,
	msgLimit int64,
	writeWait time.Duration,
	pongWait time.Duration,
) *WsServer {
	return &WsServer{
		log: log,
		upgrader: websocket.Upgrader{
			HandshakeTimeout: hsTimeout,
			ReadBufferSize:   rBuffSize,
			WriteBufferSize:  wBuffSize,
		},
		registrar: registrar,
		handler:   handler,
		clientCfg: ClientConfig{
			sendMsgBuff:  sendBuffSize,
			readMsgLimit: msgLimit,
			writeWait:    writeWait,
			pongWait:     pongWait,
			pingPeriod:   (pongWait * 9) / 10,
		},
	}
}

func (ws *WsServer) Handle(w http.ResponseWriter, r *http.Request) {
	const op = "websocket.Handle"
	log := ws.log.With(slog.String("op", op))

	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("connection upgrade failed", logger.Err(err))
		return
	}

	NewClient(ws.log, conn, &ws.clientCfg, ws.registrar, ws.handler).Serve()
}
