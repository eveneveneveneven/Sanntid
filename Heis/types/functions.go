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
<<<<<<< HEAD
		InternalOrders: []int{-1, -1, -1, -1},
=======
		InternalOrders: []bool{false, false, false, false},
>>>>>>> 23d545ecd2d0f08b3f30f1b40de27a25060d8f4c
	}
}

func NewNetworkMessage() *NetworkMessage {
	return &NetworkMessage{
		Id:       -1,
		Statuses: nil,
		Orders:   make(map[Order]struct{}),
	}
}
