package white600

import "strings"

// getLineType ...現在の行のタイプをチェックする
func (markdownInfo *MarkdownInfo) getLineType() LineType {
	// 優先度が高いものが上
	switch true {
	case markdownInfo.isCodeMarker():
		return typeCodeMarker
	case markdownInfo.isCode():
		return typeCode
	case markdownInfo.isHeader():
		return typeHeader
	case markdownInfo.isList():
		return typeList
	case markdownInfo.isQuote():
		return typeQuote
	case markdownInfo.isHorizon():
		return typeHorizon
	case markdownInfo.isTableBody():
		return typeTableBody
	case markdownInfo.isTableHead():
		return typeTableHead
	case markdownInfo.isNone():
		return typeNone
	default:
		return typeParagraph
	}
}

// isCode ...コードブロックであるかを判定する。
func (markdownInfo *MarkdownInfo) isCode() bool {
	return markdownInfo.previousData.lineType == typeCodeMarker || markdownInfo.previousData.lineType == typeCode
}

// isCodeMarker ...コードブロックの端子であるかを判定する。
func (markdownInfo *MarkdownInfo) isCodeMarker() bool {
	return len(markdownInfo.currentData.currentLine) >= 3 && markdownInfo.currentData.currentLine[:3] == "```"
}

// isHeader ...markdownLines[0] is header?
func (markdownInfo *MarkdownInfo) isHeader() bool {
	line := strings.Split(strings.Trim(markdownInfo.currentData.currentLine, " "), " ")[0]
	return line != "" && strings.Trim(line, "#") == ""
}

// isHorizon ...水平線ブロックであるかを判定する。
func (markdownInfo *MarkdownInfo) isHorizon() bool {
	var line = strings.Replace(markdownInfo.currentData.currentLine, " ", "", -1)
	return len(line) >= 3 && (len(line) == strings.Count(line, "-") || len(line) == strings.Count(line, "_") || len(line) == strings.Count(line, "*"))
}

// isList ...リストブロックであるかを判定する。
func (markdownInfo *MarkdownInfo) isList() bool {
	var line = strings.Trim(markdownInfo.currentData.currentLine, " ")
	for md := range listStyle {
		if strings.Index(line, md) == 0 {
			return true
		}
	}
	return false
}

// isNone ...空行であるかを判定する。
func (markdownInfo *MarkdownInfo) isNone() bool {
	return strings.Trim(markdownInfo.currentData.currentLine, " ") == "" && markdownInfo.previousData.lineType != typeParagraph
}

// isQuote ...引用であるかを判定する。
func (markdownInfo *MarkdownInfo) isQuote() bool {
	return (strings.Trim(markdownInfo.currentData.currentLine, " ") + "  ")[:2] == "> " || (markdownInfo.previousData.lineType == typeQuote && !markdownInfo.isNone())
}

// isTableHead ...count "|"
func (markdownInfo *MarkdownInfo) isTableHead() bool {
	var line = markdownInfo.currentData.currentLine
	return len(line) >= 2 && strings.Trim(line, " ")[:1] == "|" && strings.Count(line, "|") > 1
}

// isTableBody ...thead and before type check
func (markdownInfo *MarkdownInfo) isTableBody() bool {
	return (markdownInfo.previousData.lineType == typeTableBody || markdownInfo.previousData.lineType == typeTableHead) && markdownInfo.isTableHead()
}
