package template

import (
	"errors"
	"fmt"
	htmlplate "html/template"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/tdewolff/minify/js"

	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/xml"

	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"

	"github.com/tdewolff/minify"
)

// Renderer is an object which loads, minifies and renders templates for CSS, JS, HTML, SVG, and XML.
type Renderer struct {
	minifier   *minify.M
	paths      []string
	templates  *template.Template
	htmlplates *htmlplate.Template
}

// NewRenderer returns a Renderer ready for use.
func NewRenderer() *Renderer {
	min := minify.New()
	min.AddFunc("text/html", html.Minify)
	min.AddFunc("text/css", css.Minify)
	min.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	min.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	min.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	return &Renderer{minifier: min}
}

func existsOrError(fspc string) error {
	if _, err := os.Stat(fspc); err != nil {
		if os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func (r *Renderer) processFile(basepath string, path string, finfo os.FileInfo, err error) error {
	if finfo.IsDir() {
		return nil
	}
	rel, err := filepath.Rel(basepath, path)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	ext := filepath.Ext(path)
	mtype := mime.TypeByExtension(ext)
	_, _, mfn := r.minifier.Match(mtype)
	var dat []byte
	if mfn != nil {
		dat, err = r.minifier.Bytes(mtype, data)
	} else {
		dat = data
		err = nil
	}
	if err != nil {
		return err
	}

	if strings.HasPrefix(mtype, "text/html") {
		return r.storeHTML(rel, dat)
	}
	return r.storeOther(rel, dat)
}

func (r *Renderer) storeHTML(name string, data []byte) error {
	var tmpl *htmlplate.Template
	if r.htmlplates == nil {
		tmpl = htmlplate.New(name)
		r.htmlplates = tmpl
	} else {
		tmpl = r.htmlplates.New(name)
	}
	_, err := tmpl.Parse(string(data))
	return err
}

func (r *Renderer) storeOther(name string, data []byte) error {
	var tmpl *template.Template
	if r.templates == nil {
		tmpl = template.New(name)
		r.templates = tmpl
	} else {
		tmpl = r.templates.New(name)
	}
	_, err := tmpl.Parse(string(data))
	return err
}

// ParseFiles loads, minifies and parses all template files under the supplied base path.
// Each template is named according to its relative path under the base path.
// HTML files are parsed into html/template templates, other files are parsed into
// text/template templates.
func (r *Renderer) ParseFiles(basepath string) error {
	err := existsOrError(basepath)
	if err != nil {
		return fmt.Errorf("asked to scan template directory %s which does not exist", basepath)
	}
	r.paths = append(r.paths, basepath)
	return filepath.Walk(basepath, func(path string, finfo os.FileInfo, err error) error {
		return r.processFile(basepath, path, finfo, err)
	})
}

// Reload causes the Renderer to rescan all of the directories it has been told to scan before,
// and reload all of the template files.
//
// This method is intended for use during development. It should not be used in production, as
// the cached template update process is unsafe and template execution in other goroutines
// may fail while the files are being reloaded.
func (r *Renderer) Reload() error {
	r.templates = nil
	r.htmlplates = nil
	for _, basepath := range r.paths {
		err := r.ParseFiles(basepath)
		if err != nil {
			return err
		}
	}
	return nil
}

// ExecuteTemplate executes the template with the specified name, passing it the provided data, and sending
// the output to the provided io.Writer.
func (r *Renderer) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	ext := filepath.Ext(name)
	mtype := mime.TypeByExtension(ext)
	if strings.HasPrefix(mtype, "text/html") {
		if r.htmlplates == nil {
			return errors.New("no HTML templates found")
		}
		return r.htmlplates.ExecuteTemplate(wr, name, data)
	}
	if r.templates == nil {
		return errors.New("no text templates found")
	}
	return r.templates.ExecuteTemplate(wr, name, data)
}
