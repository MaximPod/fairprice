package pricecalcmodels

type MyModel struct {
	historyPoints []float64
}

func NewMyModel() *MyModel {
	return &MyModel{}
}

// GetFairPrice calc fair price with three points approximation
// nolint:gomnd //demo-code
func (m *MyModel) GetFairPrice(inputData map[int]float64) float64 {
	summ := 0.0
	count := 0

	for _, v := range inputData {
		summ += v
		count++
	}

	if count == 0 {
		return 0.0
	}

	avg := summ / float64(count)

	if len(m.historyPoints) == 0 {
		m.historyPoints = append(m.historyPoints, avg, avg, avg) // start case
	} else {
		m.historyPoints = m.historyPoints[1:]
		m.historyPoints = append(m.historyPoints, avg)
	}

	summXY := 2*m.historyPoints[0] + 3*m.historyPoints[1] + 4*m.historyPoints[2]
	summY := m.historyPoints[0] + m.historyPoints[1] + m.historyPoints[2]
	a := 0.5*summXY - 1.5*summY
	b := (summY - 9*a) / 3
	res := 5*a + b

	return res
}
