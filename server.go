package main

import (
	"game_server/uuid"
	"os"
	"os/signal"

	"github.com/codecat/go-enet"
	"github.com/codecat/go-libs/log"
)

var clients = make(map[enet.Peer]int)
var peers = make(map[int]enet.Peer)
var colletionIDs = uuid.New()

// if connected, add to clients
func onConnected(peer enet.Peer) {
	id := colletionIDs.Next()
	log.Info("A new client connected from %s id: %v", peer.GetAddress().String(), id)
	clients[peer] = id
	peers[id] = peer
}

// if received a message handle it
func onReceived(ev enet.Event) {
	peer := ev.GetPeer()
	packet := ev.GetPacket()
	packetBytes := packet.GetData()

	if string(packetBytes) == "ping" {
		//peer.SendString("pong", ev.GetChannelID(), enet.PacketFlagReliable)
		sendtoclient(clients[peer], "pong")
		return
	}
	if string(packetBytes) == "bye" {
		log.Info("Bye!")
		peer.Disconnect(0)
		return
	}
	packet.Destroy()
}

// if disconnected, remove from clients
func onDisconnected(peer enet.Peer) {
	id := clients[peer]
	log.Info("A client disconnected from %s %v", peer.GetAddress().String(), id)
	delete(peers, id)
	delete(clients, peer)
	colletionIDs.Free(id)
}

// Poll
func poll(host enet.Host) {
	ev := host.Service(0)
	if ev.GetType() == enet.EventNone {
		return
	}
	switch ev.GetType() {
	case enet.EventConnect:
		onConnected(ev.GetPeer())
	case enet.EventDisconnect:
		onDisconnected(ev.GetPeer())
	case enet.EventReceive:
		onReceived(ev)
	}
}

// Send a message to all clients
func sendToAll(msg string) {
	//peers = make(map[int]enet.Peer)
	for _, peer := range peers {
		peer.SendString(msg, 0, enet.PacketFlagReliable)
	}
}

// Send a message to a client
func sendtoclient(id int, msg string) {
	peer := peers[id]
	peer.SendString(msg, 0, enet.PacketFlagReliable)
}

//
func main() {
	enet.Initialize()
	host, err := enet.NewHost(enet.NewListenAddress(5000), 32, 1, 0, 0)
	if err != nil {
		log.Error("Error creating host: %s", err.Error())
		return
	}
	running := true
	// Detect Ctrl+C pressed in console?
	finalize := make(chan os.Signal, 1)
	signal.Notify(finalize, os.Interrupt)
	go func() {
		<-finalize
		running = false
	}()
	// Main loop
	for running {
		poll(host)
	}
	//
	host.Destroy()
	enet.Deinitialize()
}
