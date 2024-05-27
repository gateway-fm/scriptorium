package queue_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/gateway-fm/scriptorium/clog"
	"github.com/gateway-fm/scriptorium/queue"
)

type EventBusSuite struct {
	suite.Suite

	ctx context.Context
	bus queue.EventBus
	log *clog.CustomLogger
}

func TestEventBusSuite(t *testing.T) {
	suite.Run(t, new(EventBusSuite))
}

func (s *EventBusSuite) SetupTest() {
	s.ctx = context.Background()
	s.log = clog.NewCustomLogger(os.Stdout, slog.LevelDebug, false)

	s.bus = queue.NewEventBus(s.ctx, 100)
	s.bus.SetLogger(s.log)
}

func (s *EventBusSuite) TestRetryLogic() {
	s.Run("test one topic", func() {
		retryDelays := []int{1, 1, 1}

		received := make(chan bool, 2)

		s.bus.Subscribe("topic", func(_ context.Context, event *queue.Event) queue.AckStatus {
			if event.Retry < 3 {
				return queue.NACK
			}

			received <- true

			return queue.ACK
		}, retryDelays, time.Second)

		s.bus.Publish("topic", []byte("Test Event"))
		go func() {
			err := s.bus.StartProcessing(s.ctx)
			s.Require().NoError(err)
		}()
		defer s.bus.Stop()

		select {
		case <-received:
			s.T().Log("Event processed successfully after retries")
		case <-time.After(4 * time.Second):
			s.FailNow("Event was not processed successfully within the expected time")
		}
	})

	s.Run("test multiple topics", func() {
		retryDelays := []int{1, 1, 1}
		received := make(chan bool, 2)

		s.bus.Subscribe("first-topic", func(_ context.Context, event *queue.Event) queue.AckStatus {
			if event.Retry < 3 {
				return queue.NACK
			}

			received <- true

			return queue.ACK
		}, retryDelays, time.Second)

		s.bus.Subscribe("second-topic", func(_ context.Context, event *queue.Event) queue.AckStatus {
			if event.Retry < 3 {
				return queue.NACK
			}

			received <- true

			return queue.ACK
		}, retryDelays, time.Second)

		s.bus.Publish("first-topic", []byte("Test Event"))
		s.bus.Publish("second-topic", []byte("Test Event"))

		go func() {
			err := s.bus.StartProcessing(s.ctx)
			s.Require().NoError(err)
		}()

		defer s.bus.Stop()

		select {
		case <-received:
			s.T().Log("Event processed successfully after retries")
		case <-time.After(4 * time.Second):
			s.FailNow("Event was not processed successfully within the expected time")
		}
	})
}
