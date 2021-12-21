package exchangeapi

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
)

type FairPriceCalcModel interface {
	GetFairPrice(InputData *[]float64) float64
}

type FairPriceCalculator struct {
	ExchangeDataChan <-chan *message.Message
	BarData      *map[string]float64
	Model        FairPriceCalcModel
	CalcInterval int
	l            *log.Logger
}

func NewFairPriceCalculator(exchangeDataChan <-chan *message.Message, 
	calcInterval int, model FairPriceCalcModel) *FairPriceCalculator {
	bar := make(map[string]float64)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	return &FairPriceCalculator{
		ExchangeDataChan: exchangeDataChan,
		BarData:      &bar,
		Model:        model,
		CalcInterval: calcInterval,
		l:            infoLog,
	}
}


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

func (c *FairPriceCalculator) HandleMessage(msg *message.Message) {

	// пайлоад - строку - структуру
	// валидируем дату
	// заполняем мапу

}

func (c *FairPriceCalculator) FairPriceCalcRun(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(c.CalcInterval)*time.Second)
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

func (c *FairPriceCalculator)HandleBarData(){
	//  ссыль на мапу
			// блокировка мьютекса
			// ссыль на новую мапу
			// разблокировка мьютекса
			// мапу в массив в модель
			// вывод результата

		c.l.Println("HandleBarData")

}


