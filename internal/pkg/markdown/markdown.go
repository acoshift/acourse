package markdown

import (
	"html/template"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// Email converts markdown to html for "email"
func Email(s string) string {
	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: 0 |
			blackfriday.UseXHTML |
			blackfriday.Smartypants |
			blackfriday.SmartypantsFractions |
			blackfriday.SmartypantsDashes |
			blackfriday.SmartypantsLatexDashes,
	})

	extension := 0 |
		blackfriday.NoIntraEmphasis |
		blackfriday.FencedCode |
		blackfriday.Autolink |
		blackfriday.Strikethrough |
		blackfriday.SpaceHeadings |
		blackfriday.HeadingIDs |
		blackfriday.BackslashLineBreak |
		blackfriday.DefinitionLists

	md := blackfriday.Run([]byte(s), blackfriday.WithRenderer(renderer), blackfriday.WithExtensions(extension))
	p := bluemonday.UGCPolicy()
	return string(p.SanitizeBytes(md))
}

// HTML converts markdown to html
func HTML(s string) template.HTML {
	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: 0 |
			blackfriday.UseXHTML |
			blackfriday.Smartypants |
			blackfriday.SmartypantsFractions |
			blackfriday.SmartypantsDashes |
			blackfriday.SmartypantsLatexDashes |
			blackfriday.HrefTargetBlank,
	})

	extension := 0 |
		blackfriday.NoIntraEmphasis |
		blackfriday.Tables |
		blackfriday.FencedCode |
		blackfriday.Autolink |
		blackfriday.Strikethrough |
		blackfriday.SpaceHeadings |
		blackfriday.HeadingIDs |
		blackfriday.BackslashLineBreak |
		blackfriday.DefinitionLists

	md := blackfriday.Run([]byte(s), blackfriday.WithRenderer(renderer), blackfriday.WithExtensions(extension))
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("target").OnElements("a")
	return template.HTML(p.SanitizeBytes(md))
}
