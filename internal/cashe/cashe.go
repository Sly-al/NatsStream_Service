package cashe

import (
	"NatsStream_Service/internal/model"
	"errors"
)

type Cashe struct {
	Data map[string]model.Order_client
}

func NewCashe() Cashe {
	return Cashe{Data: make(map[string]model.Order_client)}
}

func (c *Cashe) InsertToCashe(jsonToInsert model.Order_client) error {
	_, ok := c.Data[jsonToInsert.Order_uid]
	if ok {
		return errors.New("already in cashe")
	}
	c.Data[jsonToInsert.Order_uid] = jsonToInsert
	return nil
}

func (c *Cashe) GetFromCashe(order_uid string) (model.Order_client, error) {
	order, ok := c.Data[order_uid]
	if !ok {
		return model.Order_client{}, errors.New("no such order_uid in cashe")
	}
	return order, nil
}
