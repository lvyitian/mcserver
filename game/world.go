package game

import (
	"github.com/itay2805/mcserver/math"
	"github.com/itay2805/mcserver/minecraft"
	"github.com/itay2805/mcserver/minecraft/entity"
	"github.com/itay2805/mcserver/minecraft/proto/play"
	"github.com/itay2805/mcserver/minecraft/world"
	"sync/atomic"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// EID generration
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var eidGen int32

// Generate a new entity id for this session
func generateEntityId() int32 {
	return atomic.AddInt32(&eidGen, 1)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// World related
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type World struct {
	// the minecraft world
	*world.World

	// the entities
	entities	*math.Rtree
}

func NewWorld(generator world.Generator, provider world.Provider) *World {
	return &World{
		World:    world.NewWorld(provider, generator),
		entities: math.NewRTree(10, 2000),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Entity mangaement
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (w *World) AddPlayer(player *Player) {
	player.World = w

	// TODO: load from provider

	// set the spawn point
	pos := w.World.Generator.GetSpawn()
	player.Position = math.NewPoint(float64(pos.X), float64(pos.Y), float64(pos.Z))

	// insert the player to the entity list
	w.entities.Insert(player)

	// send the join game
	player.Send(play.JoinGame{
		EntityId:         		generateEntityId(),
		Gamemode:         		1,
		Dimension:        		0,
		HashedSeed: 	  		0,
		LevelType:        		"default",
		ViewDistance: 			10,
		ReducedDebugInfo: 		false,
		EnableRespawnScreen: 	true,
	})

	// send the player coords
	// TODO: send respawn packet if need to change world
	player.Send(play.PlayerPositionAndLook{
		X:          player.Position.X(),
		Y:          player.Position.Y(),
		Z:          player.Position.Z(),
		Yaw:        0,
		Pitch:      0,
		Flags:      0,
		TeleportId: 69,
	})
}

func (w *World) UpdateEntityPosition(entity entity.IEntity) {
	w.entities.Delete(entity)
	entity.UpdateBounds()
	w.entities.Insert(entity)
}

func (w *World) RemovePlayer(p *Player) {
	w.entities.Delete(p)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Others
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (w *World) SendChunkToPlayer(x, z int, p *Player) {
	// load the chunk
	c := w.LoadChunk(x, z)

	// the chunk itself
	writer := minecraft.Writer{}
	c.MakeChunkDataPacket(&writer)

	// block entities
	writer.WriteVarint(0)

	// send the data
	p.SendRaw(writer.Bytes())
}
