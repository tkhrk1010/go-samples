package market

type Apple struct{
	Price float32
}

func (apple *Apple) GetPrice() float32 {
	return apple.Price
}

