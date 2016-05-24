package core

type Config struct {
	Binary          string
	Command         string
	CommandProvided bool
	File            string
	FileProvided    bool
	Args            []string
}
