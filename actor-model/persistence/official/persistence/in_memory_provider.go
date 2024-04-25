package provider

import (
	"github.com/asynkron/protoactor-go/persistence"
)

type Provider struct {
	providerState persistence.ProviderState
}

func NewProvider(snapshotInterval int) *Provider {
	return &Provider{
		providerState: persistence.NewInMemoryProvider(snapshotInterval),
	}
}

func (p *Provider) GetState() persistence.ProviderState {
	return p.providerState
}
