package route

import (
	"context"
	"github.com/dvid-messanger/internal/adapter/primary/frontend/ws"
	"github.com/dvid-messanger/internal/pkg/logger"
	"github.com/dvid-messanger/protos/gen/frontend"
	"github.com/golang/protobuf/proto"
	"log/slog"
	"reflect"
	"time"
)

const (
	defaultRequestTimeout = 5 * time.Second
)

type HandlerFunc = func(context.Context, *UpstreamRequest) *UpstreamResponse

type Handler struct {
	dt  frontendv1.DownstreamType
	msg proto.Message
	fun HandlerFunc
}

type Router struct {
	log    *slog.Logger
	routes map[frontendv1.UpstreamType]*Handler
}

func NewRouter(log *slog.Logger) *Router {
	return &Router{log: log, routes: make(map[frontendv1.UpstreamType]*Handler)}
}

func (r *Router) RegisterHandler(ut frontendv1.UpstreamType, dt frontendv1.DownstreamType, msg proto.Message, fun HandlerFunc) {
	const op = "route.HandlerFunc"
	log := r.log.With(slog.String("op", op))

	if _, ok := r.routes[ut]; ok {
		log.Error("handler for upstream " + ut.String() + " already registered")
		return
	}

	r.routes[ut] = &Handler{dt: dt, msg: msg, fun: fun}
}

func (r *Router) Handle(c *ws.Client, msg []byte) {
	const op = "upstream.HandleMsg"
	log := r.log.With(slog.String("op", op), slog.String("cl", c.GetAddr().String()))

	upstream := &frontendv1.Upstream{}
	if err := proto.Unmarshal(msg, upstream); err != nil {
		log.Error("failed to unmarshal upstream", logger.Err(err))
		return
	}

	log = log.With(slog.String("ut", upstream.Type.String()))

	handler := r.routes[upstream.Type]
	if handler == nil {
		log.Error("handler not found")
		return
	}

	payload := reflect.New(reflect.TypeOf(handler.msg).Elem()).Interface().(proto.Message)
	err := proto.Unmarshal(upstream.Payload, payload)
	if err != nil {
		log.Error("failed to unmarshal upstream payload", logger.Err(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	resp := handler.fun(ctx, &UpstreamRequest{ClientId: c.GetId(), Payload: payload})
	if ctx.Err() != nil {
		log.Error("request timeout")
		resp = ErrResponseTimeout
	}
	if resp == nil {
		log.Error("nil response")
		resp = ErrResponseInternal
	}

	downstream, err := MarshalResponse(resp, handler.dt)
	if err != nil {
		log.Error("failed to marshal response", logger.Err(err))
		return
	}

	if downstream != nil {
		_ = c.Send(downstream)
	}
	log.Debug("msg handled")
}
