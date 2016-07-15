package goapi

type (
	Error interface {
		error
		Status() int
	}
)
