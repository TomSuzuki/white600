package white600

import "strings"

// headText
const headText = "123456"

// convHeader ...<h1> - <h6>の解析を行う。
func (data *MarkdownInfo) convHeader() {
	data.currentData.currentLine = strings.Trim(data.currentData.currentLine, " ")
	h := strings.Count(strings.Split(data.currentData.currentLine, " ")[0], "#")
	if h <= 6 && h >= 1 {
		inner := ""
		if len(data.currentData.currentLine) > h+1 {
			inner = data.inlineConv(data.currentData.currentLine[h+1:])
		}
		data.html = append(data.html, "<h")
		data.html = append(data.html, headText[h-1:h])
		data.html = append(data.html, ">")
		data.html = append(data.html, inner)
		data.html = append(data.html, "</h")
		data.html = append(data.html, headText[h-1:h])
		data.html = append(data.html, ">")
	}
}
