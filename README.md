# rewrite_absolute_links
Caddyserver middleware to rewrite absolute links to relative ones on the fly

Directive usage:

    rewrite_absolute_links foo.example.com
    rewrite_absolute_links bar.example.com
    rewrite_absolute_links baz.example.com


OR

    rewrite_absolute_links foo.example.com bar.example.com baz.example.com

OR

    rewrite_absolute_links foo.example.com bar.example.com
    rewrite_absolute_links baz.example.com

Caddy will parse any responses with the Content-Type "text/html" as HTML, and rewrite any anchor tags with hrefs whose URL contains one
of the configured domains so that it is a relative link.

Expected behavior:

    <a href="http://foo.example.com/hello.html"> **becomes**

    <a href="/hello.html">
    
The HTML5 parser will rewrite non-conforming HTML. It can't be helped, mate.
