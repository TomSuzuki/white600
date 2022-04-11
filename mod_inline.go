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
	inlineBreak
)

type inlineObject struct {
	inlineType InlineType
	content    string
	child      *[]inlineObject
}

// inlineConv ...インライン要素の解析
func (data *MarkdownInfo) inlineConv(text string) string {

	// 初期値の生成
	inline := []inlineObject{{inlineType: inlineText, content: text}}

	// リンクの解析（画像リンクのためにデータを文字列に変換する）
	inline = inlineLinkX(inline, [3]string{"![", "](", ")"}, "<img alt='$1' src='$2'>")
	inline = inlineLinkX([]inlineObject{{inlineType: inlineText, content: inlineConvert(inline)}}, [3]string{"[", "](", ")"}, "<a href='$2'>$1</a>")

	// 再帰的にその他のインライン要素の解析
	inline = inlineTag(inline)

	// 解析結果を文字列に
	return inlineConvert(inline)
}

// inlineLink ...インラインのリンクを解析（リンク要素は前優先ではなく字句優先）
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

// 辞書用の型
type dictInfo struct {
	tag        string
	inlineType InlineType
	recursion  bool
}

// テキストのパターンと対応するタイプ（場所移動する）
var dictionary = []dictInfo{
	{tag: "`", inlineType: inlineCode, recursion: false},
	{tag: "**", inlineType: inlineBold, recursion: true},
	{tag: "__", inlineType: inlineBold, recursion: true},
	{tag: "~~", inlineType: inlineCancel, recursion: true},
	{tag: "*", inlineType: inlineItalic, recursion: true},
	{tag: "_", inlineType: inlineItalic, recursion: true},
}

// inlineTag ...インラインの要素を解析
func inlineTag(inline []inlineObject) []inlineObject {

	for i := 0; i < len(inline); i++ {
		v := inline[i]

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

				// 要素を間に追加
				nextInline := append(inline[:i], inlineObject{inlineType: inlineText, content: bf})
				nextInline = append(nextInline, inlineObject{inlineType: x.inlineType, content: "", child: &chobj})
				nextInline = append(nextInline, inlineObject{inlineType: inlineText, content: af})
				nextInline = append(nextInline, inline[i+1:]...)
				inline = nextInline

				// 現在の要素の解析を抜ける
				j = len(t)

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
