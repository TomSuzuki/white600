package white600

// closeParagraph ...パラグラフを閉じる
func (data *MarkdownInfo) closeParagraph() {
	data.html = append(data.html, "</p>")
}

// closeParagraph ...パラグラフを解析
func (data *MarkdownInfo) convParagraph() {
	// 新しいブロックなら開タグを追加
	if data.currentData.isNewBlock {
		data.html = append(data.html, "<p>")
	}

	// インライン解析
	inner := data.inlineConv(data.currentData.currentLine)

	// 解析結果を追加
	data.html = append(data.html, inner)
}
