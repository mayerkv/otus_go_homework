package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/mayerkv/otus_go_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/mayerkv/otus_go_homework/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Server struct {
	pb.UnimplementedEventsServer
	host   string
	port   string
	lsn    net.Listener
	srv    *grpc.Server
	app    *app.App
	logger Logger
}

func NewServer(host string, port string, logger Logger, app *app.App) *Server {
	return &Server{host: host, port: port, app: app, logger: logger}
}

func (s *Server) Start(ctx context.Context) error {
	lsn, err := net.Listen("tcp", net.JoinHostPort(s.host, s.port))
	if err != nil {
		return fmt.Errorf("start Server: %w", err)
	}

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			LoggingInterceptor(s.logger),
			ErrorInterceptor,
		),
	)
	pb.RegisterEventsServer(srv, s)

	s.srv = srv
	s.lsn = lsn

	s.logger.Info("grpc Server started at " + lsn.Addr().String())
	if err := srv.Serve(lsn); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.srv.GracefulStop()
	return nil
}

func (s *Server) CreateEvent(
	ctx context.Context,
	request *pb.CreateEventRequest,
) (*pb.CreateEventResponse, error) {
	eventID, err := s.app.CreateEvent(
		ctx,
		request.Event.Title,
		request.Event.Description,
		request.Event.OwnerId,
		request.Event.StartAt.AsTime(),
		request.Event.EndAt.AsTime(),
		request.Event.NotifyBefore.AsDuration(),
	)
	if err != nil {
		return nil, err
	}

	return &pb.CreateEventResponse{EventId: eventID.String()}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, request *pb.UpdateEventRequest) (*emptypb.Empty, error) {
	err := s.app.UpdateEvent(
		ctx,
		request.EventId,
		request.Event.Title,
		request.Event.Description,
		request.Event.OwnerId,
		request.Event.StartAt.AsTime(),
		request.Event.EndAt.AsTime(),
		request.Event.NotifyBefore.AsDuration(),
	)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, request *pb.DeleteEventRequest) (*emptypb.Empty, error) {
	err := s.app.DeleteEvent(ctx, request.EventId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) GetEvents(
	ctx context.Context,
	request *pb.GetEventsRequest,
) (*pb.GetEventsResponse, error) {
	list, err := s.app.GetEventList(ctx, request.OwnerId, request.StartAt.AsTime(), request.EndAt.AsTime())
	if err != nil {
		return nil, err
	}

	events := make([]*pb.Event, 0, len(list))

	for _, item := range list {
		events = append(
			events, &pb.Event{
				Id:           item.ID.String(),
				Title:        item.Title,
				StartAt:      timestamppb.New(item.StartAt),
				EndAt:        timestamppb.New(item.EndAt),
				Description:  item.Description,
				OwnerId:      item.OwnerID.String(),
				NotifyBefore: durationpb.New(item.NotifyBefore),
			},
		)
	}

	return &pb.GetEventsResponse{Events: events}, nil
}
