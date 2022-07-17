package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mayerkv/otus_go_homework/hw12_13_14_15_calendar/internal/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var ErrInternalError = status.Error(codes.Internal, "internal server error")

func ErrorInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if re := recover(); re != nil {
			err = ErrInternalError
		}
	}()

	resp, err = handler(ctx, req)
	code := status.Code(err)

	if errors.Is(err, app.ErrDateBusy) {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if errors.Is(err, app.ErrEventNotExists) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	if code == codes.Unknown || code == codes.Internal {
		return resp, ErrInternalError
	}

	return resp, err
}

func LoggingInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		p, _ := peer.FromContext(ctx)
		now := time.Now()
		resp, err = handler(ctx, req)

		msg := fmt.Sprintf(
			"%s %s %s %s",
			p.Addr,
			info.FullMethod,
			status.Code(err),
			time.Since(now),
		)
		logger.Info(msg)

		return resp, err
	}
}
