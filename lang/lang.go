package lang

type FactoryFunc func(*Generic) Element

type ElementEntry struct {
	Name       string
	ParserData []*ParserHint
	Factory    FactoryFunc
}

func GetElements(elements []string) []*ElementEntry {
	ret := []*ElementEntry{}
	for _, e := range elements {
		for _, element := range Elements {
			if e == element.Name {
				ret = append(ret, element)
				break
			}
		}
	}
	return ret
}

var Elements []*ElementEntry

func init() {
	Elements = make([]*ElementEntry, 0)
}

func registerElements(elements []*ElementEntry) {
	Elements = append(Elements, elements...)
}
