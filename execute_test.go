// package tradespread

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"os"
// 	"reflect"
// 	"strconv"
// 	"testing"
// )

// var qOne = []QueueElement{ // avg price 14.657
// 	CreateQueueElement(14.50, 500),
// 	CreateQueueElement(14.6, 250),
// 	CreateQueueElement(14.65, 500),
// 	CreateQueueElement(14.8, 250),
// 	CreateQueueElement(14.9, 250),
// }

// type Testpositions struct {
// 	Testpositions []TestOptionPosition
// }

// type TestOptionPosition struct {
// 	Position         OptionPosition
// 	Test_correct_val float64
// }

// type OptionPosition struct {
// 	Type     ActionType
// 	Price    float64
// 	Name     string
// 	Quantity float64
// 	Queue    Queue
// }

// func (p OptionPosition) GetInstrumentName() string {
// 	return p.Name
// }

// func (p OptionPosition) GetPositionType() ActionType {
// 	return p.Type
// }

// func (p OptionPosition) GetAveragePrice() float64 {
// 	return p.Price
// }

// func (p OptionPosition) GetQuantity() float64 {
// 	return p.Quantity
// }

// func (p OptionPosition) GetOppositePositionType() ActionType {
// 	if p.Type == Buy {
// 		return Sell
// 	} else {
// 		return Buy
// 	}
// }

// func (p OptionPosition) GetQueue() Queue {
// 	return p.Queue
// }

// func CreateQueueElement(price, quantity float64) QueueElement {
// 	return QueueElement{price, quantity}
// }

// func CreateQueue(qType QueueType, queueElements []QueueElement) Queue {
// 	return Queue{
// 		Type:          qType,
// 		QueueElements: queueElements,
// 	}
// }

// func CreateDerivativePosition(pType ActionType, avgPrice float64, quantity float64, name string) DerivativePostion {
// 	return OptionPosition{
// 		Name:     name,
// 		Type:     pType,
// 		Price:    avgPrice,
// 		Quantity: quantity,
// 	}
// }

// func CreateQueueAveragePosition(pos DerivativePostion, averageExecutablePrice float64) QueueAveragePosition {
// 	return QueueAveragePosition{
// 		Position:               pos,
// 		AverageExecutablePrice: averageExecutablePrice,
// 	}
// }

// func readBytesFromFile(fileLocation string) []byte {
// 	jsonFile, err := os.Open(fileLocation)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer jsonFile.Close()
// 	byteValue, _ := ioutil.ReadAll(jsonFile)
// 	return byteValue
// }

// func setPositionsFromJson(testpositions *Testpositions, fileLocation string) {
// 	byteValue := readBytesFromFile(fileLocation)
// 	json.Unmarshal(byteValue, testpositions)
// }

// func withPrecision(val float64, precision int) string {
// 	ret := strconv.FormatFloat(val, 'f', precision, 64)
// 	return ret
// }

// func TestQueueAvgPosition_GetPNL(t *testing.T) {
// 	var testpositions Testpositions
// 	setPositionsFromJson(&testpositions, "test_positions_getpnl.json")

// 	var dPos DerivativePostion

// 	for _, pos := range testpositions.Testpositions {
// 		dPos = pos.Position
// 		avgPos := CreateQueueAveragePosition(dPos, GetQueueAveragePrice(&dPos))
// 		PRECISION = 2
// 		left := withPrecision(avgPos.GetPNL(), 2)
// 		right := withPrecision(pos.Test_correct_val, 2)
// 		if left != right {
// 			t.Errorf("GetPNL() returned %s want %s", left, right)
// 		}
// 	}
// }

// func TestSetQueueAveragePositions(t *testing.T) {
// 	type customtype struct {
// 		Derivativepositions []OptionPosition
// 		Solutions           []struct {
// 			Position               OptionPosition
// 			AverageExecutablePrice float64
// 		}
// 	}

// 	tests := customtype{}

// 	json.Unmarshal(readBytesFromFile("test_positions_setqavgpositions.json"), &tests)

// 	var arrayD []DerivativePostion

// 	for _, p := range tests.Derivativepositions {
// 		arrayD = append(arrayD, p)
// 	}

// 	var toSetQAvg []QueueAveragePosition
// 	func(d []DerivativePostion) {
// 		setQueueAveragePositions(&toSetQAvg, &d)
// 	}(arrayD)

// 	for i, qavg := range toSetQAvg {
// 		fmt.Println(qavg.Position)

// 		if withPrecision(qavg.AverageExecutablePrice, 2) != withPrecision(tests.Solutions[i].AverageExecutablePrice, 2) {
// 			t.Errorf("%d AverageExecutablePrice mismatch: want %s got %s", i, withPrecision(qavg.AverageExecutablePrice, 2), withPrecision(tests.Solutions[i].AverageExecutablePrice, 2))
// 		}
// 		if reflect.DeepEqual(qavg.Position, tests.Solutions[i].Position) == false {
// 			t.Errorf("Testcase %d: Deep reflect mismatch", i)
// 		}
// 	}
// }
