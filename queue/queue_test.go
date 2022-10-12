package queue

import (
	"bytes"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/verma-kartik/Agora/message"
)

func TestQueue_CanSendAndReceiveBasicMessages(t *testing.T) {
	underTest := NewQueue("TestQueue")

	testMessagePayload := []byte("TestMessage")
	underTest.InputChannel <- message.NewHeaderlessMsg(&testMessagePayload)

	receivedMessage := <-underTest.OutputChannel

	if !bytes.Equal(*receivedMessage.Body, testMessagePayload) {
		t.Fail()
	}
}

func TestQueue_ReceiveBeforeSend_ReturnsExpectedResult(t *testing.T) {
	gomega.RegisterTestingT(t)

	underTest := NewQueue("TestQueue")

	var receivedMessage *message.Msg
	go func() {
		receivedMessage = <-underTest.OutputChannel
	}()

	time.Sleep(time.Millisecond * 10)

	testMessagePayload := []byte("TestMessage")
	underTest.InputChannel <- message.NewHeaderlessMsg(&testMessagePayload)

	gomega.Eventually(func() *message.Msg {
		return receivedMessage
	}).Should(gomega.Not(gomega.BeNil()))

	if !bytes.Equal(*receivedMessage.Body, testMessagePayload) {
		t.Fail()
	}
}

func TestQueue_CloseQueueImmediately_ThrowsNoErrors(t *testing.T) {
	gomega.RegisterTestingT(t)
	underTest := NewQueue("Test")

	close(underTest.InputChannel)

	gomega.Eventually(func() bool {
		_, open := <-underTest.OutputChannel
		return open
	}).Should(gomega.BeFalse())
}

func TestQueue_EvenNumberOfPushesAndPops_GivesZeroFinalLength(t *testing.T) {
	underTest := NewQueue("Test")
	numberOfRounds := 200

	for i := 0; i < numberOfRounds; i++ {
		dummyMessagePayLoad := []byte{byte(i)}
		dummyMessage := message.NewHeaderlessMsg(&dummyMessagePayLoad)
		underTest.InputChannel <- dummyMessage
	}

	gomega.Eventually(func() int {
		return underTest.length
	}).Should(gomega.Equal(numberOfRounds))

	for i := 0; i < numberOfRounds; i++ {
		message := <-underTest.OutputChannel
		if int((*message.Body)[0]) != i {
			t.Logf("Expected %d, got %d", i, int((*message.Body)[0]))
			t.FailNow()
		}
	}

	gomega.Eventually(func() int {
		return underTest.length
	}).Should(gomega.Equal(0))
}

func BenchmarkQueueSendRecv(b *testing.B) {
	dummyMessagePayLoad := []byte("Test")
	dummyMessage := message.NewHeaderlessMsg(&dummyMessagePayLoad)

	underTest := NewQueue("Test")

	for i := 0; i < b.N; i++ {
		underTest.InputChannel <- dummyMessage
	}

	for i := 0; i < b.N; i++ {
		_ = <-underTest.OutputChannel
	}
}
