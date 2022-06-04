package tradespread

import (
	"errors"
	"time"
)

var PRECISION int

type ActionType int

const (
	Buy ActionType = iota
	Sell
)

type QStatus int

const (
	NOT_ENOUGH_QUANTITY QStatus = iota
	ENOUGH_QUANTITY
)

type DerivativePostion interface {
	GetInstrumentName() string
	GetPositionType() ActionType
	GetOppositePositionType() ActionType
	GetAveragePrice() float64
	GetQuantity() float64
	GetQueue() Queue
}

type InputeLeg struct {
	Positions []*DerivativePostion
	TargetPNL float64
}

type Orders struct {
	Positions []QueueAveragePosition
}

type QueueAveragePosition struct {
	Position               DerivativePostion
	AverageExecutablePrice float64 // Average price in queue
}

func (q *QueueAveragePosition) GetPNL() float64 {
	if q.Position.GetPositionType() == Buy {
		return (q.AverageExecutablePrice - q.Position.GetAveragePrice()) * (q.Position.GetQuantity())
	} else {
		return (q.Position.GetAveragePrice() - q.AverageExecutablePrice) * (-1 * q.Position.GetQuantity())
	}
}

func totalQuantityExceedsPosition(currentQuantity, positionQuantity float64) bool {
	return (currentQuantity >= positionQuantity)
}

func GetQueueAveragePrice(position *DerivativePostion) (float64, QStatus) {
	queue := (*position).GetQueue()
	positionQuantity := (*position).GetQuantity()
	sumOfProducts := 0.0
	sumOfQuantity := 0.0
	for _, quote := range queue.QueueElements {
		price := quote.Price
		quantity := quote.Quantity
		if totalQuantityExceedsPosition(sumOfQuantity+quantity, positionQuantity) {
			remainingQuantity := positionQuantity - sumOfQuantity
			sumOfQuantity += remainingQuantity
			sumOfProducts += price * remainingQuantity
			break
		}
		sumOfQuantity += quantity
		sumOfProducts += price * quantity
	}
	averagePrice := sumOfProducts / sumOfQuantity

	if sumOfQuantity < positionQuantity {
		return averagePrice, NOT_ENOUGH_QUANTITY
	}

	return averagePrice, ENOUGH_QUANTITY
}

func setQueueAveragePositions(positions []*DerivativePostion) ([]QueueAveragePosition, error) {

	queueAveragePositions := make([]QueueAveragePosition, 0)

	for _, position := range positions {
		avgPrice, status := GetQueueAveragePrice(position)
		if status == NOT_ENOUGH_QUANTITY {
			empty := []QueueAveragePosition{}
			return empty, errors.New("NOT_ENOUGH_QUANTITY")
		}
		newAveragePosition := QueueAveragePosition{
			Position:               *position,
			AverageExecutablePrice: avgPrice,
		}
		queueAveragePositions = append(queueAveragePositions, newAveragePosition)
	}
	return queueAveragePositions, nil
}

func PNLSumofQueueAveragePositions(queueAveragePositions *[]QueueAveragePosition) float64 {
	PNLSum := 0.0
	for _, position := range *queueAveragePositions {
		PNLSum += position.GetPNL()
	}
	return PNLSum
}

func GetQueueAveragePositions(positions []*DerivativePostion) ([]QueueAveragePosition, error) {
	queueAveragePositions, err := setQueueAveragePositions(positions)
	return queueAveragePositions, err
}

func GetOrders(inputleg InputeLeg, queueUpdateInterval time.Duration) Orders {

	queueAveragePositions, _ := setQueueAveragePositions(inputleg.Positions)
	// Loop tille target is achievable
	for PNLSumofQueueAveragePositions(&queueAveragePositions) < inputleg.TargetPNL {
		time.Sleep(queueUpdateInterval)
		queueAveragePositions, _ = setQueueAveragePositions(inputleg.Positions)
	}

	var orders Orders
	orders.Positions = queueAveragePositions
	return orders
}
