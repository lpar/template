
# lpar/template

This is a quick wrapper for Go's `text/template` and `html/template`, providing a higher level template manager for
web applications.

## Advantages

 - It keeps track of whether each file is HTML or not, and uses `html/template` where appropriate for safety.
 - Pass it a directory and it recurses through and loads and parses all the files. Each template is named with the 
 filename relative to the specified directory.
 - Files which can be minified using `tdewolff/minify` are preminified automatically.
 - Standard familiar `html/template` style API.
 - Nested templates, partials and other template features still work.
 - A `Reload()` method provides an easy way to trigger a reload when doing development.

## Limitations

The `Reload()` method is unsafe for production use. The standard template libraries don't allow
loading of a template after any associated template has been executed. Therefore, `Reload()` has to throw away all
loaded templates and reload them from scratch. This means that template execution in other goroutines may fail until 
the appropriate template file has been reloaded. (The code shouldn't crash, though.)

## Future features

I'm thinking I should add a way to declare that files shouldn't be preminimized, in case I ever
encounter a situation where that breaks a template. Maybe putting `.min.` in the filename to say the file should be 
considered minimized already?


