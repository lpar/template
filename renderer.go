package template

import (
	"fmt"
	htmlplate "html/template"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"regexp"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/xml"
)

// TemplateSet represents a set of 1 or more Go templates.
type TemplateSet struct {
	name      string
	globs     []string
	renderer  *Renderer
	templates *htmlplate.Template
}

// Renderer is an object for loading, parsing, minifying and rendering HTML templates.
type Renderer struct {
	basePath     string
	templateSets map[string]*TemplateSet
	// Minify is true if the template content should be minified. It's possible that the minification code could have
	// problems with complex Go templates, but I've not encountered any problems of that nature so far.
	Minify       bool
	// Live is true if the template sets should be reloaded whenever Execute is used.
	// This is obviously a very bad idea in production and should only be used for development.
	Live         bool
	minifier     *minify.M
}

// NewRenderer returns an initialized Renderer. The basepath is used to locate all files subsequently added via Load().
func NewRenderer(basepath string) Renderer {
	ren := Renderer{basePath: basepath, templateSets: make(map[string]*TemplateSet)}
	min := minify.New()
	min.AddFunc("text/html", html.Minify)
	min.AddFunc("text/css", css.Minify)
	min.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	min.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	min.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	ren.minifier = min
	return ren
}

// Load loads one or more template files into the specified template set.
func (tm *Renderer) Load(templateSet string, fileglobs ... string) error {
	thing := &TemplateSet{
		name:     templateSet,
		globs:    fileglobs,
		renderer: tm,
	}
	tm.templateSets[templateSet] = thing
	return thing.Load()
}

// execute executes the named template from the named template set.
func (tm Renderer) Execute(templateSet string, wr io.Writer, tmplname string, data interface{}) error {
	var err error
	tmpl, ok := tm.templateSets[templateSet]
	if ok {
		if tm.Live {
			err = tmpl.Load()
			if err != nil {
				return err
			}
		}
		err = tmpl.execute(wr, tmplname, data)
	} else {
		err = fmt.Errorf("no such template file set %s", templateSet)
	}
	return err
}

// Reload reloads the named template set.
func (tm Renderer) Reload(templateSet string) error {
	var err error
	thing, ok := tm.templateSets[templateSet]
	if ok {
		err = thing.Load()
	} else {
		err = fmt.Errorf("reload failed, file set %s unknown", templateSet)
	}
	return err
}

// minify handles minification of loaded file data
func (tm Renderer) minify(filename string, dat []byte) (string, error) {
	if !tm.Minify {
		return string(dat), nil
	}
	var err error
	ext := filepath.Ext(filename)
	mtype := mime.TypeByExtension(ext)
	_, _, mfn := tm.minifier.Match(mtype)
	if mfn != nil {
		dat, err = tm.minifier.Bytes(mtype, dat)
	}
	return string(dat), err
}

// Load causes a template set to be loaded from the files in the filesystem.
// Minification is performed if configured, and the template files are compiled.
func (th *TemplateSet) Load() error {
	var tmpl *htmlplate.Template
	for _, glob := range th.globs {
		globpath := filepath.Join(th.renderer.basePath, glob)
		matches, err := filepath.Glob(globpath)
		if err != nil {
			return err
		}
		for _, filename := range matches {
			fi, err := os.Stat(filename)
			if err != nil {
				return err
			}
			if fi.IsDir() {
				continue
			}
			name, err := filepath.Rel(th.renderer.basePath, filename)
			if err != nil {
				return err
			}
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}
			dat, err := th.renderer.minify(filename, data)
			if err != nil {
				return err
			}
			if tmpl == nil {
				tmpl, err = htmlplate.New(name).Parse(dat)
			} else {
				tmpl, err = tmpl.New(name).Parse(dat)
			}
			if err != nil {
				return err
			}
		}
	}
	th.templates = tmpl
	return nil
}

// execute executes an individual template set.
func (th TemplateSet) execute(wr io.Writer, tmplname string, data interface{}) error {
	return th.templates.ExecuteTemplate(wr, tmplname, data)
}

