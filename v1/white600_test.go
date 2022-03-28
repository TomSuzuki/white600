/*
 * white600_test
 * available at https://github.com/TomSuzuki/white600/
 *
 * Copyright 2022 TomSuzuki
 * LICENSE: MIT
 *
 * # How to use
 * > gotest -run NONE -bench .
 * > gotest -v
 */

//package white600
package white600_test

import (
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/TomSuzuki/white600"
)

type testFile struct {
	markdown string
	html     string
}

func Test(t *testing.T) {
	// test case (markdown, html)
	dir := "./testcase/"
	var testfile []testFile

	// get list
	files, _ := ioutil.ReadDir("./testcase/")
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".md" {
			testfile = append(testfile, testFile{
				markdown: dir + file.Name(),
				html:     dir + file.Name() + ".html",
			})
		}
	}

	// test
	for i := range testfile {
		test(testfile[i], t)
	}
}

// speed test
func BenchmarkSpeed_white600(b *testing.B) {
	file := "./testcase/00.md"
	md, _ := ioutil.ReadFile(file)

	b.ResetTimer()
	for ct := 0; ct < 1500; ct++ {
		white600.MarkdownToHTML(string(md))
	}
}

func test(test testFile, t *testing.T) {
	// load
	b, _ := ioutil.ReadFile(test.html)
	sample := string(b)
	b, _ = ioutil.ReadFile(test.markdown)
	answer := string(b)

	// html -> markdown
	answer = white600.MarkdownToHTML(answer)

	// trim
	sample = strings.NewReplacer("\r\n", "", "\r", "", "\n", "", " ", "", "'", "\"").Replace(sample)
	answer = strings.NewReplacer("\r\n", "", "\r", "", "\n", "", " ", "", "'", "\"").Replace(answer)

	// html
	sampleHTML := template.HTML(sample)
	answerHTML := template.HTML(answer)

	// check
	if sampleHTML != answerHTML {
		t.Logf("☒  failed test: \t%s", test.markdown)
		t.Logf(" - sample: %s", sampleHTML)
		t.Logf(" - answer: %s", answerHTML)
		t.Logf("")
		t.Fail()
	} else {
		t.Logf("☑  success test: \t%s", test.markdown)
	}
}
