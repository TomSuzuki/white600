package white600

func MarkdownToHTML(markdown string) string {

	html := generator(markdown)

	return html
}
