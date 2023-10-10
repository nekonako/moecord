package sfu

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/jiyeyuran/mediasoup-go"
	"github.com/nekonako/moecord/config"
	"github.com/rs/zerolog/log"
)

type SFU struct {
	*sync.RWMutex
	config     *config.Config
	workers    []*mediasoup.Worker
	rooms      map[string]Room
	nextWorker int
}

type Room struct {
	ID     string
	Name   string
	Router *mediasoup.Router
}

func New(c *config.Config) *SFU {

	mediasoup.WorkerBin = c.MediaSoup.WorkerPath
	workers := []*mediasoup.Worker{}
	for i := 0; i < c.MediaSoup.NumWorkers; i++ {
		worker, err := mediasoup.NewWorker()
		if err != nil {
			panic(err)
		}
		workers = append(workers, worker)

		go func() {
			ticker := time.NewTicker(120 * time.Second)
			for {
				select {
				case <-ticker.C:
					usage, err := worker.GetResourceUsage()
					if err != nil {
						log.Err(err).Int("pid", worker.Pid()).Msg("mediasoup Worker resource usage")
						continue
					}
					log.Info().Int("pid", worker.Pid()).Interface("usage", usage).Msg("mediasoup Worker resource usage")
				}
			}
		}()
	}

	fmt.Println(workers)

	return &SFU{
		RWMutex: &sync.RWMutex{},
		config:  c,
		workers: workers,
		rooms:   make(map[string]Room),
	}
}

// user

func (s *SFU) CreateOrGetRoom(c net.Conn, roomID, name string) {
	s.Lock()
	defer s.Unlock()

	worker := s.workers[s.nextWorker]
	mediasoupRouter, err := worker.CreateRouter(mediasoupRouterConfig)
	if err != nil {
		log.Error().Msg(err.Error())
		panic(err)
	}

	_, err = mediasoupRouter.CreateWebRtcTransport(mediasoup.WebRtcTransportOptions{
		ListenIps: s.config.MediaSoup.WebRTCTransportOptions.ListenIPs,
	})
	if err != nil {
		panic(err)
	}

	if _, exist := s.rooms[roomID]; !exist {
		s.rooms[roomID] = Room{
			ID:     roomID,
			Name:   name,
			Router: mediasoupRouter,
		}
	}

	s.nextWorker = (s.nextWorker + 1) % len(s.workers)
	fmt.Println(s.rooms)
	fmt.Println("hello world from room")
}
