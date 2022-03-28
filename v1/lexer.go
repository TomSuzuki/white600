package white600

import (
	"strings"
)

// Options ...現在の状態を持つ
type Options struct {
	isCodeBlock       bool
	isParagraph       bool
	isListBlock       int
	isNumberListBlock int
}

// lexer ...構文の解析を行う
func lexer(markdown string) []Token {

	// 解析結果
	token := []Token{}

	// 初期状態
	options := Options{
		isCodeBlock:       false,
		isParagraph:       false,
		isListBlock:       0,
		isNumberListBlock: 0,
	}

	// 改行で区切る
	lines := append(strings.Split(strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(markdown), "\n"), "")

	// 解析
	for len(lines) > 0 {
		nowline := lines[0]
		token, options = lineLexer(token, options, nowline)
		lines = lines[1:]
	}

	// 全部閉じる
	token, options = lexerCloseBlock(token, options)

	return token
}

// lineLexer ...行ごとの解析を行なう
func lineLexer(token []Token, options Options, line string) ([]Token, Options) {

	// コードブロック
	if lexerCodeBlock(line) {
		token, options = lexerCloseBlock(token, options)
		options.isCodeBlock = !options.isCodeBlock
		if options.isCodeBlock {
			token = append(token, Token{elementType: typeBlockCodeOpen})
		}
		return token, options
	} else if options.isCodeBlock {
		return append(token, Token{
			content:     line,
			elementType: typeText,
		}), options
	}

	// TODO: テーブルブロック

	// リストブロック
	level, text := lexerListBlock(line)
	if level > 0 {
		// リストの開始なら他のブロックを閉じる
		if options.isListBlock == 0 {
			token, options = lexerCloseBlock(token, options)
		}

		// リストを開く
		for level > options.isListBlock {
			token = append(token, Token{elementType: typeListBlockOpen})
			options.isListBlock++
		}

		// 中身を追加
		token = append(token, Token{elementType: typeList})
		token = append(token, lexerInline([]Token{}, text)...)
		token = append(token, Token{elementType: typeListClose})

		// リストを閉じる
		for level < options.isListBlock {
			token = append(token, Token{elementType: typeListBlockOpen})
			options.isListBlock--
		}

		return token, options
	}

	// 水平線
	if lexerHorizon(line) {
		return append(token, Token{elementType: typeHorizon}), options
	}

	// ヘッダーブロック
	level, text = lexerHeaderBlock(line)
	if level > 0 {
		token, options = lexerCloseBlock(token, options)
		token = append(token, Token{
			content: text,
			elementType: map[int]elementType{
				1: typeHeader1,
				2: typeHeader2,
				3: typeHeader3,
				4: typeHeader4,
			}[level],
		})
		return token, options
	}

	// 空行
	if lexerEmpty(line) {
		if options.isParagraph {
			return append(token, Token{elementType: typeBreak}), options
		} else {
			return token, options
		}
	}

	// パラグラフの開始
	if !options.isParagraph {
		token, options = lexerCloseBlock(token, options)
		token = append(token, Token{elementType: typeParagraphOpen})
		options.isParagraph = true
	}

	// インライン要素の解析
	token = append(token, lexerInline([]Token{}, line)...)

	return token, options
}

// lexerInline ...インライン要素を解析して返す（再帰？）
func lexerInline(token []Token, text string) []Token {

	// リンク関連の解析
	token, text = lexerInlineLink(token, text)

	// 辞書用の型
	type dictInfo struct {
		tag          string
		openElement  elementType
		closeElement elementType
		recursion    bool
	}

	// テキストのパターンと対応するタイプ（場所移動する）
	dictionary := []dictInfo{
		{tag: "`", openElement: typeCodeOpen, closeElement: typeCodeClose, recursion: false},
		{tag: "**", openElement: typeBoldOpen, closeElement: typeBoldClose, recursion: true},
		{tag: "__", openElement: typeBoldOpen, closeElement: typeBoldClose, recursion: true},
		{tag: "~~", openElement: typeCancellationOpen, closeElement: typeCancellationClose, recursion: true},
		{tag: "*", openElement: typeItalicOpen, closeElement: typeItalicClose, recursion: true},
		{tag: "_", openElement: typeItalicOpen, closeElement: typeItalicClose, recursion: true},
	}

	// テキストの処理（高速化できそう）
	for i := 0; i < len(text); i++ {
		for _, d := range dictionary {
			isInline, after := lexerInlineStart(text[i:], d.tag)
			if isInline {
				token = append(token, Token{
					content:     text[:i],
					elementType: typeText,
				})
				token = lexerInline(token, text[i+len(d.tag):])
				text = after
				break
			}
		}
	}

	// 最終的に残ったテキストをプレーンテキストとして追加
	return append(token, Token{
		content:     text,
		elementType: typeText,
	})
}

