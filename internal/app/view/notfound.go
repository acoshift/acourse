package view

import (
	"math/rand"
	"net/http"

	"github.com/acoshift/header"
	"github.com/moonrhythm/hime"
)

var notFoundImages = []string{
	"https://storage.googleapis.com/acourse/static/9961f3c1-575f-4b98-af4f-447566ee1cb3.png",
	"https://storage.googleapis.com/acourse/static/b14a40c9-d3a4-465d-9453-ce7fcfbc594c.png",
}

// NotFound renders not found page
func NotFound(ctx *hime.Context) error {
	p := Page(ctx)
	p.Data["Image"] = notFoundImages[rand.Intn(len(notFoundImages))]
	ctx.SetHeader(header.XContentTypeOptions, "nosniff")
	return ctx.Status(http.StatusNotFound).View("share.not-found", p)
}
