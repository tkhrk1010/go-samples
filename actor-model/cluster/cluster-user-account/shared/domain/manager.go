package domain

type Manager struct {
	AccountMap map[string]bool
}

func (c *Manager) Init() {
	c.AccountMap = map[string]bool{}
}

func (c *Manager) RegisterAccount(account string) {
	c.AccountMap[account] = true
}

func (c *Manager) DeregisterAccount(account string) {
	delete(c.AccountMap, account)
}
