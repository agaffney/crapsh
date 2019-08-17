package input

type Input interface {
	ReadLine(bool) (string, error)
	IsAvailable() bool
}
