package domain

type Account struct {
	Name  string
	Count int64
}

func (i *Account) Add(a int64) {
	i.Count = i.Count + a
}

func (i *Account) Remove(a int64) {
	i.Count = i.Count - a
}
