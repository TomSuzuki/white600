package white600

import "strings"

// tableAlign ... : --- :
var tableAlign = map[[2]bool]string{
	{true, false}:  "right",  // : ---
	{false, true}:  "left",   //   --- :
	{true, true}:   "center", // : --- :
	{false, false}: "center", //   ---
}

// convTableHead ...make align
func (data *MarkdownInfo) convTableHead() {
	// 次行の情報が必須なので無ければエラー扱い
	if len(data.markdown) == 1 {
		return
	}

	// align
	alignLine := strings.Split(data.currentData.nextLine, "|")
	data.markdown = append([]string{data.markdown[0]}, data.markdown[2:]...) // 位置指定の行を消す
	alignLine = alignLine[1 : len(alignLine)-1]
	for _, v := range alignLine {
		data.options.tableAlign = append(data.options.tableAlign, tableAlign[[2]bool{string(v[len(v)-1]) == ":", string(v[0]) == ":"}])
	}

	// テーブルヘッダの開始
	data.html = append(data.html, "<table><thead>")

	// テーブルヘッダを解析
	data.tableGenerate("th")

	// open <table><thead>
	//data.markdownLines[0] = ("<table><thead>" + data.markdownLines[0])

}

// closeTableHead ...if table is close
func (data *MarkdownInfo) closeTableHead() {
	if data.currentData.lineType != typeTableBody {
		data.shiftLine()
		data.html = append(data.html, "</thead></table>")
		data.options.tableAlign = nil
	}
}

// convTableBody ...table generation
func (data *MarkdownInfo) convTableBody() {

	// <tbody>
	if data.currentData.isNewBlock {
		data.html = append(data.html, "</thead><tbody>")
		//data.html = append(data.html, data.inlineConv(inner)) // todo インラインの位置をチェック
	}

	// <tr>
	data.tableGenerate("td")
}

// closeTableBody ...
func (data *MarkdownInfo) closeTableBody() {
	data.shiftLine()
	data.html = append(data.html, "</tbody></table>")
	data.options.tableAlign = nil
}

// tableGenerate ...<tr>
func (data *MarkdownInfo) tableGenerate(tagType string) {
	// check
	var tr = strings.Split(data.currentData.currentLine, "|")
	if len(tr)-2 <= 1 || len(tr)-2 != len(data.options.tableAlign) || tr[0] != "" || tr[len(tr)-1] != "" {
		return
	}

	// make
	var html []string
	html = append(html, "<tr>")
	for i, v := range data.options.tableAlign {
		inner := data.inlineConv(tr[i+1])
		html = append(html, "<")
		html = append(html, tagType)
		html = append(html, " align='")
		html = append(html, v)
		html = append(html, "'>")
		html = append(html, inner)
		html = append(html, "</")
		html = append(html, tagType)
		html = append(html, ">")
	}
	html = append(html, "</tr>")

	data.html = append(data.html, html...)
}
