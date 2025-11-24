package main

//#cgo CFLAGS: -D_GNU_SOURCE
//#include "php_bridge.h"
import "C"
import (
	"fmt"
	"log"
	"net/http"
	"os"
	pathlib "path"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func checkPaths(path string, checker func(string) bool) string {
	paths := []string{
		path,
		pathlib.Join(path, "index"),
		pathlib.Join(path, "index.md"),
		pathlib.Join(path, "index.php"),
		pathlib.Join(path, "index.html"),
		strings.TrimSuffix(path, "/") + ".md",
		strings.TrimSuffix(path, "/") + ".php",
		strings.TrimSuffix(path, "/") + ".html",
	}
	for _, path := range paths {
		if checker(path) {
			return path
		}
	}
	return ""
}

func renderPHP(path string) string {
	C.render_php(C.CString(path))
	defer C.php_reset_output()
	c := C.php_get_output()
	if c == nil {
		return ""
	}
	return C.GoString(c)
}

func renderMD(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Print(err)
		return ""
	}
	p := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock)
	doc := p.Parse(data)
	renderer := html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank | html.CompletePage})

	return string(markdown.Render(doc, renderer))
}

var root string

func handler(w http.ResponseWriter, r *http.Request) {
	requestPath := pathlib.Join(root, r.URL.Path)
	filePath := checkPaths(requestPath, func(s string) bool {
		stat, err := os.Stat(s)
		return !(os.IsNotExist(err) || stat.IsDir())
	})
	parts := strings.Split(strings.TrimSuffix(filePath, "/"), ".")
	fileExt := parts[len(parts)-1]
	switch fileExt {
	case "html":
		fallthrough
	case "php":
		fmt.Fprint(w, renderPHP(filePath))
	case "md":
		fmt.Fprint(w, renderMD(filePath))
	default:
		fmt.Fprint(w, requestPath + ":" + filePath)
	}
}

func main() {
	if len(os.Args) < 2 {
		var err error
		root, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	} else {
		root = os.Args[1]
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
