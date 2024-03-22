package model

type Cart struct {
	ID    string
	Items []*Item
}

// func (c *Cart) IsDistributorToplineAdParticipant() bool {
// 	return hardcode.DistyToplineParticipants[c.ID]
// }
