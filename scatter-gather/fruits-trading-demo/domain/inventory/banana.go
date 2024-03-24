package inventory

type Banana struct {}

func (banana *Banana) GetAveragePurchasingPrice() float32 {
	return 1.0
}

func (banana *Banana) GetCount() int {
	return 1
}