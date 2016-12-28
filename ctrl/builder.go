package ctrl

import (
	"acourse/app"
	"acourse/model"
	"acourse/view"
)

// ToUserView builds a UserView from a User model
func ToUserView(m *model.User) *view.UserView {
	return &view.UserView{
		ID:       m.ID,
		Name:     m.Name,
		Username: m.Username,
		Photo:    m.Photo,
		AboutMe:  m.AboutMe,
	}
}

// ToUserMeView builds a UserMeView from a User model
func ToUserMeView(m *model.User, role *view.RoleView) *view.UserMeView {
	return &view.UserMeView{
		ID:       m.ID,
		Name:     m.Name,
		Username: m.Username,
		Photo:    m.Photo,
		AboutMe:  m.AboutMe,
		Role:     role,
	}
}

// ToUserTinyView builds a UserTinyView from a User model
func ToUserTinyView(m *model.User) *view.UserTinyView {
	return &view.UserTinyView{
		ID:       m.ID,
		Name:     m.Name,
		Username: m.Username,
		Photo:    m.Photo,
	}
}

// ToRoleView builds a RoleView fromn a Role model
func ToRoleView(m *model.Role) *view.RoleView {
	return &view.RoleView{
		Admin:      m.Admin,
		Instructor: m.Instructor,
	}
}

// ToCourseView builds a CourseView from a Course model
func ToCourseView(m *model.Course, owner *view.UserTinyView, student int, enrolled bool, owned bool) *view.CourseView {
	return &view.CourseView{
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
func ToCoursePublicView(m *model.Course, owner *view.UserTinyView, student int, purchaseStatus string) *view.CoursePublicView {
	return &view.CoursePublicView{
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
		Enroll:           m.Options.Enroll,
		Discount:         m.Options.Discount,
		PurchaseStatus:   purchaseStatus,
	}
}

// ToCourseTinyView builds a CourseTinyView from a Course model
func ToCourseTinyView(m *model.Course, owner *view.UserTinyView, student int) *view.CourseTinyView {
	return &view.CourseTinyView{
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
func ToCourseContentView(m *model.CourseContent) *view.CourseContentView {
	return &view.CourseContentView{
		Title:       m.Title,
		Description: m.Description,
		Video:       m.Video,
		DownloadURL: m.DownloadURL,
	}
}

// ToCourseContentCollectionView builds a CourseContentCollectionView from CourseContent models
func ToCourseContentCollectionView(ms []model.CourseContent) view.CourseContentCollectionView {
	r := make(view.CourseContentCollectionView, len(ms))
	for i, m := range ms {
		r[i] = ToCourseContentView(&m)
	}
	return r
}

// ToCourseContent builds a CourseContent model from CourseContent payload
func ToCourseContent(p *app.CourseContentPayload) *model.CourseContent {
	return &model.CourseContent{
		Title:       p.Title,
		Description: p.Description,
		Video:       p.Video,
		DownloadURL: p.DownloadURL,
	}
}

// ToCourseContents builds CourseContents model from CourseContents payload
func ToCourseContents(ps []*app.CourseContentPayload) []model.CourseContent {
	r := make([]model.CourseContent, len(ps))
	for i, p := range ps {
		r[i] = *ToCourseContent(p)
	}
	return r
}

// ToPaymentView builds Payment view from a Payment model
func ToPaymentView(m *model.Payment, user *view.UserTinyView, course *view.CourseMiniView) *view.PaymentView {
	return &view.PaymentView{
		ID:            m.ID,
		OriginalPrice: m.OriginalPrice,
		Price:         m.Price,
		Code:          m.Code,
		Status:        string(m.Status),
		URL:           m.URL,
		User:          user,
		Course:        course,
	}
}

// ToCourseMiniView builds Course mini view from a course model
func ToCourseMiniView(m *model.Course) *view.CourseMiniView {
	return &view.CourseMiniView{
		ID:    m.ID,
		Title: m.Title,
	}
}
