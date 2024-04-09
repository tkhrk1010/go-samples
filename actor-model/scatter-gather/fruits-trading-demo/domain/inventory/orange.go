package inventory

type Orange struct {}

func (orange *Orange) GetAveragePurchasingPrice() float32 {
	return 1.0
}

func (orange *Orange) GetCount() int {
	return 1
}