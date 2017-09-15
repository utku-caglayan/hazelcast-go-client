package internal

import (
	"github.com/hazelcast/go-client/internal/protocol"
	"log"
	"time"
)

const (
	DEFAULT_HEARTBEAT_INTERVAL = 5
	DEFAULT_HEARTBEAT_TIMEOUT  = 60
)

type HeartBeatService struct {
	client            *HazelcastClient
	heartBeatTimeout  time.Duration
	heartBeatInterval time.Duration
	cancel            chan struct{}
}

func newHeartBeatService(client *HazelcastClient) *HeartBeatService {
	//TODO:: Add listeners
	heartBeat := HeartBeatService{client: client, heartBeatInterval: DEFAULT_HEARTBEAT_INTERVAL, heartBeatTimeout: DEFAULT_HEARTBEAT_TIMEOUT,
		cancel: make(chan struct{}),
	}
	return &heartBeat
}
func (heartBeat *HeartBeatService) start() {
	go func() {
		ticker := time.NewTicker(PARTITION_UPDATE_INTERVAL * time.Second)
		for {
			if !heartBeat.client.LifecycleService.isLive {
				return
			}
			select {
			case <-ticker.C:
				heartBeat.heartBeat()
			case <-heartBeat.cancel:
				ticker.Stop()
				return
			}
		}
	}()
}
func (heartBeat *HeartBeatService) heartBeat() {
	for _, connection := range heartBeat.client.ConnectionManager.connections {
		timeSinceLastRead := time.Since(connection.lastRead)
		if time.Duration(timeSinceLastRead.Seconds()) > heartBeat.heartBeatTimeout {
			if connection.heartBeating {

				log.Println("Didnt hear back from a connection")
				heartBeat.onHeartBeatStop(connection)
			}
		}
		if time.Duration(timeSinceLastRead.Seconds()) > heartBeat.heartBeatInterval {
			request := protocol.ClientPingEncodeRequest()
			heartBeat.client.InvocationService.InvokeOnConnection(request, connection)
		} else {
			if !connection.heartBeating {
				heartBeat.onHeartBeatRestored(connection)
			}
		}
	}
}
func (heartBeat *HeartBeatService) onHeartBeatRestored(connection *Connection) {
	log.Println("Heartbeat restored for a connection")
	connection.heartBeating = true
}
func (heartBeat *HeartBeatService) onHeartBeatStop(connection *Connection) {
	connection.heartBeating = false
}
func (heartBeat *HeartBeatService) shutdown() {
	close(heartBeat.cancel)
}
