package market

type Banana struct {
	Price float32
}

func (banana *Banana) GetPrice() float32 {
	return banana.Price
}