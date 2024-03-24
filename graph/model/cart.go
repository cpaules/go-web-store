package model

import "math"

type Cart struct {
	ID    string
	Items []*Item
}

func (c *Cart) Total() *float64 {
	var total float64
	for _, item := range c.Items {
		total += item.Price
	}
	total = math.Round(total*100) / 100
	return &total
}
