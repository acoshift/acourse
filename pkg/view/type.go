package view

// Page type provides layout data like title, description, and og
type Page struct {
	Title string
	Desc  string
	Image string
	URL   string
}

// IndexData type
type IndexData struct {
	*Page
}
