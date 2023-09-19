package main

import (
	"context"
	"os"

	"github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/mailru/easyjson"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Sender struct {
	cfg  SenderCfg
	conn *amqp.Connection
	ch   *amqp.Channel
	que  amqp.Queue
	err  error
	ctx  context.Context
}

func NewSender(configuration SenderCfg) *Sender {
	return &Sender{cfg: configuration}
}

func (s *Sender) Start(ctx context.Context) error {
	s.ctx = ctx
	s.conn, s.err = amqp.Dial(s.cfg.Source.ConnectionString)
	log := logger.New(s.cfg.LoggLevel, os.Stdout)
	if s.err != nil {
		log.Error("Failed to connect to rabbitmq")
		return s.err
	}

	s.ch, s.err = s.conn.Channel()
	if s.err != nil {
		log.Error("Failed to open a channel")
		return s.err
	}

	s.err = s.ch.ExchangeDeclare(
		s.cfg.Source.ExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if s.err != nil {
		log.Error("Failed to declare an exchange")
		return s.err
	}

	s.que, s.err = s.ch.QueueDeclare(
		s.cfg.Source.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if s.err != nil {
		log.Error("Failed to declare a queue")
		return s.err
	}

	s.err = s.ch.QueueBind(
		s.cfg.Source.QueueName,
		s.cfg.Source.Key,
		s.cfg.Source.ExchangeName,
		false,
		nil,
	)
	if s.err != nil {
		log.Error("Failed to bind a queue")
		return s.err
	}

	deliveds, err := s.ch.Consume(
		s.cfg.Source.QueueName,
		s.cfg.Source.Tag,
		false,
		false,
		false,
		false,
		nil,
	)
	if s.err != nil {
		log.Error("chanel.Consume has failed")
		return err
	}

FORCYCLE:
	for {
		select {
		case <-s.ctx.Done():
			break FORCYCLE
		case message := <-deliveds:
			event := storage.Event{}
			err := easyjson.Unmarshal(message.Body, &event)
			if err != nil {
				log.Error("Error while unmarshaling message body to event")
			} else {
				log.Info("Notification!\n Event %q,\n %v, starts in %v",
					event.Title,
					event.Description,
					event.Start,
				)
			}
			message.Ack(false)
		}
	}

	return nil
}

func (s *Sender) Stop() error {
	return s.ch.Cancel(s.cfg.Source.Tag, true)
}
