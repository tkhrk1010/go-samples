// Code generated by protoc-gen-grain. DO NOT EDIT.
// versions:
//  protoc-gen-grain v0.7.0
//  protoc           v5.26.1
// source: protos.proto

package proto

import (
	fmt "fmt"
	actor "github.com/asynkron/protoactor-go/actor"
	cluster "github.com/asynkron/protoactor-go/cluster"
	proto "google.golang.org/protobuf/proto"
	slog "log/slog"
	time "time"
)

var xItemFactory func() Item

// ItemFactory produces a Item
func ItemFactory(factory func() Item) {
	xItemFactory = factory
}

// GetItemGrainClient instantiates a new ItemGrainClient with given Identity
func GetItemGrainClient(c *cluster.Cluster, id string) *ItemGrainClient {
	if c == nil {
		panic(fmt.Errorf("nil cluster instance"))
	}
	if id == "" {
		panic(fmt.Errorf("empty id"))
	}
	return &ItemGrainClient{Identity: id, cluster: c}
}

// GetItemKind instantiates a new cluster.Kind for Item
func GetItemKind(opts ...actor.PropsOption) *cluster.Kind {
	props := actor.PropsFromProducer(func() actor.Actor {
		return &ItemActor{
			Timeout: 60 * time.Second,
		}
	}, opts...)
	kind := cluster.NewKind("Item", props)
	return kind
}

// GetItemKind instantiates a new cluster.Kind for Item
func NewItemKind(factory func() Item, timeout time.Duration, opts ...actor.PropsOption) *cluster.Kind {
	xItemFactory = factory
	props := actor.PropsFromProducer(func() actor.Actor {
		return &ItemActor{
			Timeout: timeout,
		}
	}, opts...)
	kind := cluster.NewKind("Item", props)
	return kind
}

// Item interfaces the services available to the Item
type Item interface {
	Init(ctx cluster.GrainContext)
	Terminate(ctx cluster.GrainContext)
	ReceiveDefault(ctx cluster.GrainContext)
	Add(req *NumberRequest, ctx cluster.GrainContext) (*CountResponse, error)
	Remove(req *NumberRequest, ctx cluster.GrainContext) (*CountResponse, error)
	GetCurrent(req *Noop, ctx cluster.GrainContext) (*CountResponse, error)
}

// ItemGrainClient holds the base data for the ItemGrain
type ItemGrainClient struct {
	Identity string
	cluster  *cluster.Cluster
}

// Add requests the execution on to the cluster with CallOptions
func (g *ItemGrainClient) Add(r *NumberRequest, opts ...cluster.GrainCallOption) (*CountResponse, error) {
	if g.cluster.Config.RequestLog {
		g.cluster.Logger().Info("Requesting", slog.String("identity", g.Identity), slog.String("kind", "Item"), slog.String("method", "Add"), slog.Any("request", r))
	}
	bytes, err := proto.Marshal(r)
	if err != nil {
		return nil, err
	}
	reqMsg := &cluster.GrainRequest{MethodIndex: 0, MessageData: bytes}
	resp, err := g.cluster.Request(g.Identity, "Item", reqMsg, opts...)
	if err != nil {
		return nil, fmt.Errorf("error request: %w", err)
	}
	switch msg := resp.(type) {
	case *CountResponse:
		return msg, nil
	case *cluster.GrainErrorResponse:
		if msg == nil {
			return nil, nil
		}
		return nil, msg
	default:
		return nil, fmt.Errorf("unknown response type %T", resp)
	}
}

// Remove requests the execution on to the cluster with CallOptions
func (g *ItemGrainClient) Remove(r *NumberRequest, opts ...cluster.GrainCallOption) (*CountResponse, error) {
	if g.cluster.Config.RequestLog {
		g.cluster.Logger().Info("Requesting", slog.String("identity", g.Identity), slog.String("kind", "Item"), slog.String("method", "Remove"), slog.Any("request", r))
	}
	bytes, err := proto.Marshal(r)
	if err != nil {
		return nil, err
	}
	reqMsg := &cluster.GrainRequest{MethodIndex: 1, MessageData: bytes}
	resp, err := g.cluster.Request(g.Identity, "Item", reqMsg, opts...)
	if err != nil {
		return nil, fmt.Errorf("error request: %w", err)
	}
	switch msg := resp.(type) {
	case *CountResponse:
		return msg, nil
	case *cluster.GrainErrorResponse:
		if msg == nil {
			return nil, nil
		}
		return nil, msg
	default:
		return nil, fmt.Errorf("unknown response type %T", resp)
	}
}

