package Agora

type MessageHandler interface {
	Initialize(<-chan string) chan<- string
}
