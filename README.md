# white600
Markdown を HTML に変換します。

## 使い方
### インストール
```
go get -u github.com/TomSuzuki/white600
```

### 使い方
```
md, _ = ioutil.ReadFile("./markdown.md")
html := white600.MarkdownToHTML(string(md))
```

