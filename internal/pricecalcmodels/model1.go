package pricecalcmodels

type MyModel struct {}

func NewMyModel() *MyModel{
	return &MyModel{}
}

func (m *MyModel)GetFairPrice(InputData *[]float64) float64{
	return 0.0
}