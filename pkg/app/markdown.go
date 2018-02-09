package app

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// markdown converts markdown to html for "email"
func markdown(s string) string {
	renderer := blackfriday.HtmlRenderer(
		0|
			blackfriday.HTML_USE_XHTML|
			blackfriday.HTML_USE_SMARTYPANTS|
			blackfriday.HTML_SMARTYPANTS_FRACTIONS|
			blackfriday.HTML_SMARTYPANTS_DASHES|
			blackfriday.HTML_SMARTYPANTS_LATEX_DASHES,
		"", "")
	md := blackfriday.MarkdownOptions([]byte(s), renderer, blackfriday.Options{
		Extensions: 0 |
			blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
			blackfriday.EXTENSION_FENCED_CODE |
			blackfriday.EXTENSION_AUTOLINK |
			blackfriday.EXTENSION_STRIKETHROUGH |
			blackfriday.EXTENSION_SPACE_HEADERS |
			blackfriday.EXTENSION_HEADER_IDS |
			blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
			blackfriday.EXTENSION_DEFINITION_LISTS,
	})
	p := bluemonday.UGCPolicy()
	return string(p.SanitizeBytes(md))
}
