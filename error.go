package goserver

type (
	Error interface {
		error
		Status() int
	}
)
