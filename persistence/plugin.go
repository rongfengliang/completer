package persistence
/**
   This is derived from vendor/github.com/AsynkronIT/protoactor-go/persistence/plugin.go
   This has been modified to support propagating event indices to plugins
 */
import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/golang/protobuf/proto"
	"github.com/AsynkronIT/protoactor-go/persistence"
)

type persistent interface {
	init(provider Provider, context actor.Context)
	PersistReceive(message proto.Message)
	PersistSnapshot(snapshot proto.Message)
	Recovering() bool
	Name() string
}

type Mixin struct {
	eventIndex    int
	providerState ProviderState
	name          string
	receiver      receiver
	recovering    bool
}

func (mixin *Mixin) Recovering() bool {
	return mixin.recovering
}

func (mixin *Mixin) Name() string {
	return mixin.name
}
func (mixin *Mixin) PersistReceive(message proto.Message) {
	mixin.providerState.PersistEvent(mixin.Name(), mixin.eventIndex, message)
	mixin.eventIndex++
	if mixin.eventIndex%mixin.providerState.GetSnapshotInterval() == 0 {
		mixin.receiver.Receive(&persistence.RequestSnapshot{})
	}
}

func (mixin *Mixin) PersistSnapshot(snapshot proto.Message) {
	mixin.providerState.PersistSnapshot(mixin.Name(), mixin.eventIndex, snapshot)
}

func (mixin *Mixin) init(provider Provider, context actor.Context) {
	if mixin.providerState == nil {
		mixin.providerState = provider.GetState()
	}

	receiver := context.(receiver)

	mixin.name = context.Self().Id
	mixin.eventIndex = 0
	mixin.receiver = receiver
	mixin.recovering = true

	mixin.providerState.Restart()
	if snapshot, eventIndex, ok := mixin.providerState.GetSnapshot(mixin.Name()); ok {
		mixin.eventIndex = eventIndex
		receiver.Receive(snapshot)
	}
	mixin.providerState.GetEvents(mixin.Name(), mixin.eventIndex, func(index int, e interface{}) {
		receiver.Receive(e)
		mixin.eventIndex = index + 1
	})
	mixin.recovering = false
	receiver.Receive(&persistence.ReplayComplete{})
}

type receiver interface {
	Receive(message interface{})
}
