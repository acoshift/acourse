package ctrl

import (
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/payload"
	"github.com/acoshift/acourse/pkg/view"
)

// ToUserView builds an UserView from a User model
func ToUserView(m *model.User) *view.User {
	return &view.User{
		ID:       m.ID,
		Name:     m.Name,
		Username: m.Username,
		Photo:    m.Photo,
		AboutMe:  m.AboutMe,
	}
}

// ToUserCollectionView builds an UserCollection from User models
func ToUserCollectionView(ms []*model.User) view.UserCollection {
	rs := make(view.UserCollection, len(ms))
	for i := range ms {
		rs[i] = ToUserView(ms[i])
	}
	return rs
}

// ToUserMeView builds an UserMeView from a User model
func ToUserMeView(m *model.User, role *view.Role) *view.UserMe {
	return &view.UserMe{
		ID:       m.ID,
		Name:     m.Name,
		Username: m.Username,
		Photo:    m.Photo,
		AboutMe:  m.AboutMe,
		Role:     role,
	}
}

// ToUserTinyView builds an UserTinyView from a User model
func ToUserTinyView(m *model.User) *view.UserTiny {
	return &view.UserTiny{
		ID:       m.ID,
		Name:     m.Name,
		Username: m.Username,
		Photo:    m.Photo,
	}
}

// ToRoleView builds a RoleView fromn a Role model
func ToRoleView(m *model.Role) *view.Role {
	return &view.Role{
		Admin:      m.Admin,
		Instructor: m.Instructor,
	}
}

// ToCourseView builds a CourseView from a Course model
func ToCourseView(m *model.Course, owner *view.UserTiny, student int, enrolled bool, owned bool) *view.Course {
	return &view.Course{
		ID:               m.ID,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
		Owner:            owner,
		Title:            m.Title,
		ShortDescription: m.ShortDescription,
		Description:      m.Description,
		Photo:            m.Photo,
		Start:            m.Start,
		URL:              m.URL,
		Video:            m.Video,
		Type:             string(m.Type),
		Price:            m.Price,
		DiscountedPrice:  m.DiscountedPrice,
		Student:          student,
		Contents:         ToCourseContentCollectionView(m.Contents),
		EnrollDetail:     m.EnrollDetail,
		Enrolled:         enrolled,
		Owned:            owned,
		Enroll:           m.Options.Enroll,
		Public:           m.Options.Public,
		Attend:           m.Options.Attend,
		Assignment:       m.Options.Assignment,
		Discount:         m.Options.Discount,
	}
}

// ToCoursePublicView builds a CourseView from a Course model
func ToCoursePublicView(m *model.Course, owner *view.UserTiny, student int, purchaseStatus string) *view.CoursePublic {
	return &view.CoursePublic{
		ID:               m.ID,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
		Owner:            owner,
		Title:            m.Title,
		ShortDescription: m.ShortDescription,
		Description:      m.Description,
		Photo:            m.Photo,
		Start:            m.Start,
		URL:              m.URL,
		Video:            m.Video,
		Type:             string(m.Type),
		Price:            m.Price,
		DiscountedPrice:  m.DiscountedPrice,
		EnrollDetail:     m.EnrollDetail,
		Student:          student,
		Enroll:           m.Options.Enroll,
		Discount:         m.Options.Discount,
		PurchaseStatus:   purchaseStatus,
	}
}

// ToCourseTinyView builds a CourseTinyView from a Course model
func ToCourseTinyView(m *model.Course, owner *view.UserTiny, student int) *view.CourseTiny {
	return &view.CourseTiny{
		ID:               m.ID,
		Owner:            owner,
		Title:            m.Title,
		ShortDescription: m.ShortDescription,
		Photo:            m.Photo,
		Start:            m.Start,
		URL:              m.URL,
		Type:             string(m.Type),
		Price:            m.Price,
		DiscountedPrice:  m.DiscountedPrice,
		Student:          student,
		Discount:         m.Options.Discount,
	}
}

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

// ToPaymentView builds Payment view from a Payment model
func ToPaymentView(m *model.Payment, user *view.UserTiny, course *view.CourseMini) *view.Payment {
	return &view.Payment{
		ID:            m.ID,
		OriginalPrice: m.OriginalPrice,
		Price:         m.Price,
		Code:          m.Code,
		Status:        string(m.Status),
		URL:           m.URL,
		User:          user,
		Course:        course,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		At:            m.At,
	}
}

// ToCourseMiniView builds Course mini view from a course model
func ToCourseMiniView(m *model.Course) *view.CourseMini {
	return &view.CourseMini{
		ID:    m.ID,
		Title: m.Title,
	}
}
