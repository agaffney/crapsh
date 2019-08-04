package lang

const (
	HINT_TYPE_NODE = iota
	HINT_TYPE_ELEMENT
	HINT_TYPE_TOKEN
	HINT_TYPE_GROUP
	HINT_TYPE_ANY
)

type ParserHint struct {
	Type     int
	Name     string        // Name of element or token to match
	Optional bool          // Hint is optional
	Many     bool          // Hint can match multiple times
	Final    bool          // Consider the element matched if this hint matches
	Tokens   []string      // List of token names to match (for TOKEN type)
	Members  []*ParserHint // Child parser hints (used by ANY/GROUP hint types)
}
