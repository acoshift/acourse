package controller

import "net/url"

func (c *ctrl) makeLink(path string, query url.Values) string {
	link, _ := url.Parse(c.baseURL)
	link.Path = path
	if query != nil {
		link.RawQuery = query.Encode()
	}
	return link.String()
}
