package ctrl

import (
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/payload"
	"github.com/acoshift/acourse/pkg/view"
)

// ToCourseContentView builds a CourseContentView from a CourseContent model
func ToCourseContentView(m *model.CourseContent) *view.CourseContent {
	return &view.CourseContent{
		Title:       m.Title,
		Description: m.Description,
		Video:       m.Video,
		DownloadURL: m.DownloadURL,
	}
}

// ToCourseContentCollectionView builds a CourseContentCollectionView from CourseContent models
func ToCourseContentCollectionView(ms []model.CourseContent) view.CourseContentCollection {
	r := make(view.CourseContentCollection, len(ms))
	for i, m := range ms {
		r[i] = ToCourseContentView(&m)
	}
	return r
}

// ToCourseContent builds a CourseContent model from CourseContent payload
func ToCourseContent(p *payload.CourseContent) *model.CourseContent {
	return &model.CourseContent{
		Title:       p.Title,
		Description: p.Description,
		Video:       p.Video,
		DownloadURL: p.DownloadURL,
	}
}

// ToCourseContents builds CourseContents model from CourseContents payload
func ToCourseContents(ps []*payload.CourseContent) []model.CourseContent {
	r := make([]model.CourseContent, len(ps))
	for i, p := range ps {
		r[i] = *ToCourseContent(p)
	}
	return r
}

// ToCourseMiniView builds Course mini view from a course model
func ToCourseMiniView(m *model.Course) *view.CourseMini {
	return &view.CourseMini{
		ID:    m.ID,
		Title: m.Title,
	}
}
