package tradespread

type OptionPosition struct {
	Index    int
	Type     ActionType
	Price    float64
	Name     string
	Quantity float64
}

func (p OptionPosition) GetInstrumentName() string {
	return p.Name
}

func (p OptionPosition) GetPositionType() ActionType {
	return p.Type
}

func (p OptionPosition) GetAveragePrice() float64 {
	return p.Price
}

func (p OptionPosition) GetQuantity() float64 {
	return p.Quantity
}

func (p OptionPosition) GetOppositePositionType() ActionType {
	if p.Type == Buy {
		return Sell
	} else {
		return Buy
	}
}

func (p OptionPosition) GetQueue() Queue {
	queue1 := Queue{
		Type: Bid,
		QueueElements: []QueueElement{
			QueueElement{10.0, 1.0},
			QueueElement{9.0, 3.0},
			QueueElement{8.0, 4.0},
			QueueElement{7.0, 2.0},
			QueueElement{6.0, 5.0},
			QueueElement{5.0, 5.0},
		},
	}
	return queue1
}

// func TestGetOrders(t *testing.T) {
// 	leg := InputeLeg{
// 		Positions: []DerivativePostion{
// 			OptionPosition{Buy, 123.0, "A", 12.0},
// 			OptionPosition{Buy, 124.0, "B", 12.0},
// 			OptionPosition{Buy, 125.0, "C", 12.0},
// 		},
// 	}
// 	t.Logf("%s", GetOrders(leg))
// }

func CreateQueueElement(price, quantity float64) QueueElement {
	return QueueElement{price, quantity}
}

func CreateQueue(qType QueueType, queueElements []QueueElement) Queue {
	return Queue{
		Type:          qType,
		QueueElements: queueElements,
	}
}

func CreateDerivativePosition(pType ActionType, avgPrice float64, quantity float64, name string) DerivativePostion {
	return OptionPosition{
		Name:     name,
		Type:     pType,
		Price:    avgPrice,
		Quantity: quantity,
	}
}

// IndexToQueue := make(map[int]Queue)

var qOne = []QueueElement{ // avg price
	CreateQueueElement(14.50, 500),
	CreateQueueElement(14.6, 250),
	CreateQueueElement(14.65, 500),
	CreateQueueElement(14.8, 250),
	CreateQueueElement(14.9, 250),
}

// var avgOfHardcodedQueue = make(map[Queue]float64{qOne: 5}, 0)

func TestQueueAvgPositionPNL() {

}
