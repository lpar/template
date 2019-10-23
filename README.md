
# lpar/template

This is a quick wrapper for Go's `text/template` and `html/template`, providing a higher level template manager for
web applications.

## Advantages

 - Globbing of filenames relative to a specified starting directory.
 - Templates are named by their path relative to the starting directory, and extensions are included, so template name clashes are less of a problem. For example, a template might be called `layout/header.html` rather than just `header`.
 - Files which can be minified using `tdewolff/minify` are preminified automatically. (This can be disabled.)
 - Group related template files together by naming the template set when you `.Load` them.
 - Standard familiar `html/template` style template language.
 - Nested templates, partials and other standard template features still work.
 - A `Reload()` method provides an easy way to trigger a reload when doing development.
 - Alternatively, you can set `.Live = true` on the renderer, and each time a template is executed from a template set, the corresponding files will be reloaded into that set.
 
## Example usage

	rdr := template.NewRenderer("/srv/app/templates")
	err := rdr.Load("usermgr", "users.html", "layout/*.html")
	...
	err := rdr.Execute("usermgr", os.Stdout, "users.html", data)

## Limitations

The `Reload()` method is unsafe for production use. The standard template libraries don't allow
loading of a template after any associated template has been executed. Therefore, `Reload()` has to throw away all
loaded templates and recompile them from scratch. This means that template execution in other goroutines may fail until 
the appropriate template file has been reloaded. (The code shouldn't crash, though.)

## Future features

I'm thinking I should add a way to declare that individual files shouldn't be preminimized, in case I ever
encounter a situation where that breaks a template. Maybe putting `.min.` in the filename to say the file should be 
considered minimized already?

