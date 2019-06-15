package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

const html = `
<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="Refresh" content="0; url={{.RedirectURL}}/?error={{.ErrorCode}}&return={{.ReturnURL}}" />
  </head>
</html>
`

type redirectConf struct {
	RedirectURL string
	ReturnURL   string
	Options     string
	ErrorCode   string
}

func main() {
	_, found := os.LookupEnv("REDIRECT_URL")
	if !found {
		log.Fatal("REDIRECT_URL not set")
	}
	_, found = os.LookupEnv("RETURN_URL")
	if !found {
		log.Fatal("RETURN_URL not set")
	}

	http.HandleFunc("/", handler)

	log.Printf("Serving on HTTP port: %s\n", "3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	templates := template.New("error_templates")
	templates.New("html").Parse(html)
	conf := redirectConf{
		RedirectURL: os.Getenv("REDIRECT_URL"),
		ReturnURL:   os.Getenv("RETURN_URL"),
		ErrorCode:   strings.TrimPrefix(r.URL.EscapedPath(), "/"),
	}
	fmt.Printf("Sending %s/?error=%s&return=%s\n", conf.RedirectURL, conf.ErrorCode, conf.ReturnURL)
	templates.Lookup("html").Execute(w, conf)
}
