package goldmark_test

import (
	"testing"

	"github.com/yuin/goldmark/v2/extension"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/testutil"
)

func TestEscapedSpace(t *testing.T) {
	markdown := testutil.NewMarkdownToStringFunc(
		parser.New(),
		html.New(
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
	no := 1
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Without spaces around an emphasis started with east asian punctuations, it is not interpreted as an emphasis(as defined in CommonMark spec)",
			Markdown:    "太郎は**「こんにちわ」**と言った\nんです",
			Expected:    "<p>太郎は**「こんにちわ」**と言った\nんです</p>",
		},
		t,
	)

	no = 2
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "With spaces around an emphasis started with east asian punctuations, it is interpreted as an emphasis(but remains unnecessary spaces)",
			Markdown:    "太郎は **「こんにちわ」** と言った\nんです",
			Expected:    "<p>太郎は <strong>「こんにちわ」</strong> と言った\nんです</p>",
		},
		t,
	)

	// EscapedSpace
	markdown = testutil.NewMarkdownToStringFunc(
		parser.New(
			parser.WithEscapedSpace(),
		),
		html.New(
			html.WithXHTML(),
			html.WithUnsafe(),
			html.WithLineBreakStrategy(html.SimpleEastAsianLineBreakStrategy),
		),
	)

	no = 3
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "With spaces around an emphasis started with east asian punctuations,it is interpreted as an emphasis",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\nんです",
			Expected:    "<p>太郎は<strong>「こんにちわ」</strong>と言ったんです</p>",
		},
		t,
	)

	// ' ' triggers Linkify extension inline parser.
	// Escaped spaces should not trigger the inline parser.

	markdown = testutil.NewMarkdownToStringFunc(
		parser.New(
			parser.WithExtensions(
				extension.NewLinkifyParser(),
			),
			parser.WithEscapedSpace(),
		),
		html.New(
			html.WithXHTML(),
			html.WithUnsafe(),
			html.WithLineBreakStrategy(html.SimpleEastAsianLineBreakStrategy),
		),
	)

	no = 4
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Escaped space and linkfy extension",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\nんです",
			Expected:    "<p>太郎は<strong>「こんにちわ」</strong>と言ったんです</p>",
		},
		t,
	)
}

