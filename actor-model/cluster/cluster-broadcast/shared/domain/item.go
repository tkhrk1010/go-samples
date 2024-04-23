package domain

type Item struct {
	Name  string
	Count int64
}

func (i *Item) Add(a int64) {
	i.Count = i.Count + a
}

func (i *Item) Remove(a int64) {
	i.Count = i.Count - a
}
