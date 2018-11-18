package view

import (
	"context"

	"github.com/acoshift/acourse/internal/context/appctx"
)

// PageData type
type PageData struct {
	Meta PageMeta
	Data map[string]interface{}
}

// PageMeta type
type PageMeta struct {
	Title string
	Desc  string
	Image string
	URL   string
}

// Page creates new page data
func Page(ctx context.Context) *PageData {
	return &PageData{
		Meta: PageMeta{
			Title: "",
			Desc:  "แหล่งเรียนรู้ออนไลน์ ที่ทุกคนเข้าถึงได้",
			Image: "https://storage.googleapis.com/acourse/static/62b9eb0e-3668-4f9f-86b7-a11349938f7a.jpg",
			URL:   "https://acourse.io",
		},
		Data: map[string]interface{}{
			"Me":     appctx.GetUser(ctx),
			"Flash":  appctx.GetFlash(ctx),
			"Navbar": "",
		},
	}
}