// lexerInlineLink ...リンク関連の解析
func lexerInlineLink(token []Token, text string) ([]Token, string) {
	token, text = lexerInlineLinkN(token, text, "![", "](", ")")
	token, text = lexerInlineLinkN(token, text, "[", "](", ")")
	return token, text
}

// lexerInlineLinkN ...リンク関連の解析（内部処理）
func lexerInlineLinkN(token []Token, text string, start string, mid string, end string) ([]Token, string) {

	startN := strings.Index(text, start)
	if startN == -1 {
		return token, text
	}

	midN := strings.Index(text[startN:], mid)
	if midN == -1 {
		return token, text
	}
	midN += startN - 1

	endN := strings.Index(text[midN:], end)
	if endN == -1 {
		return token, text
	}
	endN += midN - 1

	token = append(token, Token{
		content:     text[startN+len(start)-1 : midN],
		elementType: typeLink,
	})
	token = lexerInline(token, text[midN+len(mid)-1:endN])
	token = append(token, Token{elementType: typeLinkClose})

	return token, text[endN+len(end)-1:]
}

// lexerInlineStart ...インラインの開始であるかをチェック（終了が存在しないなら開始扱いしない）
func lexerInlineStart(text string, tag string) (bool, string) {
	// そもそも違う
	if len(text) < len(tag)*2 || text[:len(tag)] != tag {
		return false, text
	}

	// 終了が存在するか？
	end := strings.Index(text[len(tag):], tag)
	if end == -1 {
		return false, text
	}

	after := text[end+len(tag):]
	return true, after
}

// lexerCloseBlock ...開いているブロックをすべて閉じる
func lexerCloseBlock(token []Token, options Options) ([]Token, Options) {

	if options.isCodeBlock {
		token = append(token, Token{elementType: typeBlockCodeClose})
		options.isCodeBlock = false
	}

	if options.isParagraph {
		token = append(token, Token{elementType: typeParagraphClose})
		options.isParagraph = false
	}

	for options.isListBlock > 0 {
		token = append(token, Token{elementType: typeListBlockClose})
		options.isListBlock--
	}

	return token, options
}

// lexerHeaderBlock ...ヘッダーであるかをチェックする（1以上でヘッダー）
func lexerHeaderBlock(line string) (int, string) {
	for i := 4; i >= 1; i-- {
		if len(line) > i && line[:i] == ("#####"[:i-1]+" ") {
			return i - 1, line[i:]
		}
	}
	return 0, line
}

// lexerListBlock ...リストブロックであるかをチェックする（階層1以上でリストブロック）
func lexerListBlock(line string) (int, string) {
	n := strings.Index(line, "- ")
	if n == -1 {
		return 0, line
	}
	return min(n/2+1, 3), line[n+1:]
}

// lexerCodeBlock ...コードブロックであるかをチェックする
func lexerCodeBlock(line string) bool {
	line = strings.Replace(line, " ", "", -1)
	return len(line) >= 3 && line[:3] == "```"
}

// lexerHorizon ...水平線をチェックする
func lexerHorizon(line string) bool {
	line = strings.Replace(line, " ", "", -1)
	return len(line) >= 3 && (len(line) == strings.Count(line, "-") || len(line) == strings.Count(line, "_") || len(line) == strings.Count(line, "*"))
}

// lexerEmpty ...空行をチェックする
func lexerEmpty(line string) bool {
	return strings.Replace(line, " ", "", -1) == ""
}
