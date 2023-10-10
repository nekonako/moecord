package sfu

import (
	"sync"

	"github.com/jiyeyuran/mediasoup-go"
)

type Peer struct {
	sync.Mutex
	ID              string
	Name            string
	RtpCapabilities *mediasoup.RtpCodecCapability
	transports      map[string]mediasoup.ITransport
	producers       map[string]*mediasoup.Producer
	consumers       map[string]*mediasoup.Consumer
	dataProducers   map[string]*mediasoup.DataProducer
	dataConsumers   map[string]*mediasoup.DataConsumer
}
