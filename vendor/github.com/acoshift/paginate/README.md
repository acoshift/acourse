# paginate

[![Build Status](https://travis-ci.org/acoshift/paginate.svg?branch=master)](https://travis-ci.org/acoshift/paginate)
[![Coverage Status](https://coveralls.io/repos/github/acoshift/paginate/badge.svg?branch=master)](https://coveralls.io/github/acoshift/paginate?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/acoshift/paginate)](https://goreportcard.com/report/github.com/acoshift/paginate)
[![GoDoc](https://godoc.org/github.com/acoshift/paginate?status.svg)](https://godoc.org/github.com/acoshift/paginate)

Pagination Logic for your template

## Example

```html
{{define "pagination"}}
<div>
    <a href="?page=1">First</a>
    <a href="?page={{.Prev}}">Prev</a>
    {{range .Pages 2 2}}
        {{if eq . 0}}
            <a class="disabled">...</a>
        {{else if eq $.Page .}}
            <a class="active">{{.}}</a>
        {{else}}
            <a href="?page={{.}}">{{.}}</a>
        {{end}}
    {{end}}
    <a href="?page={{.Next}}">Next</a>
    <a href="?page={{.MaxPage}}">Last</a>
</div>
{{end}}
```

![Example](https://github.com/acoshift/paginate/raw/master/demo.gif)
