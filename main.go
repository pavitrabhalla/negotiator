package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}

type indexOject struct {
	Name string
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	obj := indexOject{
		Name: req.FormValue("name"),
	}

	// Get template.
	f, err := ioutil.ReadFile("web/index.html")
	if err != nil {
		res.WriteHeader(400)
		return
	}

	tmpl := template.Must(template.New("index").Parse(string(f)))
	if err := tmpl.Execute(res, obj); err != nil {
		res.WriteHeader(400)
		return
	}
}
