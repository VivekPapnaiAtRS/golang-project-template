package realtimesocketmanager

import (
	"VivekPapnaiAtRS/template/providers"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const defaultPartitions = 5

// RealtimeHub maintains the set of active clients and broadcasts messages to the clients.
type RealtimeHub struct {
	// registered clients. Each RealtimeClient is identified by a user ID and a
	// uuID. There are a set of user IDs and each user ID has a set of uuIDs.
	// Each uuID identifies a unique RealtimeClient.
	clients map[int]map[uuid.UUID]*RealtimeClient

	// register requests from the clients.
	register chan *RealtimeClient

	// unregister clients from RealtimeHub
	unregister chan *RealtimeClient

	// getClients retrieves a gl
	getClients    chan int
	outGetClients chan getClientsResp
	DB            *sqlx.DB
	DBHelper      providers.DBHelpProvider

	Done chan bool
}

type getClientsResp struct {
	clients map[uuid.UUID]*RealtimeClient
	err     error
}

func NewRealtimeHub(db *sqlx.DB, helper providers.DBHelpProvider) providers.WebSocketHubProvider {
	return &RealtimeHub{
		register:      make(chan *RealtimeClient, 1),
		unregister:    make(chan *RealtimeClient, 1),
		getClients:    make(chan int, 1),
		outGetClients: make(chan getClientsResp, 1),
		Done:          make(chan bool, 1),
		clients:       make(map[int]map[uuid.UUID]*RealtimeClient),
		DB:            db,
		DBHelper:      helper,
	}
}

func (h *RealtimeHub) Run() {
	for {
		select {
		case <-h.Done:
			h.Stop()
			return
		case client := <-h.register:
			if _, ok := h.clients[client.userID]; !ok {
				h.clients[client.userID] = make(map[uuid.UUID]*RealtimeClient)
			}
			h.clients[client.userID][client.uuID] = client

		case client := <-h.unregister:
			delete(h.clients[client.userID], client.uuID)

		case userID := <-h.getClients:
			clients, ok := h.clients[userID]
			if !ok {
				h.outGetClients <- getClientsResp{
					clients: nil,
					err:     fmt.Errorf("no Clients are associated with userID %v", userID),
				}
				break
			}
			h.outGetClients <- getClientsResp{
				clients: clients,
				err:     nil,
			}
		}
	}
}

func (h *RealtimeHub) Get() interface{} {
	return h
}

func (h *RealtimeHub) Stop() {
	for _, userClients := range h.clients {
		for _, client := range userClients {
			close(client.send)
		}
	}
}
