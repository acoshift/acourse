package app

import "net/url"

func makeLink(path string, query url.Values) string {
	link, _ := url.Parse(baseURL)
	link.Path = path
	if query != nil {
		link.RawQuery = query.Encode()
	}
	return link.String()
}
