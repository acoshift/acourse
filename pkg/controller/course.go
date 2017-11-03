package controller

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/acoshift/header"
	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/acoshift/acourse/pkg/app"
)

func (c *ctrl) CourseView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := app.GetUser(ctx)
	link := app.GetCourseURL(ctx)

	// if id can parse to uuid get course from id
	id := link
	_, err := uuid.Parse(link)
	if err != nil {
		// link can not parse to uuid get course id from url
		id, err = c.repo.GetCourseIDFromURL(ctx, link)
		if err == app.ErrNotFound {
			c.view.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	x, err := c.repo.GetCourse(ctx, id)
	if err == app.ErrNotFound {
		c.view.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if course has url, redirect to course url
	if x.URL.Valid && x.URL.String != link {
		http.Redirect(w, r, "/course/"+x.URL.String, http.StatusFound)
		return
	}

	enrolled := false
	pendingEnroll := false
	if user != nil {
		enrolled, err = c.repo.IsEnrolled(ctx, user.ID, x.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !enrolled {
			pendingEnroll, err = c.repo.HasPendingPayment(ctx, user.ID, x.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	var owned bool
	if user != nil {
		owned = user.ID == x.UserID
	}

	// if user enrolled or user is owner fetch course contents
	if enrolled || owned {
		x.Contents, err = c.repo.GetCourseContents(ctx, x.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if owned {
		x.Owner = user
	} else {
		x.Owner, err = c.repo.GetUser(ctx, x.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	c.view.Course(w, r, x, enrolled, owned, pendingEnroll)
}

func (c *ctrl) CourseContent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := app.GetUser(ctx)
	link := app.GetCourseURL(ctx)

	// if id can parse to uuid get course from id
	id := link
	_, err := uuid.Parse(link)
	if err != nil {
		// link can not parse to uuid get course id from url
		id, err = c.repo.GetCourseIDFromURL(ctx, link)
		if err == app.ErrNotFound {
			c.view.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	x, err := c.repo.GetCourse(ctx, id)
	if err == app.ErrNotFound {
		c.view.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if course has url, redirect to course url
	if x.URL.Valid && x.URL.String != link {
		http.Redirect(w, r, "/course/"+x.URL.String+"/content", http.StatusFound)
		return
	}

	enrolled, err := c.repo.IsEnrolled(ctx, user.ID, x.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !enrolled && user.ID != x.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	x.Contents, err = c.repo.GetCourseContents(ctx, x.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	x.Owner, err = c.repo.GetUser(ctx, x.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var content *app.CourseContent
	p, _ := strconv.Atoi(r.FormValue("p"))
	if p < 0 {
		p = 0
	}
	if p > len(x.Contents)-1 {
		p = len(x.Contents) - 1
	}
	if p >= 0 {
		content = x.Contents[p]
	}

	c.view.CourseContent(w, r, x, content)
}

func (c *ctrl) EditorCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.postEditorCreate(w, r)
		return
	}
	c.view.EditorCreate(w, r)
}

func (c *ctrl) postEditorCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f := app.GetSession(ctx).Flash()
	user := app.GetUser(ctx)

	var (
		title     = r.FormValue("Title")
		shortDesc = r.FormValue("ShortDesc")
		desc      = r.FormValue("Desc")
		imageURL  string
		start     pq.NullTime
		// assignment, _ = strconv.ParseBool(r.FormValue("Assignment"))
	)
	if len(title) == 0 {
		f.Add("Errors", "title required")
		back(w, r)
		return
	}

	if v := r.FormValue("Start"); len(v) > 0 {
		t, _ := time.Parse("2006-01-02", v)
		if !t.IsZero() {
			start.Time = t
			start.Valid = true
		}
	}

	image, info, err := r.FormFile("Image")
	if err != http.ErrMissingFile {
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}

		if !strings.Contains(info.Header.Get(header.ContentType), "image") {
			f.Add("Errors", "file is not an image")
			back(w, r)
			return
		}

		imageURL, err = c.uploadCourseCoverImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}
	}

	ctx, tx, err := app.NewTransactionContext(ctx)
	if err != nil {
		f.Add("Errors", err.Error())
		back(w, r)
		return
	}
	defer tx.Rollback()

	db := app.GetDatabase(ctx)
	var id string
	err = db.QueryRowContext(ctx, `
		insert into courses
			(user_id, title, short_desc, long_desc, image, start)
		values
			($1, $2, $3, $4, $5, $6)
		returning id
	`, user.ID, title, shortDesc, desc, imageURL, start).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = db.ExecContext(ctx, `
		insert into course_options
			(course_id)
		values
			($1)
	`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var link sql.NullString
	db.QueryRowContext(ctx, `select url from courses where id = $1`, id).Scan(&link)
	if !link.Valid {
		http.Redirect(w, r, "/course/"+id, http.StatusFound)
		return
	}
	http.Redirect(w, r, "/course/"+link.String, http.StatusFound)
}

func (c *ctrl) EditorCourse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method == http.MethodPost {
		c.postEditorCourse(w, r)
		return
	}
	id := r.FormValue("id")
	course, err := c.repo.GetCourse(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.view.EditorCourse(w, r, course)
}

func (c *ctrl) postEditorCourse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.FormValue("id")

	f := app.GetSession(ctx).Flash()

	var (
		title     = r.FormValue("Title")
		shortDesc = r.FormValue("ShortDesc")
		desc      = r.FormValue("Desc")
		imageURL  string
		start     pq.NullTime
		// assignment, _ = strconv.ParseBool(r.FormValue("Assignment"))
	)
	if len(title) == 0 {
		f.Add("Errors", "title required")
		back(w, r)
		return
	}

	if v := r.FormValue("Start"); len(v) > 0 {
		t, _ := time.Parse("2006-01-02", v)
		if !t.IsZero() {
			start.Time = t
			start.Valid = true
		}
	}

	image, info, err := r.FormFile("Image")
	if err != http.ErrMissingFile {
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}

		if !strings.Contains(info.Header.Get(header.ContentType), "image") {
			f.Add("Errors", "file is not an image")
			back(w, r)
			return
		}

		imageURL, err = c.uploadCourseCoverImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}
	}

	ctx, tx, err := app.NewTransactionContext(ctx)
	if err != nil {
		f.Add("Errors", err.Error())
		back(w, r)
		return
	}
	defer tx.Rollback()

	db := app.GetDatabase(ctx)
	_, err = db.ExecContext(ctx, `
		update courses
		set
			title = $2,
			short_desc = $3,
			long_desc = $4,
			start = $5,
			updated_at = now()
		where id = $1
	`, id, title, shortDesc, desc, start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(imageURL) > 0 {
		_, err = db.ExecContext(ctx, `
			update courses
			set
				image = $2
			where id = $1
		`, id, imageURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// _, err = tx.Exec(`
	// 	upsert into course_options
	// 		(course_id, assignment)
	// 	values
	// 		($1, $2)
	// `, id, assignment)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var link sql.NullString
	db.QueryRowContext(ctx, `select url from courses where id = $1`, id).Scan(&link)
	if !link.Valid {
		http.Redirect(w, r, "/course/"+id, http.StatusFound)
		return
	}
	http.Redirect(w, r, "/course/"+link.String, http.StatusSeeOther)
}

func (c *ctrl) EditorContent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.FormValue("id")

	if r.Method == http.MethodPost {
		if r.FormValue("action") == "delete" {
			db := app.GetDatabase(ctx)
			contentID := r.FormValue("contentId")
			_, err := db.ExecContext(ctx, `delete from course_contents where id = $1 and course_id = $2`, contentID, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		back(w, r)
		return
	}

	course, err := c.repo.GetCourse(ctx, id)
	if err == app.ErrNotFound {
		c.view.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	course.Contents, err = c.repo.GetCourseContents(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.view.EditorContent(w, r, course)
}

func (c *ctrl) CourseEnroll(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.postCourseEnroll(w, r)
		return
	}
	ctx := r.Context()
	user := app.GetUser(ctx)

	link := app.GetCourseURL(ctx)

	id := link
	_, err := uuid.Parse(link)
	if err != nil {
		id, err = c.repo.GetCourseIDFromURL(ctx, link)
		if err == app.ErrNotFound {
			c.view.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	x, err := c.repo.GetCourse(ctx, id)
	if err == app.ErrNotFound {
		c.view.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if user is course owner redirect back to course page
	if user.ID == x.UserID {
		http.Redirect(w, r, "/course/"+link, http.StatusFound)
		return
	}

	// redirect enrolled user back to course page
	enrolled, err := c.repo.IsEnrolled(ctx, user.ID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if enrolled {
		http.Redirect(w, r, "/course/"+link, http.StatusFound)
		return
	}

	// check is user has pending enroll
	pendingPayment, err := c.repo.HasPendingPayment(ctx, user.ID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if pendingPayment {
		http.Redirect(w, r, "/course/"+link, http.StatusFound)
		return
	}

	c.view.CourseEnroll(w, r, x)
}

func (c *ctrl) postCourseEnroll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := app.GetUser(ctx)
	f := app.GetSession(ctx).Flash()

	link := app.GetCourseURL(ctx)

	id := link
	_, err := uuid.Parse(link)
	if err != nil {
		id, err = c.repo.GetCourseIDFromURL(ctx, link)
		if err == app.ErrNotFound {
			c.view.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	x, err := c.repo.GetCourse(ctx, id)
	if err == app.ErrNotFound {
		c.view.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if user is course owner redirect back to course page
	if user.ID == x.UserID {
		http.Redirect(w, r, "/course/"+link, http.StatusFound)
		return
	}

	// redirect enrolled user back to course page
	enrolled, err := c.repo.IsEnrolled(ctx, user.ID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if enrolled {
		http.Redirect(w, r, "/course/"+link, http.StatusFound)
		return
	}

	// check is user has pending enroll
	pendingPayment, err := c.repo.HasPendingPayment(ctx, user.ID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if pendingPayment {
		http.Redirect(w, r, "/course/"+link, http.StatusFound)
		return
	}

	originalPrice := x.Price
	if x.Option.Discount {
		originalPrice = x.Discount
	}

	price, _ := strconv.ParseFloat(r.FormValue("Price"), 64)

	if price < 0 {
		f.Add("Errors", "price can not be negative")
		back(w, r)
		return
	}

	var imageURL string
	if originalPrice != 0 {
		image, info, err := r.FormFile("Image")
		if err == http.ErrMissingFile {
			f.Add("Errors", "image required")
			back(w, r)
			return
		}
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}
		if !strings.Contains(info.Header.Get(header.ContentType), "image") {
			f.Add("Errors", "file is not an image")
			back(w, r)
			return
		}

		imageURL, err = c.uploadPaymentImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}
	}

	newPayment := false

	ctx, tx, err := app.NewTransactionContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	if x.Price == 0 {
		err = c.repo.Enroll(ctx, user.ID, x.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err = c.repo.CreatePayment(ctx, &app.Payment{
			CourseID:      x.ID,
			UserID:        user.ID,
			Image:         imageURL,
			Price:         price,
			OriginalPrice: originalPrice,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newPayment = true
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if newPayment {
			c.sendSlackMessage(ctx, fmt.Sprintf("New payment for course %s, price %.2f", x.Title, price))
		} else {
			c.sendSlackMessage(ctx, fmt.Sprintf("New enroll for course %s", x.Title))
		}
		cancel()
	}()

	http.Redirect(w, r, "/course/"+link, http.StatusFound)
}

func (c *ctrl) CourseAssignment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := app.GetUser(ctx)
	link := app.GetCourseURL(ctx)

	// if id can parse to int64 get course from id
	id := link
	_, err := uuid.Parse(link)
	if err != nil {
		// link can not parse to int64 get course id from url
		id, err = c.repo.GetCourseIDFromURL(ctx, link)
		if err == app.ErrNotFound {
			c.view.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	x, err := c.repo.GetCourse(ctx, id)
	if err == app.ErrNotFound {
		c.view.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if course has url, redirect to course url
	if x.URL.Valid && x.URL.String != link {
		http.Redirect(w, r, "/course/"+x.URL.String+"/assignment", http.StatusFound)
		return
	}

	enrolled, err := c.repo.IsEnrolled(ctx, user.ID, x.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !enrolled && user.ID != x.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	assignments, err := c.repo.GetAssignments(ctx, x.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.view.Assignment(w, r, x, assignments)
}

func (c *ctrl) EditorContentCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.FormValue("id")

	if r.Method == http.MethodPost {
		var (
			title   = r.FormValue("Title")
			desc    = r.FormValue("Desc")
			videoID = r.FormValue("VideoID")
			i       int64
		)

		ctx, tx, err := app.NewTransactionContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		db := app.GetDatabase(ctx)
		// get content index
		err = db.QueryRowContext(ctx, `
			select i from course_contents where course_id = $1 order by i desc limit 1
		`, id).Scan(&i)
		if err == sql.ErrNoRows {
			i = -1
		}
		_, err = db.ExecContext(ctx, `
			insert into course_contents
				(course_id, i, title, long_desc, video_id, video_type)
			values
				($1, $2, $3, $4, $5, $6)
		`, id, i+1, title, desc, videoID, app.Youtube)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tx.Commit()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/editor/content?id="+r.FormValue("id"), http.StatusFound)
		return
	}

	course, err := c.repo.GetCourse(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.view.EditorContentCreate(w, r, course)
}

func (c *ctrl) EditorContentEdit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// course content id
	id := r.FormValue("id")

	content, err := c.repo.GetCourseContent(ctx, id)
	if err == sql.ErrNoRows {
		c.view.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	course, err := c.repo.GetCourse(ctx, content.CourseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := app.GetUser(r.Context())
	// user is not course owner
	if user.ID != course.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if r.Method == http.MethodPost {
		var (
			title   = r.FormValue("Title")
			desc    = r.FormValue("Desc")
			videoID = r.FormValue("VideoID")
		)

		db := app.GetDatabase(ctx)
		_, err = db.ExecContext(ctx, `
			update course_contents
			set
				title = $3,
				long_desc = $4,
				video_id = $5,
				updated_at = now()
			where id = $1 and course_id = $2
		`, id, course.ID, title, desc, videoID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/editor/content?id="+course.ID, http.StatusSeeOther)
		return
	}

	c.view.EditorContentEdit(w, r, course, content)
}
