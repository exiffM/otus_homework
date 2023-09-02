package client

import (
	"context"
	"errors"
	"io"
	"log"

	eventrpcapi "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var localOpts = []grpc.CallOption{}

type Client struct {
	grpcClient eventrpcapi.EventServiceClient
	conn       *grpc.ClientConn
}

func (c *Client) Connect(dsn string) error {
	connection, err := grpc.Dial(dsn, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed on client: %v", err)
		return err
	}
	c.conn = connection
	c.grpcClient = eventrpcapi.NewEventServiceClient(c.conn)
	return nil
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateEvent(ctx context.Context, event *eventrpcapi.Event) (*eventrpcapi.Event, error) {
	return c.grpcClient.CreateEvent(ctx, event, localOpts...)
}

func (c *Client) SelectEvent(ctx context.Context, id *eventrpcapi.Id) (*eventrpcapi.Event, error) {
	return c.grpcClient.SelectEvent(ctx, id, localOpts...)
}

func (c *Client) UpdateEvent(ctx context.Context, event *eventrpcapi.Event) (*eventrpcapi.Event, error) {
	return c.grpcClient.UpdateEvent(ctx, event, localOpts...)
}

func (c *Client) DeleteEvent(ctx context.Context, id *eventrpcapi.Id) (*emptypb.Empty, error) {
	return c.grpcClient.DeleteEvent(ctx, id, localOpts...)
}

func (c *Client) Events(ctx context.Context) ([]*eventrpcapi.Event, error) {
	events, err := c.grpcClient.Events(ctx, &emptypb.Empty{}, localOpts...)
	if err != nil {
		return nil, err
	}
	pbEvents := make([]*eventrpcapi.Event, 0)
	for {
		pbEvent, err := events.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		pbEvents = append(pbEvents, pbEvent)
	}
	return pbEvents, nil
}
