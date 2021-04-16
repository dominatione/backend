package backend

import (
	"context"
	"github.com/dominati-one/backend/internal/app/backend/api/grpc"
	"github.com/dominati-one/backend/internal/pkg/blockchain"
	"github.com/dominati-one/backend/internal/pkg/blockchain/local"
	"github.com/dominati-one/backend/internal/pkg/blockchain/network"
	"github.com/dominati-one/backend/internal/pkg/game"
	"github.com/pkg/errors"
	"time"
)

type AppParameters struct {
	GrpcApiListenPort    uint32
	GrpcApiListenAddress string
}

type App struct {
	parameters          AppParameters
	game                *game.Game
	grpcApiServer       *grpc.Server
	blockchainConnector blockchain.Connector
	blockchain          *blockchain.Network
	eventPump           *EventPump
}

func NewApp(parameters AppParameters) *App {
	game := game.NewGame()

	blockchainSettings := blockchain.NetworkSettings{
		BlockInterval:       10 * time.Second,
		AuthorityPublicKeys: network.CreateTestNetAuthority(),
		GenesisBlock:        network.CreateTestNetGenesisBlock(),
	}

	blockchainConnector := local.NewConnector()
	blockchainEventStorage := NewEventStorage()
	blockchainBlockStorage := NewBlockStorage()
	blockchain := blockchain.NewNetwork(blockchainSettings, blockchainConnector, blockchainEventStorage, blockchainBlockStorage)

	gameApiHandler := grpc.NewGameApiHandler(game, blockchain.LocalEventBacklog())
	grpcApiServer := grpc.NewServer(parameters.GrpcApiListenPort, parameters.GrpcApiListenAddress, gameApiHandler)

	eventPump := NewEventPump(blockchain.EventEmitter(), game)

	return &App{
		parameters:          parameters,
		game:                game,
		grpcApiServer:       grpcApiServer,
		blockchainConnector: blockchainConnector,
		blockchain:          blockchain,
		eventPump:           eventPump,
	}
}

func (a *App) Start(ctx context.Context) error {
	if err := a.game.Start(ctx); err != nil {
		return errors.Wrap(err, "error while starting game")
	}

	if err := a.blockchain.Start(ctx); err != nil {
		return errors.Wrap(err, "error while starting blockchain")

	}

	if err := a.grpcApiServer.Start(ctx); err != nil {
		return errors.Wrap(err, "error while starting game gRPC API server")
	}

	if err := a.eventPump.Start(ctx); err != nil {
		return errors.Wrap(err, "error while starting event pump")
	}

	return nil
}
