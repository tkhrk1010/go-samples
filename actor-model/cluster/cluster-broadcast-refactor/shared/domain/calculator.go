package domain

import (

)

type Calculator struct {
	Total int64
}

func (c *Calculator) Add(a int64) {
	c.Total = c.Total + a
}

func (c *Calculator) Subtract(a int64) {
	c.Total = c.Total - a
}
