package services

type Operation uint8

const (
	SVC_OPERATION_CREATE Operation = iota
	SVC_OPERATION_UPDATE
	SVC_OPERATION_DELETE
)
