package tradespread

import "time"

type ActionType int

const (
	Buy ActionType = iota
	Sell
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
	Positions []DerivativePostion
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
		return (q.AverageExecutablePrice - q.Position.GetAveragePrice()) * q.Position.GetQuantity()
	} else {
		return (q.Position.GetAveragePrice() - q.AverageExecutablePrice) * q.Position.GetQuantity()
	}
}

func totalQuantityExceedsPosition(currentQuantity, positionQuantity float64) bool {
	return (currentQuantity >= positionQuantity)
}

func GetQueueAveragePrice(position *DerivativePostion) float64 {
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
			sumOfProducts += price * quantity
			break
		}
		sumOfQuantity += quantity
		sumOfProducts += price * quantity
	}
	averagePrice := sumOfProducts / sumOfQuantity
	return averagePrice
}

func setQueueAveragePositions(queueAveragePositions *[]QueueAveragePosition, positions *[]DerivativePostion) {
	for _, position := range *positions {
		avgPrice := GetQueueAveragePrice(&position)
		newAveragePosition := QueueAveragePosition{
			Position:               position,
			AverageExecutablePrice: avgPrice,
		}
		*queueAveragePositions = append(*queueAveragePositions, newAveragePosition)
	}
}

func PNLSumofQueueAveragePositions(queueAveragePositions *[]QueueAveragePosition) float64 {
	PNLSum := 0.0
	for _, position := range *queueAveragePositions {
		PNLSum += position.GetPNL()
	}
	return PNLSum
}

func GetOrders(inputleg InputeLeg, queueUpdateInterval time.Duration) Orders {

	// Get average prices
	var queueAveragePositions []QueueAveragePosition

	setQueueAveragePositions(&queueAveragePositions, &inputleg.Positions)

	// Loop tille target is achievable
	for PNLSumofQueueAveragePositions(&queueAveragePositions) < inputleg.TargetPNL {
		time.Sleep(queueUpdateInterval)
		setQueueAveragePositions(&queueAveragePositions, &inputleg.Positions)
	}

	var orders Orders
	orders.Positions = queueAveragePositions
	return orders
}
