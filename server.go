package main

import (
	"log"
	"net/http"
	"sync"

	"./core"
	"github.com/ant0ine/go-json-rest/rest"
)

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/api/compiler", GetAllCompilers),
		rest.Get("/api/compiler/:name", GetCompilerDetails),
		rest.Post("/api/compiler/", Compile),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

type Compiler struct {
	Name    string
	Version string
}

type PostCode struct {
	Code        string
	Language    string
	Stdin       string
	Stdout      string
	Stderr      string
	Status_code string
}

var store = map[string]*Compiler{}
var postcode_store = map[string]*PostCode{}
var lock = sync.RWMutex{}

func GetAllCompilers(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	compilers := make([]Compiler, len(store))
	i := 0
	for _, compiler := range store {
		compilers[i] = *compiler
		i++
	}
	lock.RUnlock()
	w.WriteJson(&compilers)
}

func GetCompilerDetails(w rest.ResponseWriter, r *rest.Request) {
	// code := r.PathParam("name")
}

func Compile(w rest.ResponseWriter, r *rest.Request) {
	code := PostCode{}
	err := r.DecodeJsonPayload(&code)

	// Error Handling
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if code.Code == "" {
		rest.Error(w, "code required", 400)
		return
	}
	if code.Language == "" {
		rest.Error(w, "language required", 400)
		return
	}

	// Code Push to Container
	err = core.CodePush(code.Code, code.Language)
	if err != nil {
		rest.Error(w, "Failed push code", 400)
		return
	}
	// Compile
	result := core.Compile(code.Language, code.Stdin)
	code.Stdout = result["stdout"]
	code.Status_code = result["status_code"]

	lock.Lock()
	postcode_store[code.Language] = &code
	lock.Unlock()
	w.WriteJson(&code)

}
