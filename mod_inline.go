package white600

import (
	"strings"
)

type InlineType int

const (
	inlineText InlineType = iota
	inlineCode
	inlineBold
	inlineItalic
	inlineCancel
	inlineLink
)

type inlineObject struct {
	inlineType InlineType
	content    string
	child      *[]inlineObject
}

// inlineConv ...インライン要素の解析
func (data *MarkdownInfo) inlineConv(text string) string {

	// 初期値の生成
	inline := []inlineObject{
		{
			inlineType: inlineText,
			content:    text,
			child:      nil,
		},
	}

	// リンクの解析
	inline = inlineLinkX(inline, [3]string{"![", "](", ")"}, "<img alt='$1' src='$2'>")
	inline = inlineLinkX(inline, [3]string{"[", "](", ")"}, "<a href='$2'>$1</a>")

	// 再帰的にその他のインライン要素の解析
	inline = inlineTag(inline)

	// // <img> and <a>
	// data.inlineLink(inline, "![", "](", ")", "<img alt='$1' src='$2'>")
	// data.inlineLink(inline, "[", "](", ")", "<a href='$2'>$1</a>")

	// // inline text
	// data.inlineTag(inline, "**", "strong")
	// data.inlineTag(inline, "`", "code")
	// data.inlineTag(inline, "~", "s")
	// data.inlineTag(inline, "__", "em")
	// data.inlineTag(inline, "_", "em")
	// data.inlineTag(inline, "*", "em")

	// <br>
	//todo data.markdownLines[0] = strings.Replace(data.markdownLines[0], "  ", "<br>", -1)

	// 解析結果を文字列に
	return inlineConvert(inline)
}

// inlineLink ...インラインのリンクを解析（リンク内にリンクは作れないので再帰しない、リンク要素は前優先ではなく、字句優先）
func inlineLinkX(inline []inlineObject, mdTemplate [3]string, htmlTemplate string) []inlineObject {

	// 全てのテキスト要素をチェック
	isLinkGenerated := true
	for isLinkGenerated {
		isLinkGenerated = false

		for i, v := range inline {

			// テキスト要素以外は処理しない
			if v.inlineType != inlineText {
				continue
			}

			// 分解して解析
			var mdPoint [4]int // 先頭の参照先が無いので1個ずらす
			for j, w := range mdTemplate {
				mdPoint[j+1] = strings.Index(v.content[mdPoint[j]:], w)
				if mdPoint[j+1] == -1 {
					break
				}
				mdPoint[j+1] += mdPoint[j]
			}

			// 変換が必要ない場合は処理しない
			if mdPoint[1] == -1 || mdPoint[2] == -1 || mdPoint[3] == -1 {
				continue
			}

			// インライン要素の解析
			alt := inlineTag([]inlineObject{{inlineType: inlineText, content: v.content[mdPoint[1]+len(mdTemplate[0]) : mdPoint[2]], child: nil}})

			// 置き換えテキストを作成
			s := strings.NewReplacer(
				"$1", inlineConvert(alt),
				"$2", v.content[mdPoint[2]+len(mdTemplate[1]):mdPoint[3]],
			).Replace(htmlTemplate)

			// オブジェクトを更新
			isLinkGenerated = true
			nextInline := append(inline[:i], inlineObject{inlineType: inlineText, content: v.content[:mdPoint[1]]})
			nextInline = append(nextInline, inlineObject{inlineType: inlineLink, content: s})
			nextInline = append(nextInline, inlineObject{inlineType: inlineText, content: v.content[len(mdTemplate[2])+mdPoint[3]:]})
			nextInline = append(nextInline, inline[i+1:]...)
			inline = nextInline

			break
		}
	}

	return inline
}

// inlineConvert ...インライン要素の解析結果を文字列に変換
func inlineConvert(inline []inlineObject) string {
	text := []string{}

	for _, v := range inline {
		// 開きタグを入れる
		text = append(text, inlineTagOpen(v.inlineType))

		// 内部テキストを入れる
		text = append(text, v.content)

		// 子要素を処理
		if v.child != nil {
			text = append(text, inlineConvert(*v.child))
		}

		// 閉じタグを入れる
		text = append(text, inlineTagClose(v.inlineType))

	}

	// まとめる
	return strings.Join(text, "")
}

// inlineTagClose ...閉じタグ
func inlineTagClose(inlineType InlineType) string {
	switch inlineType {
	case inlineCode:
		return "</code>"
	case inlineBold:
		return "</strong>"
	case inlineItalic:
		return "</em>"
	case inlineCancel:
		return "</s>"
	default:
		return ""
	}
}

// inlineTagOpen ...開きタグ
func inlineTagOpen(inlineType InlineType) string {
	switch inlineType {
	case inlineCode:
		return "<code>"
	case inlineBold:
		return "<strong>"
	case inlineItalic:
		return "<em>"
	case inlineCancel:
		return "<s>"
	default:
		return ""
	}
}

