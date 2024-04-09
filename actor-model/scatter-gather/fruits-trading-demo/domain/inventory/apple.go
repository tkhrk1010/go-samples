package inventory

type Apple struct {}

func (apple *Apple) GetAveragePurchasingPrice() float32 {
	return 1.0
}

func (apple *Apple) GetCount() int {
	return 1
}