package domain

type Cart struct {
	ItemMap map[string]bool
}

func (c *Cart) Init() {
	c.ItemMap = map[string]bool{}
}

func (c *Cart) RegisterItem(item string) {
	c.ItemMap[item] = true
}

func (c *Cart) DeregisterItem(item string) {
	delete(c.ItemMap, item)
}