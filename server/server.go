package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/MarkTBSS/assessment-tax/config"
	"github.com/MarkTBSS/assessment-tax/databases"
	"github.com/labstack/echo/v4"
)

type echoServer struct {
	app  *echo.Echo
	db   databases.Database
	conf *config.Config
}

func (s *echoServer) healthCheck(pctx echo.Context) error {
	return pctx.String(http.StatusOK, "OK")
}

func (s *echoServer) gracefullyShutdown(quitCh <-chan os.Signal) {
	ctx := context.Background()
	<-quitCh
	//s.app.Logger.Infof("Shutting down service...")
	log.Printf("shutting down the server")
	if err := s.app.Shutdown(ctx); err != nil {
		s.app.Logger.Fatalf("Error: %s", err.Error())
	}
}

func (s *echoServer) httpListening() {
	url := fmt.Sprintf(":%s", s.conf.Server.Port)

	if err := s.app.Start(url); err != nil && err != http.ErrServerClosed {
		s.app.Logger.Fatalf("Error: %v", err)
	}
}

func (s *echoServer) Start() {
	s.app.GET("/health", s.healthCheck)
	s.initTaxRouter()

	// Graceful shutdown
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)
	go s.gracefullyShutdown(quitCh)

	s.httpListening()
}

var (
	server *echoServer
	once   sync.Once
)

func NewEchoServer(conf *config.Config, db databases.Database) *echoServer {
	echoApp := echo.New()
	//echoApp.Logger.SetLevel(log.DEBUG)
	once.Do(func() {
		server = &echoServer{
			app:  echoApp,
			db:   db,
			conf: conf,
		}
	})
	log.Printf("Admin Username : %s", conf.Server.Username)
	log.Printf("Admin Password : %s", conf.Server.Password)
	return server
}
