package blockchain

import (
	"context"
	blockchainProtocol "github.com/dominati-one/backend/pkg/protocol/blockchain"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

type blockBlockchainBacklogReceiverDependencies struct {
	connector           Connector
	blockStorage        BlockStorage
	blockValidator      *BlockValidator
	localBlockBacklog   *LocalBlockBacklog
	localEventBacklog   *LocalEventBacklog
	networkEventBacklog *NetworkEventBacklog
	eventEmitter        *EventEmitter
}

type BlockBlockchainBacklogReceiver struct {
	blockBlockchainBacklogReceiverDependencies
	log zerolog.Logger
}

func NewBlockBlockchainBacklogReceiver(dependencies blockBlockchainBacklogReceiverDependencies) *BlockBlockchainBacklogReceiver {
	return &BlockBlockchainBacklogReceiver{
		log: log.With().Str("applicationComponent", "blockchain").Str("blockchainComponent", "blockBlockchainBacklogReceiver").Logger(),
		blockBlockchainBacklogReceiverDependencies: dependencies,
	}
}

func (r *BlockBlockchainBacklogReceiver) Start(ctx context.Context) error {
	go func() {
		for range time.NewTicker(time.Millisecond * 100).C {
			r.loop(ctx)
		}
	}()
	return nil
}

func (r *BlockBlockchainBacklogReceiver) loop(ctx context.Context) {
	blockchainBlock, err := r.connector.GetBacklogBlock(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("Unable to read block from blockchain backlog.")
		return
	}

	blockId, err := NewBlockId(blockchainBlock)
	if err != nil {
		log.Warn().Err(err).Msg("Unable to generate block id from data.")
		return
	}
	log := r.log.With().Str("blockId", blockId.String()).Logger()

	if r.blockStorage.Exists(*blockId) {
		log.Warn().Msg("Block received from blockchain backlog already exists in block storage.")
		return
	}

	latestBlock, err := r.blockStorage.GetLatestBlock()
	if err != nil {
		log.Warn().Err(err).Msg("Unable to read latest block from storage.")
		return
	}

	err = r.blockValidator.Validate(latestBlock, blockchainBlock)
	if err != nil {
		log.Warn().Err(err).Msg("Unable to validate block from blockchain backlog.")
		return
	}

	err = r.eventEmitter.emitBlock(blockchainBlock)
	if err != nil {
		log.Warn().Err(err).Msg("Unable to emit block.")
		return
	}

	for _, blockEvent := range blockchainBlock.Body.Events {
		if err := r.processEvent(blockEvent.Event); err != nil {
			log.Panic().Err(err).Msg("Unable to process event. Inconsistency detected.")
			return
		}
	}

	_, err = r.blockStorage.Add(blockchainBlock)
	if err != nil {
		log.Warn().Err(err).Msg("Unable to add block from blockchain backlog to storage.")
		return
	}

	if r.localBlockBacklog.Exists(*blockId) {
		if err := r.localBlockBacklog.MarkAsConfirmed(*blockId); err != nil {
			log.Warn().Err(err).Msg("Unable to mark block from blockchain backlog as received in local backlog.")
			return
		}
	}

	log.Debug().Msg("Block received from blockchain backlog and add added to storage. Events processed without errors.")

}

func (r *BlockBlockchainBacklogReceiver) processEvent(blockEvent *blockchainProtocol.Event) error {
	eventId, err := NewEventId(blockEvent)
	if err != nil {
		return errors.Wrap(err, "unable to generate event id")
	}

	if r.localEventBacklog.Exists(eventId) {
		if err := r.localEventBacklog.MarkAsConfirmed(eventId); err != nil {
			return errors.Wrap(err, "unable to confirm event in local backlog")
		}
	}

	if r.networkEventBacklog.Exists(eventId) {
		if err := r.networkEventBacklog.MarkAsConfirmed(eventId); err != nil {
			return errors.Wrap(err, "unable to confirm event in network backlog")
		}
	}

	if err := r.eventEmitter.emitEvent(blockEvent); err != nil {
		return errors.Wrap(err, "error while event emission")
	}

	return nil
}