// inlineTag ...インラインの要素を解析
func inlineTag(inline []inlineObject) []inlineObject {

	// 辞書用の型
	type dictInfo struct {
		tag        string
		inlineType InlineType
		recursion  bool
	}

	// テキストのパターンと対応するタイプ（場所移動する）
	dictionary := []dictInfo{
		{tag: "`", inlineType: inlineCode, recursion: false},
		{tag: "**", inlineType: inlineBold, recursion: true},
		{tag: "__", inlineType: inlineBold, recursion: true},
		{tag: "~~", inlineType: inlineCancel, recursion: true},
		{tag: "*", inlineType: inlineItalic, recursion: true},
		{tag: "_", inlineType: inlineItalic, recursion: true},
	}

	isGenerated := true
	for isGenerated {
		isGenerated = false

		for i, v := range inline {

			// 解析しない物を除外
			if v.inlineType == inlineLink || v.inlineType == inlineCode {
				continue
			}

			// 現在の要素を解析
			t := v.content
			for j := 0; j < len(t); j++ {
				for _, x := range dictionary {
					flg, ct, af := lexerInlineStart(v.content[j:], x.tag)
					bf := v.content[:j]
					if !flg {
						continue
					}

					// 中のテキストを処理
					chobj := []inlineObject{{inlineType: inlineText, content: ct}}
					if x.recursion {
						chobj = inlineTag(chobj)
					}

					// オブジェクトを更新
					isGenerated = true
					nextInline := append(inline[:i], inlineObject{inlineType: x.inlineType, content: bf, child: &chobj})
					nextInline = append(nextInline, inlineObject{inlineType: inlineText, content: af})
					nextInline = append(nextInline, inline[i+1:]...)
					inline = nextInline

					break
				}

				if isGenerated {
					break
				}
			}

			if isGenerated {
				break
			}
		}

	}

	return inline
}

// lexerInlineStart ...インラインの開始であるかをチェック（終了が存在しないなら開始扱いしない）
func lexerInlineStart(text string, tag string) (bool, string, string) {
	// そもそも違う
	if len(text) < len(tag)*2 || text[:len(tag)] != tag {
		return false, text, ""
	}

	// 終了が存在するか？
	end := strings.Index(text[len(tag):], tag)
	if end == -1 {
		return false, text, ""
	}

	inner := text[len(tag) : end+len(tag)]
	after := text[end+len(tag)*2:]
	return true, inner, after
}

// ...インライン要素のタグを解析

// // inlineLink ...![]() and []()
// func (data *MarkdownInfo) inlineLink(start, middle, end string, format string) {
// 	var line = data.currentData.currentLine
// 	var p = 0

// 	for true {
// 		// index
// 		var iMiddle = p + strings.Index(line[p:], middle)
// 		if iMiddle == -1 {
// 			break
// 		}
// 		var iStart = strings.LastIndex(line[:iMiddle], start)
// 		var iEnd = iMiddle + strings.Index(line[iMiddle:], end)

// 		// replace
// 		if iMiddle != -1 && iStart != -1 && iEnd != -1 {
// 			var s = strings.NewReplacer(
// 				"$1", line[len(start)+iStart:iMiddle],
// 				"$2", line[len(middle)+iMiddle:iEnd],
// 			).Replace(format)
// 			line = line[:iStart] + s + line[len(end)+iEnd:]
// 		}

// 		// error check
// 		p = iMiddle + len(middle)
// 		if p > len(line) {
// 			break
// 		}
// 	}

// 	//todo data.markdownLines[0] = line
// }

// // indexList
// func indexList(s, substr string) []int {
// 	var n []int
// 	for true {
// 		var m = strings.Index(s, substr)
// 		if m == -1 {
// 			break
// 		}
// 		n = append(n, m)
// 		s = s[m:]
// 	}
// 	return n
// }

// // inlineTag ...md -> html
// func (data *MarkdownInfo) inlineTag(md string, html string) {
// 	var codeList = strings.Split(data.currentData.currentLine, md)
// 	var isEven = len(codeList)%2 == 0
// 	var text []string
// 	data.currentData.currentLine = ""

// 	// insert tags
// 	for i, v := range codeList {
// 		if isEven && i == len(codeList)-1 {
// 			text = append(text, md)
// 			text = append(text, v)
// 		} else if i%2 == 0 {
// 			text = append(text, v)
// 		} else if isNotBrokenHTML(v) {
// 			text = append(text, "<")
// 			text = append(text, html)
// 			text = append(text, ">")
// 			text = append(text, v)
// 			text = append(text, "</")
// 			text = append(text, html)
// 			text = append(text, ">")
// 		} else {
// 			text = append(text, md)
// 			text = append(text, v)
// 		}
// 	}

// 	// join
// 	//todo data.markdownLines[0] = strings.Join(text, "")
// }

// // isNotBrokenHTML ..."<s></s><em></em>" <<< true, "<s><em></em>" <<< false, "<img ...>" <<< false...?
// func isNotBrokenHTML(html string) bool {
// 	var nest = 0
// 	var open = false
// 	for i := 0; i < len(html); i++ {
// 		s := html[i]
// 		if open {
// 			if string(s) == "/" {
// 				nest--
// 			} else {
// 				nest++
// 			}
// 		}
// 		open = string(s) == "<"
// 		if nest == -1 {
// 			return false
// 		}
// 	}
// 	return nest == 0
// }
