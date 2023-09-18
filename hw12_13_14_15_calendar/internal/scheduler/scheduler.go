package scheduler

import (
	"os"
	"time"

	"github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/logger"
	rpcclient "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/server/grpc/client"
	rpcs "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/server/grpc/server"
	"github.com/mailru/easyjson"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/net/context"
)

type Scheduler struct {
	cfg  Cfg
	conn *amqp.Connection
	ch   *amqp.Channel
	err  error
	ctx  context.Context
}

func NewScheduler(configuration Cfg) *Scheduler {
	return &Scheduler{cfg: configuration}
}

func confirmOne(confirms <-chan amqp.Confirmation, log *logger.Logger) {
	if confirmed := <-confirms; confirmed.Ack {
		log.Info("Confirmation delivery")
	} else {
		log.Error("Confirmation delivery error")
	}
}

func (s *Scheduler) Start(ctx context.Context) error { //nolint: gocognit
	s.ctx = ctx
	s.conn, s.err = amqp.Dial(s.cfg.Target.ConnectionString)
	log := logger.New(s.cfg.Logger, os.Stdout)
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
		s.cfg.Target.ExchangeName,
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
	if s.err = s.ch.Confirm(false); s.err != nil {
		log.Error("Checking comnfirmation failed")
	}
	confirms := s.ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	// Create rpc client connect and get events list
	src := rpcclient.Client{}
	src.Connect(s.cfg.Source.ConnectionString)
	// Cycling in default case and checking unscheduled events every timeout value
	// Then mark events as scheduled and publish in rqbbit
FORCYCLE:
	for {
		select {
		case <-ctx.Done():
			break FORCYCLE
		case <-time.After(time.Second * time.Duration(s.cfg.Timeout)):
			rpcevents, err := src.NotScheduledEvents(ctx)
			if err != nil { //nolint: nestif
				log.Error("Error calling rpc method gets unscheduled events")
			} else {
				for _, rpcevent := range rpcevents {
					event := rpcs.ConvertToEvent(*rpcevent)
					val := time.Now().Add(time.Duration(event.NotificationTime) * time.Minute)
					if val.Compare(event.Start) >= 0 {
						// Mark as scheduled and publish message
						event.Scheduled = true
						*rpcevent = rpcs.ConvertFromEvent(event)
						e, err := src.UpdateEvent(s.ctx, rpcevent)
						_ = e
						if err != nil {
							log.Error("Error calling rpc method updates uscheduled event")
						}
						body, err := easyjson.Marshal(event)
						if err != nil {
							log.Error("Error while marshaling event")
						}
						s.err = s.ch.PublishWithContext(
							s.ctx,
							s.cfg.Target.ExchangeName,
							s.cfg.Target.Key,
							false,
							false,
							amqp.Publishing{
								Headers:         amqp.Table{},
								ContentType:     "text/plain",
								ContentEncoding: "",
								Body:            body,
								DeliveryMode:    amqp.Transient,
								Priority:        0,
							},
						)
						if s.err != nil {
							log.Error("Error while publishing event")
						}
						confirmOne(confirms, log)
					}
				}
			}
		}
	}
	return nil
}

func (s *Scheduler) Stop() error {
	return s.conn.Close()
}
