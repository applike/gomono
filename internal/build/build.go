package build

type Builder interface {
	Build() error
	Test() error
	Deploy() error
}
