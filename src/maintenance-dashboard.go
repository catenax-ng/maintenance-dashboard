package main

import (
    "fmt"
    "os"
)

type Page struct {
    Title string
    Body  []byte
}

func (p *Page) save() error {
    filename := p.Title + ".txt"
    return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func main() {
    p1 := &Page{Title: "DevSecOps", Body: []byte("This is the page we want to have information on the currently deployed and the latest releases of all our applications.")}
    p1.save()
    p2, _ := loadPage("DevSecOps")
    fmt.Println(string(p2.Body))
}
