package controller

import "net/url"

func parsePath(p string) string {
	l, err := url.ParseRequestURI(p)
	if err != nil {
		return "/"
	}
	if len(l.Path) == 0 {
		return "/"
	}
	return l.Path
}
