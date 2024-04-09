package market

type Orange struct {
	Price float32
}

func (orange *Orange) GetPrice() float32 {
	return orange.Price
}