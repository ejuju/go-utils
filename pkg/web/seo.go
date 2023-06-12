package web

import "net/http"

// Generate sitemap XML
func SitemapXML(host string, routes ...string) string {
	out := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
	out += "<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n"
	for _, r := range routes {
		out += "\t<url>\n"
		out += "\t\t<loc>" + "https://" + host + r + "</loc>\n"
		out += "\t\t<priority>1</priority>\n"
		out += "\t</url>\n"
	}
	out += "</urlset>\n"
	return out
}

func ServeSitemapXML(host string, routes ...string) http.HandlerFunc {
	return ServeRaw([]byte(SitemapXML(host, routes...)))
}

// Generate robots.txt with disallowed routes
func DisallowedRobotsTXT(routes []string) string {
	out := "User-agent: *\n"
	for _, r := range routes {
		out += "Disallow: " + r + "\n"
	}
	return out
}

// Website JSON+LD schema
type JSONLDWebsiteSchema struct {
	Context string `json:"@context"`
	Type    struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"@type"`
}
