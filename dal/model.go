package dal

type Model interface {
	Identifier() int64
	Validate() error
}
