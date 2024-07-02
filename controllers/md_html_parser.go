package controllers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
    "strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func MdToHTML(mdLink string) (string, error) {
    var path string
    data, err := os.ReadFile(mdLink)
    if err != nil {
        return path, errors.New(fmt.Sprintf("Error reading markdown link %s while parsing to html with error: %s", mdLink, err))
    }

    ext := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
    p := parser.NewWithExtensions(ext)
    parsedDoc := p.Parse(data)

    htmlFlags := html.CommonFlags | html.HrefTargetBlank
    opts := html.RendererOptions{
        Flags: htmlFlags,
    }

    renderer := html.NewRenderer(opts)

    root, err := os.Getwd()
    path = filepath.Join(root, "views", "htmlfiles", strings.TrimSuffix(filepath.Base(mdLink), ".md") + ".html")

    htmlFile, err := os.Create(path)
    if err != nil {
        return path, errors.New(fmt.Sprintf("Error creating html file with link %s with error: %s", path, err))
    }
    defer htmlFile.Close()

    _, err = htmlFile.Write(markdown.Render(parsedDoc, renderer))

    return path, nil
}
