package white600

import (
	"strings"
)

// convQuote ...引用の解析。
func (data *MarkdownInfo) convQuote() {
	// >
	nest := strings.Count(data.currentData.currentLine, "> ")
	inQuote := data.currentData.currentLine
	if nest == 0 {
		nest = data.options.nestQuote
	} else {
		inQuote = inQuote[2*nest:]
	}

	// インライン解析
	inQuote = data.inlineConv(inQuote)

	// open
	if data.options.nestQuote < nest {
		var text []string
		var oldNest = data.options.nestQuote
		for data.options.nestQuote < nest {
			data.options.nestQuote++
			if oldNest != 0 {
				text = append(text, "</p>")
			}
			text = append(text, "<blockquote>")

		}
		text = append(text, "<p>")
		data.html = append(data.html, text...)
	}

	// インライン要素を追加
	data.html = append(data.html, inQuote)

	// close
	data.quoteTagClose(nest)

	// inline
	// convData.inlineConv()
}

// closeQuote ...引用ブロックを閉じる。
func (data *MarkdownInfo) closeQuote() {
	data.shiftLine()
	data.quoteTagClose(0)
}

// quoteTagClose ...引用タグを閉じる。
func (data *MarkdownInfo) quoteTagClose(nest int) {
	if data.options.nestQuote > nest {
		var text []string
		text = append(text, "</p>")
		text = append(text, data.currentData.currentLine)
		for data.options.nestQuote > nest {
			text = append(text, "</blockquote>")
			data.options.nestQuote--
		}
		//data.currentData.currentLine = strings.Join(text, "")
		data.html = append(data.html, text...)
	}
}
