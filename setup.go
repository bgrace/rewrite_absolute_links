package rewrite_absolute_links

import (
	"github.com/mholt/caddy"
	"log"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"net/http"
	"strings"
	"os"
	"bytes"
)

var info = log.New(os.Stdout, "", 0)
var warn = log.New(os.Stderr, "WARNING: ", 0)

func readDomains(c *caddy.Controller) []string {
	domains := []string{}

	for c.NextArg() {
		domains = append(domains, c.Val())
	}

	return domains
}

func setup(c *caddy.Controller) (err error) {
	cfg := httpserver.GetConfig(c)

	domains := []string{}

	c.Next() // consume the directive name
	domains = append(domains, readDomains(c)...)

	for c.NextLine() { // consumes directive name
		domains = append(domains, readDomains(c)...)
	}

	if len(domains) == 0 {
		warn.Printf("The rewrite_absolute_links directive was specified for %v but no domains were found", cfg.ListenHost)
		return c.ArgErr()
	} else {
		info.Printf("Rewriting absolute hrefs to relative ones for the following domains: %v", domains)
	}

	mid := func(next httpserver.Handler) httpserver.Handler {
		return myHandler{Next: next, domains: domains}
	}

	cfg.AddMiddleware(mid)

	return nil
}

type myHandler struct {
	Next httpserver.Handler
	domains []string
}

type htmlInterceptResponseWriter struct {
	buff *bytes.Buffer
	http.ResponseWriter
}

func newHTMLInterceptResponseWriter(w http.ResponseWriter) *htmlInterceptResponseWriter {
	return &htmlInterceptResponseWriter{bytes.NewBufferString(""), w}
}

func (w htmlInterceptResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w htmlInterceptResponseWriter) Write(bytes []byte) (int, error) {
	if strings.HasPrefix(w.Header().Get("Content-Type"), "text/html") {
		return w.buff.Write(bytes)
	}

	return w.ResponseWriter.Write(bytes)
}

func (h myHandler) ServeHTTP(handlerWriter http.ResponseWriter, r *http.Request) (int, error) {

	w := newHTMLInterceptResponseWriter(handlerWriter)
	status, err := h.Next.ServeHTTP(w, r)

	if contentType := w.Header().Get("Content-Type"); strings.HasPrefix(contentType, "text/html") {
		transformed := transform(w.buff.Bytes(), h.domains)
		handlerWriter.Write(transformed)
	}

	return status, err

}

func init() {
	caddy.RegisterPlugin("rewrite_absolute_links", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}


