package white600

// convCodeMarker ...start code lines
func (data *MarkdownInfo) convCode() {
	if data.currentData.isNewBlock {
		data.html = append(data.html, "<pre><code>")
		data.currentData.lineType = typeCode
	} else {
		data.html = append(data.html, data.currentData.currentLine)
		data.html = append(data.html, "\n")
	}
}

// closeCode ...close code lines
func (data *MarkdownInfo) closeCode() {
	data.html = append(data.html, "</code></pre>")
	data.currentData.lineType = typeNone
}
