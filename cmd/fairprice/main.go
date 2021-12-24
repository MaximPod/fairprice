package main

import (
	"context"
	"fmt"

	"fairprice/internal/eventsource"
	"fairprice/internal/exchangeapi"
	"fairprice/internal/pricecalcmodels"
	apptools "fairprice/internal/tools/app"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

const (
	topicName    = "btc_usd_course.topic"
	calcInterval = 60 // second, min 5s
)

func main() {
	// context
	ctx := context.Background()

	// one of many message transports, see https://watermill.io
	pubSub := gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(false, false),
	)

	// subscriber is a data collector
	messages, err := pubSub.Subscribe(context.Background(), topicName)
	if err != nil {
		panic(err)
	}

	// model implement secret formula fair price calculation
	model := pricecalcmodels.NewMyModel()

	fairPriceCalculator := exchangeapi.NewFairPriceCalculator(messages, calcInterval, model)

	// publisher is a mock message generator
	messageGen := eventsource.NewMessageGen(pubSub, topicName, calcInterval)

	fmt.Printf("application started with interval %d second\n", calcInterval)

	// application
	apptools.RunParallel(ctx,
		apptools.SignalNotify,
		fairPriceCalculator.SubscribeRun,
		fairPriceCalculator.FairPriceCalcRun,
		messageGen.PublisherMockRun,
	)
}
