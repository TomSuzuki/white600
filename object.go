package white600

type LineType int

// LineType ...行のタイプ
const (
	typeNone LineType = iota
	typeParagraph
	typeList
	typeCode
	typeCodeMarker
	typeTableHead
	typeTableBody
	typeQuote
	typeHeader
	typeHorizon
)

// Data ...処理中の状態をすべて持つ
type MarkdownInfo struct {
	markdown     []string
	html         []string
	currentData  LineInfo
	previousData LineInfo
	options      Options
}

// Options ...全体の処理状態を表す（次の行へ引き継ぐ）
type Options struct {
	listNest   []string
	tableAlign []string
	nestQuote  int
}

// LineInfo ...1つの行の状態を表す
type LineInfo struct {
	lineType    LineType
	isNewBlock  bool
	currentLine string
	nextLine    string
}

// listStyle ...リストの種類を列挙。
var listStyle = map[string]string{
	"- ":  "ul",
	"* ":  "ul",
	"1. ": "ol",
	"+ ":  "ul",
}
