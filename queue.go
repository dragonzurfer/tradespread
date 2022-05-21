package tradespread

type QueueType int

const (
	Bid QueueType = iota
	Offer
)

type QueueElement struct {
	Price    float64
	Quantity float64
}

type Queue struct {
	Type          QueueType
	QueueElements []QueueElement
}
