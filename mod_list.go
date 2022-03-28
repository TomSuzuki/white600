package white600

import (
	"strings"
)

// convList ...list generation
func (data *MarkdownInfo) convList() {
	var text []string
	var line = data.currentData.currentLine
	var nest = 0
	var oldNest = len(data.options.listNest)

	// list type and open list
	for md, tag := range listStyle {
		if strings.Index(strings.Trim(line, " "), md) == 0 {
			nest = 1 + strings.Index(line, md)/2
			line = line[strings.Index(line, md)+len(md):]

			// open <ul> or <ol>
			for nest > len(data.options.listNest) {
				data.options.listNest = append(data.options.listNest, tag)
				text = append(text, "<")
				text = append(text, tag)
				text = append(text, ">")
			}
		}
	}

	// close
	data.listTagClose(nest, oldNest)

	// open <li>
	text = append(text, "<li>")
	text = append(text, data.inlineConv(line))
	data.html = append(data.html, text...)
}

// closeList ...close list
func (data *MarkdownInfo) closeList() {
	data.listTagClose(0, len(data.options.listNest))
}

// listTagClose
func (data *MarkdownInfo) listTagClose(nest int, oldNest int) {
	var text = []string{""}

	// close
	for nest < len(data.options.listNest) {
		text = append(text, "</li></")
		text = append(text, data.options.listNest[len(data.options.listNest)-1])
		text = append(text, ">")
		data.options.listNest = data.options.listNest[:len(data.options.listNest)-1]
	}

	// </li>
	if nest <= oldNest && nest != 0 {
		if nest == oldNest {
			text = append(text[:1], text...)
			text[0] = "</li>"
		} else {
			text = append(text, "</li>")
		}
	}

	// append
	//text = append(text, data.markdownLines[0])

	// join
	//data.markdownLines[0] = strings.Join(text, "")
	data.html = append(data.html, text...)
}
