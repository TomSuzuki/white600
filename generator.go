package white600

import (
	"strings"
)

// generator ...本体
func generator(markdown string) string {

	// 初期状態の生成
	data := initInfo(markdown)

	// 本体
	for len(data.markdown) > 0 {
		// 現在の行を処理
		data.currentData = LineInfo{currentLine: data.markdown[0]}
		data.currentData.lineType = data.getLineType()
		data.currentData.isNewBlock = data.currentData.lineType != data.previousData.lineType
		if len(data.markdown) > 1 {
			data.currentData.nextLine = data.markdown[1]
		}

		// タイプが変わったらブロックを閉じる
		if data.currentData.isNewBlock {
			data.closeBlock()
		}

		// 現在の行を解析
		data.convBlock()

		// 次の行へ引き続き
		data.previousData = data.currentData
		data.markdown = data.markdown[1:]

	}

	// ブロックを閉じる処理
	data.closeBlock()

	return strings.Join(data.html, "")
}

// convBlock ...ブロックを解析
func (data *MarkdownInfo) convBlock() {
	switch data.currentData.lineType {
	case typeParagraph:
		data.convParagraph()
	case typeHeader:
		data.convHeader()
	// case typeTableBody:
	// 	convData.closeTableBody()
	// case typeTableHead:
	// 	data.convTableHead()
	case typeCode, typeCodeMarker:
		data.convCode()
	case typeList:
		data.convList()
	case typeQuote:
		data.convQuote()
	case typeHorizon:
		data.convHorizon()
	}
}

// closeBlock ...ブロックを閉じる
func (data *MarkdownInfo) closeBlock() {
	switch data.previousData.lineType {
	case typeParagraph:
		data.closeParagraph()
	// case typeTableBody:
	// 	convData.closeTableBody()
	// case typeTableHead:
	// 	convData.closeTableHead()
	case typeCode, typeCodeMarker:
		data.closeCode()
	case typeList:
		data.closeList()
	case typeQuote:
		data.closeQuote()
	}
}

// shiftLine ...現在の行に空行を挿入する。
func (data *MarkdownInfo) shiftLine() {
	data.markdown = append([]string{""}, data.markdown...)
	data.currentData.nextLine = data.currentData.currentLine
	data.currentData.currentLine = ""
}

// initInfo ...状態の初期化処理
func initInfo(markdown string) MarkdownInfo {
	return MarkdownInfo{
		markdown: append(strings.Split(strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(markdown), "\n"), ""),
		html:     []string{},
		currentData: LineInfo{
			lineType:    typeNone,
			currentLine: "",
			nextLine:    "",
		},
		previousData: LineInfo{
			lineType:    typeNone,
			currentLine: "",
			nextLine:    "",
		},
		options: Options{
			listNest:   []string{},
			tableAlign: []string{},
		},
	}
}
