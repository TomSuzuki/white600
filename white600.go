package white600

func MarkdownToHTML(markdown string) string {

	token := lexer(markdown)
	html := generator(token)

	return html
}
