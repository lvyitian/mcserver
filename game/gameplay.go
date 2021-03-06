package game

import (
	"github.com/google/uuid"
	"github.com/itay2805/mcserver/minecraft/world/generator/flatgrass"
	"github.com/itay2805/mcserver/minecraft/world/provider/nullprovider"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Game loop api
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// list of all players
var players = make(map[uuid.UUID]*Player)
var playerCount int32

var newPlayers []*Player
var newPlayersMutex	sync.Mutex

var leftPlayers []*Player
var leftPlayersMutex sync.Mutex

func GetPlayerCount() int32 {
	return atomic.LoadInt32(&playerCount)
}

//
// Add a player to the server
//
func JoinPlayer(player *Player) {
	newPlayersMutex.Lock()
	defer newPlayersMutex.Unlock()
	newPlayers = append(newPlayers, player)
}

func LeftPlayer(player *Player) {
	leftPlayersMutex.Lock()
	defer leftPlayersMutex.Unlock()
	leftPlayers = append(leftPlayers, player)
}

// the world we are going to use
var OurWorld = NewWorld(
	&flatgrass.FlatgraassGenerator{},
	&nullprovider.NullProvider{},
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Game loop sync stages
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func syncServerState() {
	// sync any changes made by the clients in
	// the last tick
	for _, p := range players {
		p.syncChanges()
	}

	// add all new players to the world
	for _, p := range newPlayers {
		OurWorld.AddPlayer(p)
		players[p.UUID] = p
		playerCount++
	}

	// remove any players that have left from
	// the world
	for _, p := range leftPlayers {
		delete(players, p.UUID)
		p.World.RemovePlayer(p)
		playerCount--
	}
}

func syncClientsState() {
	for _, p := range players {
		p.syncClient()
	}
}

func cleanupTick() {
	for _, p := range players {
		p.cleanupTick()
	}

	// finished updating players
	// for this tick
	newPlayers = nil
	leftPlayers = nil
}

func tickAllPlayers() {
	for _, p := range players {
		p.tick()
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// The game loop
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func StartGameLoop() {
	log.Println("Starting gameloop")

	ticker := time.NewTicker(time.Second / 20)
	lastLog := time.Now()
	ticks := 0
	for range ticker.C {

		// any thing that can't be changed while we are
		// doing a game tick need to be locked here
		leftPlayersMutex.Lock()
		newPlayersMutex.Lock()

		// start by syncing all of the states, this is required
		// to process any change sent out from players that we
		// want to know about this tick
		syncServerState()

		// process all of the player actions
		tickAllPlayers()

		// cause all the objects that need to run this tick to run
		// note that some objects (like player actions) do not run
		// this tick
		tickAllObjects()

		// finally sync all of the clients states, this is called
		// after the server did all of its processing
		syncClientsState()

		// Do any last cleanups before we are finishing up with
		// this tick
		cleanupTick()

		// increment the game ticks
		ticks++

		// log the tps, just for sanity
		if time.Since(lastLog) > time.Second {
			log.Printf("%d tp%s",
				ticks,
				time.Since(lastLog).Round(time.Second))
			ticks = 0
			lastLog = time.Now()
		}

		// and released here
		newPlayersMutex.Unlock()
		leftPlayersMutex.Unlock()
	}
}
