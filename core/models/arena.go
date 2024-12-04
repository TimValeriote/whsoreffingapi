package models

type Arena struct {
	Id   int
	Name string
}

type ArenaService interface {
	GetArenaById(arenaId int) (Arena, error)
}
