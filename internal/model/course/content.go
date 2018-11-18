package course

// Content type
type Content struct {
	ID          string
	CourseID    string
	Title       string
	Desc        string
	VideoID     string
	VideoType   int
	DownloadURL string
}

// CreateContent creates new course content
type CreateContent struct {
	ID        string
	Title     string
	LongDesc  string
	VideoID   string
	VideoType int

	Result string // content id
}

// UpdateContent updates a course content
type UpdateContent struct {
	ContentID string
	Title     string
	Desc      string
	VideoID   string
}

// GetContent gets a course's content
type GetContent struct {
	ContentID string

	Result *Content
}

// ListContents lists course's contents
type ListContents struct {
	ID string

	Result []*Content
}

// DeleteContent deletes a course's content
type DeleteContent struct {
	ContentID string
}
