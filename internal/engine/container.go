package engine

type ContainerStatus int

const (
	ContainerStatusUnknown    ContainerStatus = iota
	ContainerStatusRunning    ContainerStatus = iota
	ContainerStatusRestarting ContainerStatus = iota
)

type Container struct {
	ID     string
	Status ContainerStatus
}