func TestEastAsianLineBreaks(t *testing.T) {
	markdown := testutil.NewMarkdownToStringFunc(
		parser.New(),
		html.New(
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
	no := 1
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{

			No:          no,
			Description: "Soft line breaks are rendered as a newline, so some asian users will see it as an unnecessary space",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\nんです",
			Expected:    "<p>太郎は\\ <strong>「こんにちわ」</strong>\\ と言った\nんです</p>",
		},
		t,
	)

	// Enables EastAsianLineBreaks
	markdown = testutil.NewMarkdownToStringFunc(
		parser.New(
			parser.WithEscapedSpace(),
		),
		html.New(
			html.WithXHTML(),
			html.WithUnsafe(),
			html.WithLineBreakStrategy(html.SimpleEastAsianLineBreakStrategy),
		),
	)

	no = 2
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between east asian wide characters are ignored",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\nんです",
			Expected:    "<p>太郎は<strong>「こんにちわ」</strong>と言ったんです</p>",
		},
		t,
	)

	no = 3
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between east asian wide characters(w/ strong) are ignored",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\n**ん**です",
			Expected:    "<p>太郎は<strong>「こんにちわ」</strong>と言った<strong>ん</strong>です</p>",
		},
		t,
	)

	no = 4
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between east asian wide characters(w/ link) are ignored",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\n[ん](./test)です",
			Expected:    "<p>太郎は<strong>「こんにちわ」</strong>と言った<a href=\"./test\">ん</a>です</p>",
		},
		t,
	)

	no = 5
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between east asian wide characters(w/ leading spaces) are ignored",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\n んです",
			Expected:    "<p>太郎は<strong>「こんにちわ」</strong>と言ったんです</p>",
		},
		t,
	)

	no = 6
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between east asian wide characters(w/ codespan) are ignored",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\n`ん`です",
			Expected:    "<p>太郎は<strong>「こんにちわ」</strong>と言った<code>ん</code>です</p>",
		},
		t,
	)

	no = 7
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between east asian wide characters(w/ codespans) are ignored",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言a`た`\n`ん`です",
			Expected:    "<p>太郎は<strong>「こんにちわ」</strong>と言a<code>た</code><code>ん</code>です</p>",
		},
		t,
	)

	no = 8
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between ascii characters(w/ codespan) are not ignored",
			Markdown:    "Access is enforced in\n`internal/app/server.go`.",
			Expected:    "<p>Access is enforced in\n<code>internal/app/server.go</code>.</p>",
		},
		t,
	)

	no = 9
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between east asian wide characters(w/ numeric references) are ignored",
			Markdown:    "&#x65e5;\n&#x672c;",
			Expected:    "<p>日本</p>",
		},
		t,
	)

	no = 10
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between western characters are rendered as a newline",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言ったa\nbんです",
			Expected:    "<p>太郎は<strong>「こんにちわ」</strong>と言ったa\nbんです</p>",
		},
		t,
	)

	no = 11
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between a western character and an east asian wide character are rendered as a newline",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言ったa\nんです",
			Expected:    "<p>太郎は<strong>「こんにちわ」</strong>と言ったa\nんです</p>",
		},
		t,
	)

	no = 12
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between an east asian wide character and a western character are rendered as a newline",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\nbんです",
			Expected:    "<p>太郎は<strong>「こんにちわ」</strong>と言った\nbんです</p>",
		},
		t,
	)

	// WithHardWraps take precedence over WithEastAsianLineBreaks
	markdown = testutil.NewMarkdownToStringFunc(
		parser.New(
			parser.WithEscapedSpace(),
		),
		html.New(
			html.WithXHTML(),
			html.WithUnsafe(),
			html.WithHardWraps(),
			html.WithLineBreakStrategy(html.SimpleEastAsianLineBreakStrategy),
		),
	)

	no = 13
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "WithHardWraps take precedence over WithEastAsianLineBreaks",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\nんです",
			Expected:    "<p>太郎は<strong>「こんにちわ」</strong>と言った<br />\nんです</p>",
		},
		t,
	)

	markdown = testutil.NewMarkdownToStringFunc(
		parser.New(
			parser.WithExtensions(extension.NewLinkifyParser()),
		),
		html.New(
			html.WithXHTML(),
			html.WithUnsafe(),
			html.WithLineBreakStrategy(html.SimpleEastAsianLineBreakStrategy),
		),
	)
	no = 14
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "WithEastAsianLineBreaks and linkfy extension",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\r\nんです",
			Expected:    "<p>太郎は\\ <strong>「こんにちわ」</strong>\\ と言ったんです</p>",
		},
		t,
	)
	no = 15
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between east asian wide characters or punctuations are ignored",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と、\r\n言った\r\nんです",
			Expected:    "<p>太郎は\\ <strong>「こんにちわ」</strong>\\ と、言ったんです</p>",
		},
		t,
	)
	no = 16
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between an east asian wide character and a western character are ignored",
			Markdown:    "私はプログラマーです。\n東京の会社に勤めています。\nGoでWebアプリケーションを開発しています。",
			Expected:    "<p>私はプログラマーです。東京の会社に勤めています。\nGoでWebアプリケーションを開発しています。</p>",
		},
		t,
	)

	// Tests with EastAsianLineBreaksCSS3Draft
	markdown = testutil.NewMarkdownToStringFunc(
		parser.New(),
		html.New(
			html.WithXHTML(),
			html.WithUnsafe(),
			html.WithLineBreakStrategy(html.CSSText3LineBreakStrategy),
		),
	)
	no = 17
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between a western character and an east asian wide character are ignored",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言ったa\nんです",
			Expected:    "<p>太郎は\\ <strong>「こんにちわ」</strong>\\ と言ったaんです</p>",
		},
		t,
	)

	no = 18
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between an east asian wide character and a western character are ignored",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\nbんです",
			Expected:    "<p>太郎は\\ <strong>「こんにちわ」</strong>\\ と言ったbんです</p>",
		},
		t,
	)

	no = 19
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between an east asian wide character and a western character are ignored",
			Markdown:    "私はプログラマーです。\n東京の会社に勤めています。\nGoでWebアプリケーションを開発しています。",
			Expected:    "<p>私はプログラマーです。東京の会社に勤めています。GoでWebアプリケーションを開発しています。</p>",
		},
		t,
	)

	no = 20
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks around western punctuation are converted to spaces",
			Markdown:    "This is one sentence.\nThis is another one,\nwith a comma.",
			Expected:    "<p>This is one sentence.\nThis is another one,\nwith a comma.</p>",
		},
		t,
	)

	no = 21
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks around inline containers ending or starting with CodeSpan use the code span runes",
			Markdown:    "commit\n[`e7d50326de51`](https://example.com)的\n\n的[`e7d50326de51`](https://example.com)\ncommit",
			Expected:    "<p>commit\n<a href=\"https://example.com\"><code>e7d50326de51</code></a>的</p>\n<p>的<a href=\"https://example.com\"><code>e7d50326de51</code></a>\ncommit</p>",
		},
		t,
	)

	no = 22
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks around nested inline containers wrapping CodeSpan inspect innermost inline rune",
			Markdown:    "foo\n[*`code`*](https://example.com)的\n\n的[*`code`*](https://example.com)\nbar",
			Expected:    "<p>foo\n<a href=\"https://example.com\"><em><code>code</code></em></a>的</p>\n<p>的<a href=\"https://example.com\"><em><code>code</code></em></a>\nbar</p>",
		},
		t,
	)

	no = 23
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks around AutoLink inspect AutoLink label",
			Markdown:    "See\n<https://example.com>的\n\n的<https://example.com>\nagain",
			Expected:    "<p>See\n<a href=\"https://example.com\">https://example.com</a>的</p>\n<p>的<a href=\"https://example.com\">https://example.com</a>\nagain</p>",
		},
		t,
	)

	no = 24
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks around Image treat Image as replaced object rather than its alt text",
			Markdown:    "See\n![图](a.png)\nagain\n\nSee\n![](a.png)的\n\n的![](a.png)\nagain",
			Expected:    "<p>See\n<img src=\"a.png\" alt=\"图\" />\nagain</p>\n<p>See\n<img src=\"a.png\" alt=\"\" />的</p>\n<p>的<img src=\"a.png\" alt=\"\" />\nagain</p>",
		},
		t,
	)

	markdown = testutil.NewMarkdownToStringFunc(
		parser.New(),
		html.New(
			html.WithXHTML(),
			html.WithUnsafe(),
			html.WithLineBreakStrategy(html.SimpleEastAsianLineBreakStrategy),
		),
	)

	no = 25
	testutil.DoTestCase(
		markdown,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks around Image under SimpleEastAsianLineBreakStrategy treat Image as replaced object",
			Markdown:    "中文\n![图](a.png)\n中文",
			Expected:    "<p>中文\n<img src=\"a.png\" alt=\"图\" />\n中文</p>",
		},
		t,
	)

}
