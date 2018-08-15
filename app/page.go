package app

import (
	"context"

	"github.com/acoshift/acourse/context/appctx"
)

func page(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"Title":  "",
		"Desc":   "แหล่งเรียนรู้ออนไลน์ ที่ทุกคนเข้าถึงได้",
		"Image":  "https://storage.googleapis.com/acourse/static/62b9eb0e-3668-4f9f-86b7-a11349938f7a.jpg",
		"URL":    "https://acourse.io",
		"Me":     appctx.GetUser(ctx),
		"Flash":  appctx.GetSession(ctx).Flash(),
		"Navbar": "",
	}
}
