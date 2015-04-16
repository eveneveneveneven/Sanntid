package types

import (
	"bytes"
	"encoding/gob"
)

func Clone(dst, src interface{}) {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	enc.Encode(src)
	dec.Decode(dst)
}

func NewElevStat() *ElevStat {
	return &ElevStat{
		Dir:            0,
		Floor:          -1,
		InternalOrders: []bool{false, false, false, false},
	}
}

func NewNetworkMessage() *NetworkMessage {
	return &NetworkMessage{
		Id:       -1,
		Statuses: nil,
		Orders:   make(map[Order]struct{}),
	}
}
