package models

type GameType struct {
	Id          int
	Description string
}

type GameTypeService interface {
	GetGameTypeById(gameTypeId int) (GameType, error)
}
