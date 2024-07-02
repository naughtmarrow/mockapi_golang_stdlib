package controllers

import(
    "testing"
    "fmt"
)

func TestParser(t *testing.T) {
    path := "../views/mdfiles/testfile.md"
    newPath, err := MdToHTML(path)
    if err != nil {
        t.Fatal(fmt.Sprintf("Something went wrong in parser with newPath: %s\n and error: %s", newPath, err))
    }
}
