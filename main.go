package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

const htmlResp = `
<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="Refresh" content="0; url=%s" />
  </head>
</html>
`

const redirTmpl = "{{.RedirectURL}}/{{if .ErrorCode}}?error={{.ErrorCode}}{{end}}{{if .ReturnURL}}&return={{.ReturnURL}}{{end}}"

type redirectConf struct {
	RedirectURL string
	ReturnURL   string
	Options     string
	ErrorCode   string
}

var (
	redirURL         = template.Must(template.New("redirURL").Parse(redirTmpl))
	redirectURL      string
	defaultReturnURL string
)

func main() {
	var found bool
	redirectURL, found = os.LookupEnv("REDIRECT_URL")
	if !found {
		log.Fatal("REDIRECT_URL not set")
	}
	defaultReturnURL, found = os.LookupEnv("DEFAULT_RETURN_URL")
	if !found {
		log.Fatal("DEFAULT_RETURN_URL not set")
	}

	http.HandleFunc("/", redirectHandler())

	log.Printf("Serving on HTTP port: %s\n", "3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func redirectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		returnScheme := r.Header.Get("X-Forwarded-Proto")
		returnURL := r.Header.Get("X-Forwarded-Host")
		if returnURL == "" {
			returnURL = defaultReturnURL
		}

		queryValues := r.URL.Query()
		URLconf := redirectConf{
			RedirectURL: redirectURL,
			ReturnURL:   returnScheme + "://" + returnURL,
			ErrorCode:   queryValues.Get("error"),
		}

		fmt.Printf("Sending %s (%s) to %s\n", r.RemoteAddr, r.UserAgent(), returnURL)
		renderedURLTmpl := bytes.Buffer{}
		redirURL.Execute(&renderedURLTmpl, URLconf)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, htmlResp, renderedURLTmpl.String())
	}
}
