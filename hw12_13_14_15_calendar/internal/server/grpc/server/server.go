package rpcserver

import (
	"context"
	"net"

	interfaces "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/interface"
	eventrpcapi "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/server/grpc/pb"
	mdl "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Converts protobuf Event to application Event type.
func ConvertToEvent(e eventrpcapi.Event) mdl.Event {
	return mdl.Event{
		ID:               int(e.Id),
		Title:            e.Tittle,
		Start:            e.Start.AsTime(),
		Duration:         int(e.Duration),
		Description:      e.Description,
		NotificationTime: int(e.Notification),
		Scheduled:        e.Scheduled,
	}
}

// Converts application Event to protobuf Event type.
func ConvertFromEvent(e mdl.Event) eventrpcapi.Event {
	return eventrpcapi.Event{
		Id:           int32(e.ID),
		Tittle:       e.Title,
		Start:        timestamppb.New(e.Start),
		Duration:     int32(e.Duration),
		Description:  e.Description,
		Notification: int32(e.NotificationTime),
		Scheduled:    e.Scheduled,
	}
}

type GRPCServer struct {
	eventrpcapi.UnimplementedEventServiceServer
	logger interfaces.Logger
	app    interfaces.Application
	server *grpc.Server
	lis    net.Listener
}

func NewGRPCServer(
	l interfaces.Logger, a interfaces.Application) *GRPCServer { //nolint: gofumpt
	return &GRPCServer{logger: l, app: a}
}

func (s *GRPCServer) CreateEvent(ctx context.Context, e *eventrpcapi.Event) (*eventrpcapi.Event, error) {
	event := ConvertToEvent(*e)
	ev, err := s.app.CreateEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	*e = ConvertFromEvent(ev)
	return e, nil
}

func (s *GRPCServer) SelectEvent(ctx context.Context, id *eventrpcapi.Id) (*eventrpcapi.Event, error) {
	event, err := s.app.SelectEvent(ctx, int(id.Id))
	if err != nil {
		return nil, err
	}
	res := ConvertFromEvent(event)
	return &res, nil
}

func (s *GRPCServer) UpdateEvent(ctx context.Context, e *eventrpcapi.Event) (*eventrpcapi.Event, error) {
	event := ConvertToEvent(*e)
	ev, err := s.app.UpdateEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	*e = ConvertFromEvent(ev)
	return e, nil
}

func (s *GRPCServer) DeleteEvent(ctx context.Context, id *eventrpcapi.Id) (*emptypb.Empty, error) {
	if err := s.app.DeleteEvent(ctx, int(id.Id)); err != nil {
		return new(emptypb.Empty), err
	}
	return new(emptypb.Empty), nil
}

func (s *GRPCServer) Events(_ *emptypb.Empty, stream eventrpcapi.EventService_EventsServer) error {
	events, err := s.app.Events()
	if err != nil {
		return err
	}
	for _, event := range events {
		pbEvent := ConvertFromEvent(event)
		if err := stream.Send(&pbEvent); err != nil {
			return err
		}
	}
	return nil
}

func (s *GRPCServer) NotScheduledEvents(
	_ *emptypb.Empty,
	stream eventrpcapi.EventService_NotScheduledEventsServer,
) error {
	events, err := s.app.NotScheduledEvents()
	if err != nil {
		return err
	}
	for _, event := range events {
		pbEvent := ConvertFromEvent(event)
		if err := stream.Send(&pbEvent); err != nil {
			return err
		}
	}
	return nil
}

func (s *GRPCServer) Start(address string) error {
	var err error
	s.lis, err = net.Listen("tcp", address)
	if err != nil {
		return err
	}
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryInterceptor(s.logger)),
		grpc.StreamInterceptor(StreamInterceptor(s.logger)),
	)
	eventrpcapi.RegisterEventServiceServer(gRPCServer, s)
	s.server = gRPCServer
	s.logger.Info("GRPCServer.Start()")
	err = s.server.Serve(s.lis)
	return err
}

func (s *GRPCServer) Stop() {
	s.logger.Info("GRPCServer.Stop()")
	s.server.Stop()
	s.logger.Error("After RPC stop!")
}

func (s *GRPCServer) GracefulStop() {
	s.logger.Info("GRPCServer.GracefulStop()")
	// s.lis.Close()
	s.server.GracefulStop()
	s.logger.Error("After RPC Graceful stop!")
}

func UnaryInterceptor(logger interfaces.Logger) func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger.Warn("UnaryInterceptor: Method - %q Object - %s", info.FullMethod, req)
		return handler(ctx, req)
	}
}

func StreamInterceptor(logger interfaces.Logger) func(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		logger.Warn("StreamInterceptor: Method - %q", info.FullMethod)
		return handler(srv, stream)
	}
}
