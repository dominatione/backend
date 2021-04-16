package blockchain

import (
	"context"
	"github.com/dominati-one/backend/internal/pkg/security"
	blockchainProtocol "github.com/dominati-one/backend/pkg/protocol/blockchain"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

type NetworkSettings struct {
	BlockInterval       time.Duration
	AuthorityPublicKeys *security.PublicKeysBag
	GenesisBlock        *blockchainProtocol.Block
}

type Network struct {
	log                            zerolog.Logger
	settings                       NetworkSettings
	connector                      Connector
	eventEmitter                   *EventEmitter
	localEventBacklog              *LocalEventBacklog
	networkEventBacklog            *NetworkEventBacklog
	eventValidator                 *EventValidator
	eventStorage                   EventStorage
	blockTicker                    *BlockTicker
	blockStorage                   BlockStorage
	blockValidator                 *BlockValidator
	localBlockBacklog              *LocalBlockBacklog
	blockBlockchainBacklogReceiver *BlockBlockchainBacklogReceiver
}

func NewNetwork(settings NetworkSettings, connector Connector, eventStorage EventStorage, blockStorage BlockStorage) *Network {
	eventValidator := NewEventValidator()
	blockValidator := NewBlockValidator(eventValidator)

	localBlockBacklog := NewLocalBlockBacklog(blockValidator)
	localEventBacklog := NewLocalEventBacklog(eventValidator)
	networkEventBacklog := NewNetworkEventBacklog(eventValidator)

	eventEmitter := NewEventEmitter()

	blockchainBlockBacklogReceiver := NewBlockBlockchainBacklogReceiver(blockBlockchainBacklogReceiverDependencies{
		connector:           connector,
		blockStorage:        blockStorage,
		eventEmitter:        eventEmitter,
		blockValidator:      blockValidator,
		localBlockBacklog:   localBlockBacklog,
		localEventBacklog:   localEventBacklog,
		networkEventBacklog: networkEventBacklog,
	})

	return &Network{
		log:                            log.With().Str("applicationComponent", "blockchain").Logger(),
		connector:                      connector,
		settings:                       settings,
		eventStorage:                   eventStorage,
		localEventBacklog:              localEventBacklog,
		networkEventBacklog:            networkEventBacklog,
		eventValidator:                 eventValidator,
		eventEmitter:                   eventEmitter,
		localBlockBacklog:              localBlockBacklog,
		blockTicker:                    NewBlockTicker(settings),
		blockStorage:                   blockStorage,
		blockValidator:                 blockValidator,
		blockBlockchainBacklogReceiver: blockchainBlockBacklogReceiver,
	}
}

func (n *Network) Start(ctx context.Context) error {
	if n.blockStorage.Count() == 0 {
		if blockId, err := n.blockStorage.Add(n.settings.GenesisBlock); err != nil {
			return errors.Wrap(err, "error while starting blockchain network")
		} else {
			n.log.Info().Str("blockId", blockId.String()).Msg("Initialized block storage with Genesis block.")
		}
	}

	go n.localEventBacklogSendLoop(ctx)
	go n.blockchainEventBacklogReceiveLoop(ctx)
	go n.blockBuildLoop(ctx)
	go n.localBlockBacklogSendLoop(ctx)

	if err := n.blockBlockchainBacklogReceiver.Start(ctx); err != nil {
		return errors.Wrap(err, "unable to start blockchain block backlog receiver")
	}

	n.log.Info().Msg("Started.")

	return nil
}

func (n *Network) LocalEventBacklog() *LocalEventBacklog {
	return n.localEventBacklog
}

func (n *Network) EventEmitter() *EventEmitter {
	return n.eventEmitter
}

func (n *Network) localEventBacklogSendLoop(ctx context.Context) {
	for range time.NewTicker(time.Millisecond * 100).C {
		for eventId, event := range n.localEventBacklog.Unsent() {
			log := n.log.With().
				Str("eventData", event.String()).
				Str("eventId", eventId.String()).
				Logger()

			err := n.connector.SendEventToBacklog(event)
			if err != nil {
				log.Warn().Err(err).Msg("Unable to send event from local backlog to blockchain backlog.")
				continue
			}

			err = n.localEventBacklog.MarkAsSent(eventId)
			if err != nil {
				log.Warn().Err(err).Msg("Unable to mark event as sent.")
				continue
			}

			log.Debug().Msg("Event sent from local backlog to blockchain backlog.")
		}
	}

}

func (n *Network) blockchainEventBacklogReceiveLoop(ctx context.Context) {
	for range time.NewTicker(time.Millisecond * 100).C {
		blockchainEvent, err := n.connector.GetBacklogEvent(ctx)
		if err != nil {
			log.Warn().Err(err).Msg("Unable to read event from blockchain backlog.")
			continue
		}
		log := log.With().Str("eventData", blockchainEvent.String()).Logger()

		eventId, err := NewEventId(blockchainEvent)
		if err != nil {
			log.Warn().Err(err).Msg("Unable to generate event id from data.")
			continue
		}
		log = log.With().Str("eventId", eventId.String()).Logger()

		if n.eventStorage.Exists(eventId) {
			log.Warn().Msg("Event received from blockchain backlog already exists in event storage.")
			continue
		}

		if n.localEventBacklog.Exists(eventId) {
			if err := n.localEventBacklog.MarkAsReceived(eventId); err != nil {
				log.Warn().Err(err).Msg("Unable to mark event from blockchain backlog as received in local backlog.")
				continue
			}
		}

		err = n.networkEventBacklog.Add(blockchainEvent)
		if err != nil {
			log.Warn().Err(err).Msg("Unable to add event to network backlog.")
			continue
		}

		log.Debug().Msg("Event received from blockchain backlog.")
	}
}

func (n *Network) blockBuildLoop(ctx context.Context) {
	for {
		lastBlock, err := n.blockStorage.GetLatestBlock()
		if err != nil {
			log.Warn().Err(err).Msg("Error while reading last block from storage.")
			continue
		}

		lastBlockTimestamp := CreateBlockTimestampFromUnixMilliseconds(lastBlock.Body.Timestamp)

		blockTimestamp, err := n.blockTicker.WaitForNext(ctx, lastBlockTimestamp)
		if err != nil {
			log.Warn().Err(err).Msg("Error while calculating next block timestamp.")
			continue
		}

		events := []*blockchainProtocol.Event{}

		for _, event := range n.networkEventBacklog.Unconfirmed() {
			events = append(events, event)
		}

		blockBuilder := NewBlockBuilder(lastBlock, events)

		newBlock, err := blockBuilder.Build(*blockTimestamp)
		if err != nil {
			log.Warn().Err(err).Msg("Unable to build block.")
			continue
		}

		blockId, err := n.localBlockBacklog.Add(*newBlock)
		if err != nil {
			log.Warn().Err(err).Msg("Unable to add block to local backlog.")
			continue
		}

		backOffDuration := n.settings.BlockInterval - time.Second

		log.Info().
			Str("blockId", blockId.String()).
			Uint64("blockTimestamp", blockTimestamp.UnixMilliseconds()).
			Dur("backOffDuration", backOffDuration).
			Msg("Block added to local backlog. Backing off.")

		time.Sleep(backOffDuration)
	}
}

func (n *Network) localBlockBacklogSendLoop(ctx context.Context) {
	for range time.NewTicker(time.Millisecond * 100).C {
		for blockId, block := range n.localBlockBacklog.Unsent() {
			log := n.log.With().
				Str("blockId", blockId.String()).
				Logger()

			err := n.connector.SendBlockToBacklog(block)
			if err != nil {
				log.Warn().Err(err).Msg("Unable to send block from local backlog to blockchain backlog.")
				continue
			}

			err = n.localBlockBacklog.MarkAsSent(blockId)
			if err != nil {
				log.Warn().Err(err).Msg("Unable to mark block as sent.")
				continue
			}

			log.Debug().Msg("Block sent from local backlog to blockchain backlog.")
		}
	}
}
