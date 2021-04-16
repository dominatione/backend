package grpc

import (
	"context"
	"fmt"
	"github.com/dominati-one/backend/pkg/protocol/gameapi"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"time"
)

type Server struct {
	log           zerolog.Logger
	listenPort    uint32
	listenAddress string
	socket        net.Listener

	gameapiHandler *GameApiHandler
}

func NewServer(listenPort uint32, listenAddress string, gameapiHandler *GameApiHandler) *Server {
	return &Server{
		log:            log.With().Str("applicationComponent", "grpcApiServer").Logger(),
		listenPort:     listenPort,
		listenAddress:  listenAddress,
		gameapiHandler: gameapiHandler,
	}
}

func (s *Server) Start(ctx context.Context) error {
	var err error
	listenAddress := fmt.Sprintf("%s:%d", s.listenAddress, s.listenPort)
	log := log.With().Str("listenAddress", listenAddress).Logger()

	s.socket, err = net.Listen("tcp", listenAddress)
	if err != nil {
		return errors.Wrap(err, "unable to start network listener")
	}

	server := grpc.NewServer()

	gameapi.RegisterApiServer(server, s.gameapiHandler)

	go func() {
		for {
			log.Info().Msg("Listening for connections.")
			if err := server.Serve(s.socket); err != nil {
				backOffDuration := time.Second * 5
				log.Error().Dur("backOffDuration", backOffDuration).
					Err(err).
					Msg("Unable to start gRPC API server. Backing off.")
				time.Sleep(backOffDuration)
			}
		}
	}()

	return nil
}
