package html_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/extension"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/text"
	"github.com/yuin/goldmark/v2/util"
)

func render(t *testing.T, r html.Renderer, source string) string {
	t.Helper()
	var buf bytes.Buffer
	src := []byte(source)
	if err := r.Render(&buf, src, parser.New().Parse(src)); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

func TestWithNodeRendererOverridesCommonMark(t *testing.T) {
	strong := html.NodeRendererFunc(func(
		w io.Writer, _ []byte, _ ast.Node, entering bool, _ renderer.Context,
	) (ast.WalkStatus, error) {
		if entering {
			_, _ = io.WriteString(w, `<b class="important">`)
		} else {
			_, _ = io.WriteString(w, `</b>`)
		}
		return ast.WalkContinue, nil
	})
	r := html.New(html.WithNodeRenderer(ast.KindStrong, strong))
	got := render(t, r, "This is **important**.\n")
	want := "<p>This is <b class=\"important\">important</b>.</p>\n"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestWithNodeRenderersOverridesCommonMark(t *testing.T) {
	emphasis := html.NodeRendererFunc(func(
		w io.Writer, _ []byte, _ ast.Node, entering bool, _ renderer.Context,
	) (ast.WalkStatus, error) {
		if entering {
			_, _ = io.WriteString(w, `<i class="important">`)
		} else {
			_, _ = io.WriteString(w, `</i>`)
		}
		return ast.WalkContinue, nil
	})
	r := html.New(html.WithNodeRenderers(map[ast.NodeKind]html.NodeRenderer{
		ast.KindEmphasis: emphasis,
	}))
	got := render(t, r, "*important*\n")
	want := "<p><i class=\"important\">important</i></p>\n"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestWithNodeRendererOverridesExtension(t *testing.T) {
	paragraph := html.NodeRendererFunc(func(
		w io.Writer, _ []byte, _ ast.Node, entering bool, _ renderer.Context,
	) (ast.WalkStatus, error) {
		if entering {
			_, _ = io.WriteString(w, `<div>`)
		} else {
			_, _ = io.WriteString(w, `</div>`)
		}
		return ast.WalkContinue, nil
	})
	r := html.New(
		html.WithExtensions(extension.NewTaskListItemHTMLRenderer()),
		html.WithNodeRenderer(ast.KindParagraph, paragraph),
	)
	got := render(t, r, "- [ ] todo\n")
	want := "<ul>\n<li><div>[ ] todo</div></li>\n</ul>\n"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestWithNodeRendererDecorator(t *testing.T) {
	decorator := func(next html.NodeRenderer) html.NodeRenderer {
		return html.NodeRendererFunc(func(
			w io.Writer, source []byte, n ast.Node, entering bool, rc renderer.Context,
		) (ast.WalkStatus, error) {
			if entering {
				_, _ = io.WriteString(w, `<!-- decorated -->`)
			}
			return next.Render(w, source, n, entering, rc)
		})
	}
	r := html.New(html.WithNodeRendererDecorator(ast.KindStrong, decorator))
	got := render(t, r, "**important**\n")
	want := "<p><!-- decorated --><strong>important</strong></p>\n"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

type readerOnly struct {
	io.Reader
}

func TestContextWritersSatisfyIOWriterContract(t *testing.T) {
	const value = "a < b && c > d \"q\" \x00 あ"
	copied := 0
	codeSpan := html.NodeRendererFunc(func(
		_ io.Writer, _ []byte, _ ast.Node, entering bool, rc renderer.Context,
	) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		writers := []util.BufWriter{
			html.ContextTextWriter(rc),
			html.ContextHTMLWriter(rc),
			html.ContextLinkURLWriter(rc),
		}
		for _, cw := range writers {
			for _, r := range []io.Reader{
				strings.NewReader(value),
				readerOnly{strings.NewReader(value)},
			} {
				n, err := io.Copy(cw, r)
				if err != nil {
					return ast.WalkStop, err
				}
				if n != int64(len(value)) {
					return ast.WalkStop, fmt.Errorf("wrote %d bytes, want %d", n, len(value))
				}
				copied++
			}
		}
		return ast.WalkSkipChildren, nil
	})

	source := []byte("`x`\n")
	r := html.New(html.WithNodeRenderer(ast.KindCodeSpan, codeSpan))
	if err := r.Render(io.Discard, source, parser.New().Parse(source)); err != nil {
		t.Fatal(err)
	}
	if copied != 6 {
		t.Fatalf("io.Copy called %d times, want 6", copied)
	}
}

func TestSoftLineBreakAdjacentRuneScanIsBounded(t *testing.T) {
	r := html.New(html.WithLineBreakStrategy(html.SimpleEastAsianLineBreakStrategy))

	// A few invisible inline nodes between East Asian runes: the soft line
	// breaks look through them and are removed.
	links := strings.Repeat(`[](u)`, 4)
	want := "<p>中" + strings.Repeat(`<a href="u"></a>`, 4) + "文</p>\n"
	if got := render(t, r, "中\n"+links+"\n文"); got != want {
		t.Errorf("short run:\nwant %q\ngot  %q", want, got)
	}

	// A long run of invisible inline nodes exceeds the scan budget, so the
	// adjacent runes are unknown and the soft line breaks are kept. Without the
	// budget every soft line break rescans the run, which is quadratic.
	links = strings.Repeat(`[](u)`, 64)
	want = "<p>中\n" + strings.Repeat(`<a href="u"></a>`, 64) + "\n文</p>\n"
	if got := render(t, r, "中\n"+links+"\n文"); got != want {
		t.Errorf("long run:\nwant %q\ngot  %q", want, got)
	}

	renderNested := func(depth int) string {
		doc := ast.NewDocument()
		p := ast.NewParagraph()
		doc.AppendChild(p)
		p.AppendChild(ast.NewText(text.NewSingleLineValueFromString("中", nil)))
		var cur ast.Node = p
		for range depth {
			em := ast.NewEmphasis()
			cur.AppendChild(em)
			cur = em
		}
		sb := ast.NewText(text.SingleLineValue{})
		sb.SetSoftLineBreak(true)
		cur.AppendChild(sb)
		p.AppendChild(ast.NewText(text.NewSingleLineValueFromString("文", nil)))

		var buf bytes.Buffer
		if err := r.Render(&buf, nil, doc); err != nil {
			t.Fatal(err)
		}
		return buf.String()
	}

	// Shallow parent climbs look through enclosing inline containers.
	want = "<p>中" + strings.Repeat("<em>", 4) + strings.Repeat("</em>", 4) + "文</p>\n"
	if got := renderNested(4); got != want {
		t.Errorf("shallow parent climb:\nwant %q\ngot  %q", want, got)
	}

	// Deep parent climbs also consume the scan budget and fall back to keeping
	// the soft line break.
	want = "<p>中" + strings.Repeat("<em>", 64) + "\n" + strings.Repeat("</em>", 64) + "文</p>\n"
	if got := renderNested(64); got != want {
		t.Errorf("deep parent climb:\nwant %q\ngot  %q", want, got)
	}
}
