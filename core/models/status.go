package models

type Status struct {
	Id          int
	Description string
}

type StatusService interface {
	GetStatusById(statusId int) (Status, error)
}
