package main

import "fmt"

type Sender struct {
	SendMethodImpl func(message string)
}

func NewCourierSender() *Sender {
	return &Sender{SendMethodImpl: func(message string) {
		fmt.Printf("Message %s was sent by courier\n", message)
	}}
}

func NewEmailSender() *Sender {
	return &Sender{SendMethodImpl: func(message string) {
		fmt.Printf("Message %s was sent by email\n", message)
	}}
}

func (s *Sender) Send(message string) {
	s.SendMethodImpl(message)
}

func main() {
	var sender Sender

	sender = *NewEmailSender()
	sender.Send("1")

	sender = *NewCourierSender()
	sender.Send("2")
}
