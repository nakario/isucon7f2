package main

import (
	"math/big"
)

type mItem struct {
	ItemID int   `db:"item_id"`
	Power1 int64 `db:"power1"`
	Power2 int64 `db:"power2"`
	Power3 int64 `db:"power3"`
	Power4 int64 `db:"power4"`
	Price1 int64 `db:"price1"`
	Price2 int64 `db:"price2"`
	Price3 int64 `db:"price3"`
	Price4 int64 `db:"price4"`
}

var itemLists []mItem = []mItem{
	// ItemID:0は存在しないアイテムID
	mItem{ItemID: 0, Power1: 0, Power2: 1, Power3: 0, Power4: 1, Price1: 0, Price2: 1, Price3: 1, Price4: 1},
	mItem{ItemID: 1, Power1: 0, Power2: 1, Power3: 0, Power4: 1, Price1: 0, Price2: 1, Price3: 1, Price4: 1},
	mItem{ItemID: 2, Power1: 0, Power2: 1, Power3: 1, Power4: 1, Price1: 0, Price2: 1, Price3: 2, Price4: 1},
	mItem{ItemID: 3, Power1: 1, Power2: 10, Power3: 0, Power4: 2, Price1: 1, Price2: 3, Price3: 1, Price4: 2},
	mItem{ItemID: 4, Power1: 1, Power2: 24, Power3: 1, Power4: 2, Price1: 1, Price2: 10, Price3: 0, Price4: 3},
	mItem{ItemID: 5, Power1: 1, Power2: 25, Power3: 100, Power4: 3, Price1: 2, Price2: 20, Price3: 20, Price4: 2},
	mItem{ItemID: 6, Power1: 1, Power2: 30, Power3: 147, Power4: 13, Price1: 1, Price2: 22, Price3: 69, Price4: 17},
	mItem{ItemID: 7, Power1: 5, Power2: 80, Power3: 128, Power4: 6, Price1: 6, Price2: 61, Price3: 200, Price4: 5},
	mItem{ItemID: 8, Power1: 20, Power2: 340, Power3: 180, Power4: 3, Price1: 9, Price2: 105, Price3: 134, Price4: 14},
	mItem{ItemID: 9, Power1: 55, Power2: 520, Power3: 335, Power4: 5, Price1: 48, Price2: 243, Price3: 600, Price4: 7},
	mItem{ItemID: 10, Power1: 157, Power2: 1071, Power3: 1700, Power4: 12, Price1: 157, Price2: 625, Price3: 1000, Price4: 13},
	mItem{ItemID: 11, Power1: 2000, Power2: 7500, Power3: 2600, Power4: 3, Price1: 2001, Price2: 5430, Price3: 1000, Price4: 3},
	mItem{ItemID: 12, Power1: 1000, Power2: 9000, Power3: 0, Power4: 17, Price1: 963, Price2: 7689, Price3: 1, Price4: 19},
	mItem{ItemID: 13, Power1: 11000, Power2: 11000, Power3: 11000, Power4: 23, Price1: 10000, Price2: 2, Price3: 2, Price4: 29},
}

func exp4(a, b, c, d, x int64) *big.Int {
	// power(x):=(cx+1)*d^(ax+b)
	s := big.NewInt(c*x + 1)
	t := new(big.Int).Exp(big.NewInt(d), big.NewInt(a*x+b), nil)
	return new(big.Int).Mul(s, t)
}

const bufferNum = 50

var inittedBuffer = false

func initPowerBuffer() [][]*big.Int {
	result := make([][]*big.Int, len(itemLists))
	for i := range result {
		if i == 0 {
			continue
		}
		result[i] = make([]*big.Int, bufferNum)
		item := itemLists[i]
		for j := range result[i] {
			result[i][j] = exp4(item.Power1, item.Power2, item.Power3, item.Power4, int64(j))
		}
	}
	return result
}
func initPriceBuffer() [][]*big.Int {
	result := make([][]*big.Int, len(itemLists))
	for i := range result {
		if i == 0 {
			continue
		}
		result[i] = make([]*big.Int, bufferNum)
		item := itemLists[i]
		for j := range result[i] {
			result[i][j] = exp4(item.Price1, item.Price2, item.Price3, item.Price4, int64(j))
		}
	}
	inittedBuffer = true
	return result
}

var powerBuffer = initPowerBuffer()
var priceBuffer = initPriceBuffer()

func (item *mItem) GetPower(count int) *big.Int {
	if inittedBuffer && count < bufferNum {
		return powerBuffer[item.ItemID][count]
	}
	return exp4(item.Power1, item.Power2, item.Power3, item.Power4, int64(count))
}

func (item *mItem) GetPrice(count int) *big.Int {
	if inittedBuffer && count < bufferNum {
		return priceBuffer[item.ItemID][count]
	}
	return exp4(item.Price1, item.Price2, item.Price3, item.Price4, int64(count))
}