// GetCurrent requests the execution on to the cluster with CallOptions
func (g *ItemGrainClient) GetCurrent(r *Noop, opts ...cluster.GrainCallOption) (*CountResponse, error) {
	if g.cluster.Config.RequestLog {
		g.cluster.Logger().Info("Requesting", slog.String("identity", g.Identity), slog.String("kind", "Item"), slog.String("method", "GetCurrent"), slog.Any("request", r))
	}
	bytes, err := proto.Marshal(r)
	if err != nil {
		return nil, err
	}
	reqMsg := &cluster.GrainRequest{MethodIndex: 2, MessageData: bytes}
	resp, err := g.cluster.Request(g.Identity, "Item", reqMsg, opts...)
	if err != nil {
		return nil, fmt.Errorf("error request: %w", err)
	}
	switch msg := resp.(type) {
	case *CountResponse:
		return msg, nil
	case *cluster.GrainErrorResponse:
		if msg == nil {
			return nil, nil
		}
		return nil, msg
	default:
		return nil, fmt.Errorf("unknown response type %T", resp)
	}
}

// ItemActor represents the actor structure
type ItemActor struct {
	ctx     cluster.GrainContext
	inner   Item
	Timeout time.Duration
}

// Receive ensures the lifecycle of the actor for the received message
func (a *ItemActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started: //pass
	case *cluster.ClusterInit:
		a.ctx = cluster.NewGrainContext(ctx, msg.Identity, msg.Cluster)
		a.inner = xItemFactory()
		a.inner.Init(a.ctx)

		if a.Timeout > 0 {
			ctx.SetReceiveTimeout(a.Timeout)
		}
	case *actor.ReceiveTimeout:
		ctx.Poison(ctx.Self())
	case *actor.Stopped:
		a.inner.Terminate(a.ctx)
	case actor.AutoReceiveMessage: // pass
	case actor.SystemMessage: // pass

	case *cluster.GrainRequest:
		switch msg.MethodIndex {
		case 0:
			req := &NumberRequest{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				ctx.Logger().Error("[Grain] Add(NumberRequest) proto.Unmarshal failed.", slog.Any("error", err))
				resp := cluster.NewGrainErrorResponse(cluster.ErrorReason_INVALID_ARGUMENT, err.Error()).
					WithMetadata(map[string]string{
						"argument": req.String(),
					})
				ctx.Respond(resp)
				return
			}

			r0, err := a.inner.Add(req, a.ctx)
			if err != nil {
				resp := cluster.FromError(err)
				ctx.Respond(resp)
				return
			}
			ctx.Respond(r0)
		case 1:
			req := &NumberRequest{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				ctx.Logger().Error("[Grain] Remove(NumberRequest) proto.Unmarshal failed.", slog.Any("error", err))
				resp := cluster.NewGrainErrorResponse(cluster.ErrorReason_INVALID_ARGUMENT, err.Error()).
					WithMetadata(map[string]string{
						"argument": req.String(),
					})
				ctx.Respond(resp)
				return
			}

			r0, err := a.inner.Remove(req, a.ctx)
			if err != nil {
				resp := cluster.FromError(err)
				ctx.Respond(resp)
				return
			}
			ctx.Respond(r0)
		case 2:
			req := &Noop{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				ctx.Logger().Error("[Grain] GetCurrent(Noop) proto.Unmarshal failed.", slog.Any("error", err))
				resp := cluster.NewGrainErrorResponse(cluster.ErrorReason_INVALID_ARGUMENT, err.Error()).
					WithMetadata(map[string]string{
						"argument": req.String(),
					})
				ctx.Respond(resp)
				return
			}

			r0, err := a.inner.GetCurrent(req, a.ctx)
			if err != nil {
				resp := cluster.FromError(err)
				ctx.Respond(resp)
				return
			}
			ctx.Respond(r0)
		}
	default:
		a.inner.ReceiveDefault(a.ctx)
	}
}

// onError should be used in ctx.ReenterAfter
// you can just return error in reenterable method for other errors
func (a *ItemActor) onError(err error) {
	resp := cluster.FromError(err)
	a.ctx.Respond(resp)
}

