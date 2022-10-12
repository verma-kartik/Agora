package queue

import "github.com/verma-kartik/Agora/message"

type queuemsg struct {
	next *queuemsg
	data *message.Msg
}

func newQueueMsg(giveMsg *message.Msg) *queuemsg {
	return &queuemsg{data: giveMsg}
}
