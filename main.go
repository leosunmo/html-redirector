package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

const redirURL = `{{.RedirectURL}}/{{if .ErrorCode}}?error={{.ErrorCode}}{{end}}{{if .ReturnURL}}&return={{.ReturnURL}}{{end}}`

const html = `
<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="Refresh" content="0; url={{.URL}}" />
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

	http.HandleFunc("/", handler)

	log.Printf("Serving on HTTP port: %s\n", "3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	templates := template.New("error_templates")
	redirURL, err := templates.New("url").Parse(redirURL)
	if err != nil {
		fmt.Printf("url tmpl err: %s", err.Error())
	}
	htmlResponse, err := templates.New("html").Parse(html)
	if err != nil {
		fmt.Printf("html tmpl err: %s", err.Error())
	}
	reqValues := r.URL.Query()
	URLconf := redirectConf{
		RedirectURL: os.Getenv("REDIRECT_URL"),
		ReturnURL:   reqValues.Get("return"),
		ErrorCode:   reqValues.Get("error"),
	}
	renderedURLTmpl := bytes.Buffer{}
	redirURL.Execute(&renderedURLTmpl, URLconf)
	fmt.Printf("Sending %s (%s) to %s\n", r.RemoteAddr, r.UserAgent(), renderedURLTmpl.String())
	htmlConf := struct {
		URL string
	}{
		URL: renderedURLTmpl.String(),
	}
	htmlResponse.Execute(w, htmlConf)
}
