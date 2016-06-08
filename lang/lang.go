package lang

type FactoryFunc func(*Generic) Element

type ElementEntry struct {
	Name       string
	ParserData []*ParserHint
	Factory    FactoryFunc
}

func GetElementEntry(element string) *ElementEntry {
	for _, e := range Elements {
		if e.Name == element {
			return e
		}
	}
	return nil
}

var Elements []*ElementEntry

func init() {
	Elements = make([]*ElementEntry, 0)
}

func registerElements(elements []*ElementEntry) {
	Elements = append(Elements, elements...)
}
