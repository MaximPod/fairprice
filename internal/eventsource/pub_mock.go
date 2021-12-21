package eventsource

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type MessageGen struct {
	publisher message.Publisher
	topicName string
	errLog    *log.Logger
}

func NewMessageGen(publisher message.Publisher, topicName string) *MessageGen {
	errLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	return &MessageGen{
		publisher: publisher,
		topicName: topicName,
		errLog:    errLog,
	}
}

func (m *MessageGen) PublisherMockRun(ctx context.Context) error {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			m.PublisherMock()
		}
	}
}

const maxChannelCount = 5

func (m *MessageGen) PublisherMock() {
	for i := 1; i < maxChannelCount; i++ {
		msg := message.NewMessage(watermill.NewUUID(), []byte("Hello, world!"))

		if err := m.publisher.Publish(m.topicName, msg); err != nil {
			m.errLog.Println("eventsource.PublisherMock.Publish", err)
		}
	}
}
