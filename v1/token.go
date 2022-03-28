package white600

type elementType int

const (
	typeRoot elementType = iota
	typeText
	typeBreak
	typeBlockCodeOpen
	typeBlockCodeClose
	typeListBlockOpen
	typeListBlockClose
	typeNumberListOpen
	typeNumberListClose
	typeList
	typeListClose
	typeHeader1
	typeHeader2
	typeHeader3
	typeHeader4
	typeBoldOpen
	typeBoldClose
	typeCancellationOpen
	typeCancellationClose
	typeItalicOpen
	typeItalicClose
	typeCodeOpen
	typeCodeClose
	typeHorizon
	typeLink
	typeLinkClose
	typeParagraphOpen
	typeParagraphClose
)

type Token struct {
	content     string
	elementType elementType
}
