package template_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/lpar/template"
)

func testRender(rdr *template.Renderer, tmplname string, data interface{}) (string, error) {
	var result strings.Builder
	err := rdr.ExecuteTemplate(&result, tmplname, data)
	return result.String(), err
}

const expectedButton = `<button class=ds-button>Submit</button>`
const expectedPage = `<!doctype html><html lang=en><meta charset=utf-8><title>Hello world</title><button class=ds-button>Submit</button>`
const expectedCSS = `button{border:solid red 3px;min-width:10em}`
const expectedJS = `function nth(o){return o+(['st','nd','rd'][(o+'').match(/1?\d\b/)-1]||'th');}`
const expectedSafety = `<p>Written by &lt;script&gt;alert(&#39;You have been pwned&#39;);&lt;/script&gt;.`

func TestRenderHTML(t *testing.T) {
	rdr := template.NewRenderer()
	rdr.ParseFiles("testdata")

	btn, err := testRender(rdr, "subdir/button.html", nil)
	if err != nil {
		t.Errorf("button.html render failed: %v", err)
	}
	if btn != expectedButton {
		t.Errorf("button.html render failed, expected '%s' got '%s'", expectedButton, btn)
	}

	html, err := testRender(rdr, "index.html", "Hello world")
	if err != nil {
		t.Errorf("index.html render failed: %v", err)
	}
	if html != expectedPage {
		t.Errorf("index.html render failed, expected '%s' got '%s'", expectedPage, html)
	}
}

func TestRenderSafety(t *testing.T) {
	rdr := template.NewRenderer()
	rdr.ParseFiles("testdata")
	html, err := testRender(rdr, "footer.html", "<script>alert('You have been pwned');</script>")
	if err != nil {
		t.Errorf("footer.html render with injected tags failed: %v", err)
	}
	if html != expectedSafety {
		t.Errorf("footer.html render with injected tags failed, expected '%s' got '%s'", expectedSafety, html)
	}
}

func TestRenderCSS(t *testing.T) {
	rdr := template.NewRenderer()
	rdr.ParseFiles("testdata")

	css, err := testRender(rdr, "button.css", nil)
	if err != nil {
		t.Errorf("button.css render failed: %v", err)
	}
	if css != expectedCSS {
		t.Errorf("button.css render failed, expected '%s' got '%s'", expectedCSS, css)
	}
}

func TestRenderJS(t *testing.T) {
	rdr := template.NewRenderer()
	rdr.ParseFiles("testdata")

	js, err := testRender(rdr, "nth.js", nil)
	if err != nil {
		t.Errorf("nth.js render failed: %v", err)
	}
	if js != expectedJS {
		t.Errorf("nth.js render failed, expected '%s' got '%s'", expectedJS, js)
	}
}

const TMP1 = "<button class=\"btn ds-button\">Primary</button>"
const TMP2 = "<a class=\"btn ds-button\">Secondary</a>"

func TestReload(t *testing.T) {
	err := ioutil.WriteFile("testdata/reload.html", []byte(TMP1), 0644)
	if err != nil {
		t.Errorf("reload test failed: %v", err)
	}
	rdr := template.NewRenderer()
	rdr.ParseFiles("testdata")
	html, err := testRender(rdr, "reload.html", nil)
	if err != nil {
		t.Errorf("reload test render failed: %v", err)
	}
	if html != TMP1 {
		t.Errorf("reload test part 1 failed, expected %s got %s", TMP1, html)
	}
	err = ioutil.WriteFile("testdata/reload.html", []byte(TMP2), 0644)
	if err != nil {
		t.Errorf("reload test failed: %v", err)
	}
	rdr.Reload()
	html, err = testRender(rdr, "reload.html", nil)
	if err != nil {
		t.Errorf("reload test render failed: %v", err)
	}
	if html != TMP2 {
		t.Errorf("reload test part 2 failed, expected %s got %s", TMP2, html)
	}
}
