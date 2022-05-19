package interfaces

type Executor interface {
	Setup() error
	Teardown() error
}
