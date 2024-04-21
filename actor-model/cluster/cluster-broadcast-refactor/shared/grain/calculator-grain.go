package grain

import (
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-broadcast-refactor/shared/domain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-broadcast-refactor/shared/proto"
)

type CalcGrain struct {
	calculator domain.Calculator
}

func (c *CalcGrain) ReceiveDefault(ctx cluster.GrainContext) {
}

func (c *CalcGrain) Init(ctx cluster.GrainContext) {
	c.calculator.Total = 0

	// register with the tracker
	trackerGrain := proto.GetTrackerGrainClient(ctx.Cluster(), "singleTrackerGrain")
	trackerGrain.RegisterGrain(&proto.RegisterMessage{GrainId: ctx.Identity()})
}

func (c *CalcGrain) Terminate(ctx cluster.GrainContext) {
	// deregister with the tracker
	trackerGrain := proto.GetTrackerGrainClient(ctx.Cluster(), "singleTrackerGrain")
	trackerGrain.DeregisterGrain(&proto.RegisterMessage{GrainId: ctx.Identity()})
}

func (c *CalcGrain) Add(n *proto.NumberRequest, ctx cluster.GrainContext) (*proto.CountResponse, error) {
	c.calculator.Add(n.Number)
	return &proto.CountResponse{Number: c.calculator.Total}, nil
}

func (c *CalcGrain) Subtract(n *proto.NumberRequest, ctx cluster.GrainContext) (*proto.CountResponse, error) {
	c.calculator.Subtract(n.Number)
	return &proto.CountResponse{Number: c.calculator.Total}, nil
}

func (c *CalcGrain) GetCurrent(n *proto.Noop, ctx cluster.GrainContext) (*proto.CountResponse, error) {
	return &proto.CountResponse{Number: c.calculator.Total}, nil
}
