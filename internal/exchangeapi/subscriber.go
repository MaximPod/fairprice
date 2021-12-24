package exchangeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
)

// FairPriceCalcModel interface with method of calc fair price
type FairPriceCalcModel interface {
	GetFairPrice(InputData map[int]float64) float64
}

// FairPriceCalculator is the main entity app. It is collecting data and calculating faire price
type FairPriceCalculator struct {
	sync.RWMutex
	ExchangeDataChan <-chan *message.Message
	BarData          map[int]float64    // basket for data collection by sourceID
	StartBarDataTime time.Time          // time since backet collect data
	Model            FairPriceCalcModel // model type
	CalcInterval     int                // interval in seconds to change basket and calc FP
}

// NewFairPriceCalculator - constructor FairPriceCalculator
func NewFairPriceCalculator(exchangeDataChan <-chan *message.Message,
	calcInterval int, model FairPriceCalcModel) *FairPriceCalculator {
	return &FairPriceCalculator{
		ExchangeDataChan: exchangeDataChan,
		BarData:          make(map[int]float64),
		StartBarDataTime: time.Now(),
		Model:            model,
		CalcInterval:     calcInterval,
	}
}

// SubscribeRun listen channel and call HandleMessage func
// nolint:wrapcheck // demo mode
func (c *FairPriceCalculator) SubscribeRun(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-c.ExchangeDataChan:
			c.HandleMessage(msg)
			msg.Ack()
		}
	}
}

// HandleMessage - parse, validate and save data to basket
// nolint:gomnd // demo-mode
func (c *FairPriceCalculator) HandleMessage(msg *message.Message) {
	// parse
	var data TickerPrice

	err := json.Unmarshal(msg.Payload, &data)
	if err != nil {
		return
	}

	// validate
	if c.StartBarDataTime.Sub(data.Time) > 0 {
		return
	}

	price, err := strconv.ParseFloat(data.Price, 64)
	if err != nil {
		return
	}

	// save
	c.Lock()
	c.BarData[data.SourceID] = price
	c.Unlock()
}

// FairPriceCalcRun - runner winch HandleBarData every CalcInterval
// nolint:wrapcheck // demo mode
func (c *FairPriceCalculator) FairPriceCalcRun(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(c.CalcInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			c.HandleBarData()
		}
	}
}

// HandleBarData call model to calculate FP and show result
// nolint:forbidigo // demo-code
func (c *FairPriceCalculator) HandleBarData() {
	bardata := c.BarData
	startTime := c.StartBarDataTime

	c.Lock()
	c.BarData = make(map[int]float64)
	c.StartBarDataTime = time.Now()
	c.Unlock()

	res := c.Model.GetFairPrice(bardata)
	timestamp := startTime.Unix()

	fmt.Println(timestamp, fmt.Sprintf("%.2f", res))
}
