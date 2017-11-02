package rewrite_absolute_links

import (
	"golang.org/x/net/html"
	"bytes"
	"net/url"
	"net"
	"strings"
)

func transform(in []byte, domains []string) []byte {

	tokenizer := html.NewTokenizer(bytes.NewReader(in))

	var buf []byte
	out := bytes.NewBuffer(buf)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return out.Bytes()
		}

		mutatedToken := false

		if tokenType == html.StartTagToken {
			token := tokenizer.Token()

			if token.Data == "a" {
				for index, a := range token.Attr {
					if a.Key == "href" {
						u, err := url.Parse(a.Val)
						if err != nil {
							continue
						}

						var host string
						if strings.Contains(u.Host, ":") {
							h, _, err := net.SplitHostPort(u.Host)
							if err != nil {
								continue
							}
							host = h

						} else {
							host = u.Host
						}

						// see if the host is in the list of domains
						for _, domain := range domains {
							if host == domain {
								u.Host = ""
								u.Scheme = ""
								mutatedToken = true
								break
							}
						}

						if mutatedToken {
							token.Attr[index].Val = u.String()
							out.WriteString(token.String())
						}
					}
				}

			}
		}

		if !mutatedToken {
			out.Write(tokenizer.Raw())
		}

	}

	return out.Bytes()
}
