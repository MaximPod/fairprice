package eventsource

// maxPeriodCount - counts of source channels in data example
const maxChannelCount = 5

// maxPeriodCount - counts of periods in data example
const maxPeriodCount = 5

// dataExample1 is matrix, rows - periods, columns - source channels
// 0 index is exclude
// nolint:gochecknoglobals //demo-code
var dataExample1 = [][]int{
	{0, 0, 0, 0, 0, 0},
	{0, 10, 10, 10, 10, 10},
	{0, 20, 20, 20, 20, 20},
	{0, 30, 30, 30, 30, 30},
	{0, 20, 20, 20, 20, 20},
	{0, 10, 10, 10, 10, 10},
}
