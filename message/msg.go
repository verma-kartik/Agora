package message

import (
	"time"
)

type Msg struct {
	Body       *[]byte
	Head       *map[string]string
	ReceivedAt time.Time
}

func NewMsg(givenHead *map[string]string, givenBody *[]byte) *Msg {
	return &Msg{Head: givenHead, Body: givenBody, ReceivedAt: time.Now()}
}

func NewHeaderlessMsg(givenBody *[]byte) *Msg {
	emptyHeader := make(map[string]string)
	return &Msg{Head: &emptyHeader, Body: givenBody, ReceivedAt: time.Now()}
}
