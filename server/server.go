package server

import (
	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"gobot/config"
	telegabot "gobot/telegram"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	Host      string
	Port      string
	Log       *logrus.Logger
	Router    *router.Router
	TelegaApi *telegabot.TelegaApi
}

func NewServer(log *logrus.Logger, telega *telegabot.TelegaApi, config *config.Config) *Server {
	s := &Server{
		Host:      config.Server.Host,
		Port:      config.Server.Port,
		Log:       log,
		TelegaApi: telega,
		Router:    router.New(),
	}

	s.Router.GET("/", s.HandlerIndex)
	s.Router.GET("/send", s.HandlerSendTelegramMessage)

	return s
}

func (s *Server) ListenAndServe() error {
	var g errgroup.Group

	g.Go(func() error {
		return fasthttp.ListenAndServe(s.Host+":"+s.Port, s.Router.Handler)
	})
	s.Log.Info("Server started!")

	return g.Wait()
}

func (s *Server) HandlerIndex(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("Welcome!")
}

func (s *Server) HandlerSendTelegramMessage(ctx *fasthttp.RequestCtx) {
	msg := string(ctx.QueryArgs().Peek("msg"))
	token := string(ctx.QueryArgs().Peek("token"))

	if token != "77QrYVUwBwWcsszzqw586x7Y9dukHNCr" || msg == "" {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	s.TelegaApi.SendMessages(msg)

	ctx.WriteString("Sended!")
}
