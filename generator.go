package white600

// generator ...
func generator(token []Token) string {

	var html string

	for _, v := range token {
		html += generatorTag(v)
	}

	return html
}

// generatorTag ...
func generatorTag(token Token) string {
	switch token.elementType {
	case typeText:
		return token.content
	case typeBreak:
		return "<br>"
	case typeBlockCodeOpen:
		return "<pre><code>"
	case typeBlockCodeClose:
		return "</code></pre>"
	case typeListBlockOpen:
		return "<ul>"
	case typeListBlockClose:
		return "</ul>"
	case typeNumberListOpen:
		return "<ol>"
	case typeNumberListClose:
		return "</ol>"
	case typeList:
		return "<li>"
	case typeListClose:
		return "</li>"
	case typeHeader1:
		return "<h1>" + token.content + "</h1>"
	case typeHeader2:
		return "<h2>" + token.content + "</h2>"
	case typeHeader3:
		return "<h3>" + token.content + "</h3>"
	case typeHeader4:
		return "<h4>" + token.content + "</h4>"
	case typeBoldOpen:
		return "<strong>"
	case typeBoldClose:
		return "</strong>"
	case typeCancellationOpen:
		return "<s>"
	case typeCancellationClose:
		return "</s>"
	case typeItalicOpen:
		return "<em>"
	case typeItalicClose:
		return "</em>"
	case typeCodeOpen:
		return "<code>"
	case typeCodeClose:
		return "</code>"
	case typeHorizon:
		return "<hr>"
	case typeLink:
		return "<a href='" + token.content + "'>"
	case typeLinkClose:
		return "</a>"
	case typeParagraphOpen:
		return "<p>"
	case typeParagraphClose:
		return "</p>"
	}

	return ""
}
