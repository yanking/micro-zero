package contract

type IContainer interface {
	Name() string
	Run() error
	Register(Component)
}
