package queue

import "github.com/verma-kartik/Agora/message"

type Queue struct {
	Name          string
	InputChannel  chan *message.Msg
	OutputChannel chan *message.Msg
	head          *queuemsg
	tail          *queuemsg
	length        int
}

func NewQueue(name string) *Queue {
	q := Queue{}
	q.InputChannel = make(chan *message.Msg)
	q.OutputChannel = make(chan *message.Msg)

	go q.pump()

	return &q
}

func (q *Queue) pump() {
pump:
	for {
		// If we have no messages - block until we receive one
		if q.head == nil {
			newMsg, ok := <-q.InputChannel
			newQueueMsg := newQueueMsg(newMsg)

			// Someone closed our input channel
			if !ok {
				break pump
			}
			q.head = newQueueMsg
			q.tail = newQueueMsg
			q.length++
		}

		select {

		case newMsg, ok := <-q.InputChannel:
			if !ok {
				break pump
			}
			newQueueMsg := newQueueMsg(newMsg)
			q.tail.next = newQueueMsg
			q.tail = newQueueMsg
			q.length++

		case q.OutputChannel <- q.head.data:
			q.head = q.head.next
			q.length--
		}
	}

	nextMsgToSend := q.head
	for nextMsgToSend != nil {
		q.OutputChannel <- nextMsgToSend.data
		nextMsgToSend = nextMsgToSend.next
	}

	close(q.OutputChannel)
}

func (q *Queue) PendingMessages() int {
	return q.length
}
