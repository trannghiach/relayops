package broker

import "github.com/nats-io/nats.go"

func NewNATS(url string) *nats.Conn {
	nc, err := nats.Connect(url)
	if err != nil {
		panic(err)
	}
	return nc
}