var xCartFactory func() Cart

// CartFactory produces a Cart
func CartFactory(factory func() Cart) {
	xCartFactory = factory
}

// GetCartGrainClient instantiates a new CartGrainClient with given Identity
func GetCartGrainClient(c *cluster.Cluster, id string) *CartGrainClient {
	if c == nil {
		panic(fmt.Errorf("nil cluster instance"))
	}
	if id == "" {
		panic(fmt.Errorf("empty id"))
	}
	return &CartGrainClient{Identity: id, cluster: c}
}

// GetCartKind instantiates a new cluster.Kind for Cart
func GetCartKind(opts ...actor.PropsOption) *cluster.Kind {
	props := actor.PropsFromProducer(func() actor.Actor {
		return &CartActor{
			Timeout: 60 * time.Second,
		}
	}, opts...)
	kind := cluster.NewKind("Cart", props)
	return kind
}

// GetCartKind instantiates a new cluster.Kind for Cart
func NewCartKind(factory func() Cart, timeout time.Duration, opts ...actor.PropsOption) *cluster.Kind {
	xCartFactory = factory
	props := actor.PropsFromProducer(func() actor.Actor {
		return &CartActor{
			Timeout: timeout,
		}
	}, opts...)
	kind := cluster.NewKind("Cart", props)
	return kind
}

// Cart interfaces the services available to the Cart
type Cart interface {
	Init(ctx cluster.GrainContext)
	Terminate(ctx cluster.GrainContext)
	ReceiveDefault(ctx cluster.GrainContext)
	RegisterGrain(req *RegisterMessage, ctx cluster.GrainContext) (*Noop, error)
	DeregisterGrain(req *RegisterMessage, ctx cluster.GrainContext) (*Noop, error)
	BroadcastGetCounts(req *Noop, ctx cluster.GrainContext) (*TotalsResponse, error)
}

// CartGrainClient holds the base data for the CartGrain
type CartGrainClient struct {
	Identity string
	cluster  *cluster.Cluster
}

// RegisterGrain requests the execution on to the cluster with CallOptions
func (g *CartGrainClient) RegisterGrain(r *RegisterMessage, opts ...cluster.GrainCallOption) (*Noop, error) {
	if g.cluster.Config.RequestLog {
		g.cluster.Logger().Info("Requesting", slog.String("identity", g.Identity), slog.String("kind", "Cart"), slog.String("method", "RegisterGrain"), slog.Any("request", r))
	}
	bytes, err := proto.Marshal(r)
	if err != nil {
		return nil, err
	}
	reqMsg := &cluster.GrainRequest{MethodIndex: 0, MessageData: bytes}
	resp, err := g.cluster.Request(g.Identity, "Cart", reqMsg, opts...)
	if err != nil {
		return nil, fmt.Errorf("error request: %w", err)
	}
	switch msg := resp.(type) {
	case *Noop:
		return msg, nil
	case *cluster.GrainErrorResponse:
		if msg == nil {
			return nil, nil
		}
		return nil, msg
	default:
		return nil, fmt.Errorf("unknown response type %T", resp)
	}
}

// DeregisterGrain requests the execution on to the cluster with CallOptions
func (g *CartGrainClient) DeregisterGrain(r *RegisterMessage, opts ...cluster.GrainCallOption) (*Noop, error) {
	if g.cluster.Config.RequestLog {
		g.cluster.Logger().Info("Requesting", slog.String("identity", g.Identity), slog.String("kind", "Cart"), slog.String("method", "DeregisterGrain"), slog.Any("request", r))
	}
	bytes, err := proto.Marshal(r)
	if err != nil {
		return nil, err
	}
	reqMsg := &cluster.GrainRequest{MethodIndex: 1, MessageData: bytes}
	resp, err := g.cluster.Request(g.Identity, "Cart", reqMsg, opts...)
	if err != nil {
		return nil, fmt.Errorf("error request: %w", err)
	}
	switch msg := resp.(type) {
	case *Noop:
		return msg, nil
	case *cluster.GrainErrorResponse:
		if msg == nil {
			return nil, nil
		}
		return nil, msg
	default:
		return nil, fmt.Errorf("unknown response type %T", resp)
	}
}

