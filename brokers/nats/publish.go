// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Package nats contains NATS message publisher implementation.
package nats

import (
	"context"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/mainflux/mainflux"
	"github.com/nats-io/nats.go"
)

// NatsPublisher specifies a message publishing API.
type NatsPublisher interface {
	// Publish publishes message to the msessage broker.
	Publish(context.Context, string, mainflux.Message) error

	PubConn() *nats.Conn
}

const prefix = "channel"

var _ NatsPublisher = (*natsPub)(nil)

type natsPub struct {
	conn *nats.Conn
}

// NewPublisher NATS message publisher.
func NewPublisher(url string) (NatsPublisher, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &natsPub{conn: nc}, nil
}

func (np natsPub) Publish(_ context.Context, _ string, msg mainflux.Message) error {
	data, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("%s.%s", prefix, msg.Channel)
	if msg.Subtopic != "" {
		subject = fmt.Sprintf("%s.%s", subject, msg.Subtopic)
	}
	return np.conn.Publish(subject, data)
}

func (np natsPub) PubConn() *nats.Conn {
	return np.conn
}