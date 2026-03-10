package modules

type Status string

const (
	StatusIdle  Status = "idle"
	StatusBusy  Status = "busy"
	StatusError Status = "error"
)

type Module interface {
	Name() string
	Status() Status
	Handle(req any) (any, error)
}
