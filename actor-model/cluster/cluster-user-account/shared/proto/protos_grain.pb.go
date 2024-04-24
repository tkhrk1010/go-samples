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

var xManagerFactory func() Manager

// ManagerFactory produces a Manager
func ManagerFactory(factory func() Manager) {
	xManagerFactory = factory
}

// GetManagerGrainClient instantiates a new ManagerGrainClient with given Identity
func GetManagerGrainClient(c *cluster.Cluster, id string) *ManagerGrainClient {
	if c == nil {
		panic(fmt.Errorf("nil cluster instance"))
	}
	if id == "" {
		panic(fmt.Errorf("empty id"))
	}
	return &ManagerGrainClient{Identity: id, cluster: c}
}

// GetManagerKind instantiates a new cluster.Kind for Manager
func GetManagerKind(opts ...actor.PropsOption) *cluster.Kind {
	props := actor.PropsFromProducer(func() actor.Actor {
		return &ManagerActor{
			Timeout: 60 * time.Second,
		}
	}, opts...)
	kind := cluster.NewKind("Manager", props)
	return kind
}

// GetManagerKind instantiates a new cluster.Kind for Manager
func NewManagerKind(factory func() Manager, timeout time.Duration, opts ...actor.PropsOption) *cluster.Kind {
	xManagerFactory = factory
	props := actor.PropsFromProducer(func() actor.Actor {
		return &ManagerActor{
			Timeout: timeout,
		}
	}, opts...)
	kind := cluster.NewKind("Manager", props)
	return kind
}

// Manager interfaces the services available to the Manager
type Manager interface {
	Init(ctx cluster.GrainContext)
	Terminate(ctx cluster.GrainContext)
	ReceiveDefault(ctx cluster.GrainContext)
	GetAllAccountEmails(req *Noop, ctx cluster.GrainContext) (*EmailsResponse, error)
	CreateAccount(req *Noop, ctx cluster.GrainContext) (*AccountIdResponse, error)
	GetAccount(req *AccountIdResponse, ctx cluster.GrainContext) (*AccountResponse, error)
}

// ManagerGrainClient holds the base data for the ManagerGrain
type ManagerGrainClient struct {
	Identity string
	cluster  *cluster.Cluster
}

// GetAllAccountEmails requests the execution on to the cluster with CallOptions
func (g *ManagerGrainClient) GetAllAccountEmails(r *Noop, opts ...cluster.GrainCallOption) (*EmailsResponse, error) {
	if g.cluster.Config.RequestLog {
		g.cluster.Logger().Info("Requesting", slog.String("identity", g.Identity), slog.String("kind", "Manager"), slog.String("method", "GetAllAccountEmails"), slog.Any("request", r))
	}
	bytes, err := proto.Marshal(r)
	if err != nil {
		return nil, err
	}
	reqMsg := &cluster.GrainRequest{MethodIndex: 0, MessageData: bytes}
	resp, err := g.cluster.Request(g.Identity, "Manager", reqMsg, opts...)
	if err != nil {
		return nil, fmt.Errorf("error request: %w", err)
	}
	switch msg := resp.(type) {
	case *EmailsResponse:
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

// CreateAccount requests the execution on to the cluster with CallOptions
func (g *ManagerGrainClient) CreateAccount(r *Noop, opts ...cluster.GrainCallOption) (*AccountIdResponse, error) {
	if g.cluster.Config.RequestLog {
		g.cluster.Logger().Info("Requesting", slog.String("identity", g.Identity), slog.String("kind", "Manager"), slog.String("method", "CreateAccount"), slog.Any("request", r))
	}
	bytes, err := proto.Marshal(r)
	if err != nil {
		return nil, err
	}
	reqMsg := &cluster.GrainRequest{MethodIndex: 1, MessageData: bytes}
	resp, err := g.cluster.Request(g.Identity, "Manager", reqMsg, opts...)
	if err != nil {
		return nil, fmt.Errorf("error request: %w", err)
	}
	switch msg := resp.(type) {
	case *AccountIdResponse:
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

// GetAccount requests the execution on to the cluster with CallOptions
func (g *ManagerGrainClient) GetAccount(r *AccountIdResponse, opts ...cluster.GrainCallOption) (*AccountResponse, error) {
	if g.cluster.Config.RequestLog {
		g.cluster.Logger().Info("Requesting", slog.String("identity", g.Identity), slog.String("kind", "Manager"), slog.String("method", "GetAccount"), slog.Any("request", r))
	}
	bytes, err := proto.Marshal(r)
	if err != nil {
		return nil, err
	}
	reqMsg := &cluster.GrainRequest{MethodIndex: 2, MessageData: bytes}
	resp, err := g.cluster.Request(g.Identity, "Manager", reqMsg, opts...)
	if err != nil {
		return nil, fmt.Errorf("error request: %w", err)
	}
	switch msg := resp.(type) {
	case *AccountResponse:
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

// ManagerActor represents the actor structure
type ManagerActor struct {
	ctx     cluster.GrainContext
	inner   Manager
	Timeout time.Duration
}

// Receive ensures the lifecycle of the actor for the received message
func (a *ManagerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started: //pass
	case *cluster.ClusterInit:
		a.ctx = cluster.NewGrainContext(ctx, msg.Identity, msg.Cluster)
		a.inner = xManagerFactory()
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
			req := &Noop{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				ctx.Logger().Error("[Grain] GetAllAccountEmails(Noop) proto.Unmarshal failed.", slog.Any("error", err))
				resp := cluster.NewGrainErrorResponse(cluster.ErrorReason_INVALID_ARGUMENT, err.Error()).
					WithMetadata(map[string]string{
						"argument": req.String(),
					})
				ctx.Respond(resp)
				return
			}

			r0, err := a.inner.GetAllAccountEmails(req, a.ctx)
			if err != nil {
				resp := cluster.FromError(err)
				ctx.Respond(resp)
				return
			}
			ctx.Respond(r0)
		case 1:
			req := &Noop{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				ctx.Logger().Error("[Grain] CreateAccount(Noop) proto.Unmarshal failed.", slog.Any("error", err))
				resp := cluster.NewGrainErrorResponse(cluster.ErrorReason_INVALID_ARGUMENT, err.Error()).
					WithMetadata(map[string]string{
						"argument": req.String(),
					})
				ctx.Respond(resp)
				return
			}

			r0, err := a.inner.CreateAccount(req, a.ctx)
			if err != nil {
				resp := cluster.FromError(err)
				ctx.Respond(resp)
				return
			}
			ctx.Respond(r0)
		case 2:
			req := &AccountIdResponse{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				ctx.Logger().Error("[Grain] GetAccount(AccountIdResponse) proto.Unmarshal failed.", slog.Any("error", err))
				resp := cluster.NewGrainErrorResponse(cluster.ErrorReason_INVALID_ARGUMENT, err.Error()).
					WithMetadata(map[string]string{
						"argument": req.String(),
					})
				ctx.Respond(resp)
				return
			}

			r0, err := a.inner.GetAccount(req, a.ctx)
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
func (a *ManagerActor) onError(err error) {
	resp := cluster.FromError(err)
	a.ctx.Respond(resp)
}

func respond[T proto.Message](ctx cluster.GrainContext) func(T) {
	return func(resp T) {
		ctx.Respond(resp)
	}
}
