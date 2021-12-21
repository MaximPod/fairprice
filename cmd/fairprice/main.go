package main

import (
    "context"
  
    "github.com/ThreeDotsLabs/watermill"
    "github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	

	"fairprice/internal/exchangeapi"
	"fairprice/internal/eventsource"
	"fairprice/internal/pricecalcmodels"
	apptools "fairprice/internal/tools/app"
)

const (
	topicName = "btc_usd_course.topic"
)

func main() {
	// context
	ctx := context.Background()

	
    pubSub := gochannel.NewGoChannel(
        gochannel.Config{},
        watermill.NewStdLogger(false, false),
    )

	// subscriber
    messages, err := pubSub.Subscribe(context.Background(), topicName)
    if err != nil {
        panic(err)
    }


	model := pricecalcmodels.NewMyModel()

	fairPriceCalculator := exchangeapi.NewFairPriceCalculator(messages, 60, model)

	// publisher
	messageGen := eventsource.NewMessageGen(pubSub, topicName)


	// application
	apptools.RunParallel(ctx,
		apptools.SignalNotify,
		fairPriceCalculator.SubscribeRun,
		fairPriceCalculator.FairPriceCalcRun,
		messageGen.PublisherMockRun,
	)
}