// BroadcastGetCounts requests the execution on to the cluster with CallOptions
func (g *CartGrainClient) BroadcastGetCounts(r *Noop, opts ...cluster.GrainCallOption) (*TotalsResponse, error) {
	if g.cluster.Config.RequestLog {
		g.cluster.Logger().Info("Requesting", slog.String("identity", g.Identity), slog.String("kind", "Cart"), slog.String("method", "BroadcastGetCounts"), slog.Any("request", r))
	}
	bytes, err := proto.Marshal(r)
	if err != nil {
		return nil, err
	}
	reqMsg := &cluster.GrainRequest{MethodIndex: 2, MessageData: bytes}
	resp, err := g.cluster.Request(g.Identity, "Cart", reqMsg, opts...)
	if err != nil {
		return nil, fmt.Errorf("error request: %w", err)
	}
	switch msg := resp.(type) {
	case *TotalsResponse:
		return msg, nil
	case *cluster.GrainErrorResponse:
		if msg == nil {
			return nil, nil
		}
		return nil, msg
	default:
		return nil, fmt.Errorf("unknown response type %T", resp)
	}
}

// CartActor represents the actor structure
type CartActor struct {
	ctx     cluster.GrainContext
	inner   Cart
	Timeout time.Duration
}

// Receive ensures the lifecycle of the actor for the received message
func (a *CartActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started: //pass
	case *cluster.ClusterInit:
		a.ctx = cluster.NewGrainContext(ctx, msg.Identity, msg.Cluster)
		a.inner = xCartFactory()
		a.inner.Init(a.ctx)

		if a.Timeout > 0 {
			ctx.SetReceiveTimeout(a.Timeout)
		}
	case *actor.ReceiveTimeout:
		ctx.Poison(ctx.Self())
	case *actor.Stopped:
		a.inner.Terminate(a.ctx)
	case actor.AutoReceiveMessage: // pass
	case actor.SystemMessage: // pass

	case *cluster.GrainRequest:
		switch msg.MethodIndex {
		case 0:
			req := &RegisterMessage{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				ctx.Logger().Error("[Grain] RegisterGrain(RegisterMessage) proto.Unmarshal failed.", slog.Any("error", err))
				resp := cluster.NewGrainErrorResponse(cluster.ErrorReason_INVALID_ARGUMENT, err.Error()).
					WithMetadata(map[string]string{
						"argument": req.String(),
					})
				ctx.Respond(resp)
				return
			}

			r0, err := a.inner.RegisterGrain(req, a.ctx)
			if err != nil {
				resp := cluster.FromError(err)
				ctx.Respond(resp)
				return
			}
			ctx.Respond(r0)
		case 1:
			req := &RegisterMessage{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				ctx.Logger().Error("[Grain] DeregisterGrain(RegisterMessage) proto.Unmarshal failed.", slog.Any("error", err))
				resp := cluster.NewGrainErrorResponse(cluster.ErrorReason_INVALID_ARGUMENT, err.Error()).
					WithMetadata(map[string]string{
						"argument": req.String(),
					})
				ctx.Respond(resp)
				return
			}

			r0, err := a.inner.DeregisterGrain(req, a.ctx)
			if err != nil {
				resp := cluster.FromError(err)
				ctx.Respond(resp)
				return
			}
			ctx.Respond(r0)
		case 2:
			req := &Noop{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				ctx.Logger().Error("[Grain] BroadcastGetCounts(Noop) proto.Unmarshal failed.", slog.Any("error", err))
				resp := cluster.NewGrainErrorResponse(cluster.ErrorReason_INVALID_ARGUMENT, err.Error()).
					WithMetadata(map[string]string{
						"argument": req.String(),
					})
				ctx.Respond(resp)
				return
			}

			r0, err := a.inner.BroadcastGetCounts(req, a.ctx)
			if err != nil {
				resp := cluster.FromError(err)
				ctx.Respond(resp)
				return
			}
			ctx.Respond(r0)
		}
	default:
		a.inner.ReceiveDefault(a.ctx)
	}
}

// onError should be used in ctx.ReenterAfter
// you can just return error in reenterable method for other errors
func (a *CartActor) onError(err error) {
	resp := cluster.FromError(err)
	a.ctx.Respond(resp)
}

func respond[T proto.Message](ctx cluster.GrainContext) func(T) {
	return func(resp T) {
		ctx.Respond(resp)
	}
}
