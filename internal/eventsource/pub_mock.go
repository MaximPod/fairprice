package eventsource

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

// Payload implement data from exchange source
type Payload struct {
	SourceID int       `json:"source_id"` // id exchange source
	Ticker   string    `json:"ticker"`
	Time     time.Time `json:"time"`
	Price    string    `json:"price"`
}

// MessageGen generate messages to transport channel
type MessageGen struct {
	publisher      message.Publisher
	topicName      string
	currentPeriod  int
	periodInterval int
	errLog         *log.Logger
}

// NewMessageGen is a constructor MessageGen
func NewMessageGen(publisher message.Publisher, topicName string, periodInterval int) *MessageGen {
	errLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	m := &MessageGen{
		publisher:      publisher,
		topicName:      topicName,
		currentPeriod:  1,
		periodInterval: periodInterval,
		errLog:         errLog,
	}

	return m
}

// PublisherMockRun - runner witch public new messages
// nolint:wrapcheck // demo-code
func (m *MessageGen) PublisherMockRun(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(m.periodInterval) * time.Second)
	defer ticker.Stop()

	m.GenerateMessages()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			m.GenerateMessages()
		}
	}
}

// GenerateMessages call PublisherMock and switch currentPeriod
// nolint:gomnd //demo-code
func (m *MessageGen) GenerateMessages() {
	time.Sleep(2000) // network delay emulator
	m.PublisherMock()
	m.currentPeriod++

	if m.currentPeriod > maxPeriodCount {
		m.currentPeriod = 1
	}
}

// PublisherMock generate messages for maxChannelCount for currentPeriod
// nolint:gomnd // demo-code
func (m *MessageGen) PublisherMock() {
	for i := 1; i <= maxChannelCount; i++ {
		time.Sleep(200) // network delay emulator
		msgTime := time.Now()

		price := dataExample1[m.currentPeriod][i]

		payload := Payload{
			SourceID: i,
			Ticker:   "BTC_USD",
			Time:     msgTime,
			Price:    fmt.Sprintf("%d", price),
		}

		messagePayload, err := json.Marshal(payload)
		if err != nil {
			m.errLog.Println("eventsource.PublisherMock.json.Marshal", err)
		}

		msg := message.NewMessage(watermill.NewUUID(), messagePayload)

		err = m.publisher.Publish(m.topicName, msg)
		if err != nil {
			m.errLog.Println("eventsource.PublisherMock.Publish", err)
		}
	}
}
