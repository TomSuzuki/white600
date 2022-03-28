# white600
Markdown を HTML に変換します。

## 使い方
### インストール
```
go get -u github.com/TomSuzuki/gomarkdown
```

### 使い方
```
md, _ = ioutil.ReadFile("./markdown.md")
html := gomarkdown.MarkdownToHTML(string(md))
```

