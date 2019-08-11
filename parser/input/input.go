package input

type Input interface {
	ReadLine() (string, error)
	ReadAnotherLine() (string, error)
}
