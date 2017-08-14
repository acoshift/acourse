package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/acoshift/acourse/pkg/appctx"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/view"
	"github.com/acoshift/header"
	"github.com/acoshift/session"
	"github.com/lib/pq"
)

func course(w http.ResponseWriter, r *http.Request) {
	s := strings.SplitN(r.URL.Path, "/", 2)
	var p string
	if len(s) > 1 {
		p = strings.TrimSuffix(s[1], "/")
	}

	r = r.WithContext(appctx.WithCourseURL(r.Context(), s[0]))
	switch p {
	case "":
		courseView(w, r)
	case "content":
		mustSignedIn(http.HandlerFunc(courseContent)).ServeHTTP(w, r)
	case "enroll":
		mustSignedIn(http.HandlerFunc(courseEnroll)).ServeHTTP(w, r)
	case "assignment":
		mustSignedIn(http.HandlerFunc(courseAssignment)).ServeHTTP(w, r)
	default:
		view.NotFound(w, r)
	}
}

func courseView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := appctx.GetUser(ctx)
	link := appctx.GetCourseURL(ctx)

	// if id can parse to uuid get course from id
	id := intID(link)
	var err error
	if id == 0 {
		// link can not parse to uuid get course id from url
		id, err = model.GetCourseIDFromURL(ctx, db, link)
		if err == model.ErrNotFound {
			view.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	x, err := model.GetCourse(ctx, db, id)
	if err == model.ErrNotFound {
		view.NotFound(w, r)
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
		enrolled, err = model.IsEnrolled(ctx, db, user.ID, x.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !enrolled {
			pendingEnroll, err = model.HasPendingPayment(ctx, db, user.ID, x.ID)
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
		x.Contents, err = model.GetCourseContents(ctx, db, x.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if owned {
		x.Owner = user
	} else {
		x.Owner, err = model.GetUser(ctx, db, x.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	view.Course(w, r, x, enrolled, owned, pendingEnroll)
}

func courseContent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := appctx.GetUser(ctx)
	link := appctx.GetCourseURL(ctx)

	// if id can parse to uuid get course from id
	id := intID(link)
	var err error
	if id == 0 {
		// link can not parse to uuid get course id from url
		id, err = model.GetCourseIDFromURL(ctx, db, link)
		if err == model.ErrNotFound {
			view.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	x, err := model.GetCourse(ctx, db, id)
	if err == model.ErrNotFound {
		view.NotFound(w, r)
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

	enrolled, err := model.IsEnrolled(ctx, db, user.ID, x.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !enrolled && user.ID != x.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	x.Contents, err = model.GetCourseContents(ctx, db, x.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	x.Owner, err = model.GetUser(ctx, db, x.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var content *model.CourseContent
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

	view.CourseContent(w, r, x, content)
}

func editorCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postEditorCreate(w, r)
		return
	}
	view.EditorCreate(w, r)
}

func postEditorCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f := session.Get(ctx, sessName).Flash()
	user := appctx.GetUser(ctx)

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

		imageURL, err = UploadCourseCoverImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		f.Add("Errors", err.Error())
		back(w, r)
		return
	}
	defer tx.Rollback()

	var id int64
	err = tx.QueryRow(`
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
	_, err = tx.Exec(`
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
	db.QueryRowContext(ctx, `select url from courses where id = $1`, strID(id)).Scan(&link)
	if !link.Valid {
		http.Redirect(w, r, "/course/"+strID(id), http.StatusFound)
		return
	}
	http.Redirect(w, r, "/course/"+link.String, http.StatusFound)
}

func editorCourse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method == http.MethodPost {
		postEditorCourse(w, r)
		return
	}
	id := intID(r.FormValue("id"))
	course, err := model.GetCourse(ctx, db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.EditorCourse(w, r, course)
}

func postEditorCourse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.FormValue("id")

	f := session.Get(ctx, sessName).Flash()

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

		imageURL, err = UploadCourseCoverImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		f.Add("Errors", err.Error())
		back(w, r)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
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
		_, err = tx.Exec(`
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

func editorContent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := intID(r.FormValue("id"))

	if r.Method == http.MethodPost {
		if r.FormValue("action") == "delete" {
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

	course, err := model.GetCourse(ctx, db, id)
	if err == model.ErrNotFound {
		view.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	course.Contents, err = model.GetCourseContents(ctx, db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.EditorContent(w, r, course)
}

func courseEnroll(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postCourseEnroll(w, r)
		return
	}
	ctx := r.Context()
	user := appctx.GetUser(ctx)

	link := appctx.GetCourseURL(ctx)

	id := intID(link)
	var err error
	if id == 0 {
		id, err = model.GetCourseIDFromURL(ctx, db, link)
		if err == model.ErrNotFound {
			view.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	x, err := model.GetCourse(ctx, db, id)
	if err == model.ErrNotFound {
		view.NotFound(w, r)
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
	enrolled, err := model.IsEnrolled(ctx, db, user.ID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if enrolled {
		http.Redirect(w, r, "/course/"+link, http.StatusFound)
		return
	}

	view.CourseEnroll(w, r, x)
}

func postCourseEnroll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := appctx.GetUser(ctx)
	f := session.Get(ctx, sessName).Flash()

	link := appctx.GetCourseURL(ctx)

	id := intID(link)
	var err error
	if id == 0 {
		id, err = model.GetCourseIDFromURL(ctx, db, link)
		if err == model.ErrNotFound {
			view.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	x, err := model.GetCourse(ctx, db, id)
	if err == model.ErrNotFound {
		view.NotFound(w, r)
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
	enrolled, err := model.IsEnrolled(ctx, db, user.ID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if enrolled {
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

		imageURL, err = UploadPaymentImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}
	}

	newPayment := false
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	if x.Price == 0 {
		err = model.Enroll(ctx, tx, user.ID, x.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err = model.CreatePayment(ctx, tx, &model.Payment{
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
			sendSlackMessage(ctx, fmt.Sprintf("New payment for course %s, price %.2f", x.Title, price))
		} else {
			sendSlackMessage(ctx, fmt.Sprintf("New enroll for course %s", x.Title))
		}
		cancel()
	}()

	http.Redirect(w, r, "/course/"+link, http.StatusFound)
}

func courseAssignment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := appctx.GetUser(ctx)
	link := appctx.GetCourseURL(ctx)

	// if id can parse to int64 get course from id
	id := intID(link)
	var err error
	if id == 0 {
		// link can not parse to int64 get course id from url
		id, err = model.GetCourseIDFromURL(ctx, db, link)
		if err == model.ErrNotFound {
			view.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	x, err := model.GetCourse(ctx, db, id)
	if err == model.ErrNotFound {
		view.NotFound(w, r)
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

	enrolled, err := model.IsEnrolled(ctx, db, user.ID, x.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !enrolled && user.ID != x.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	assignments, err := model.GetAssignments(ctx, db, x.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.Assignment(w, r, x, assignments)
}

func editorContentCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := intID(r.FormValue("id"))

	if r.Method == http.MethodPost {
		var (
			title   = r.FormValue("Title")
			desc    = r.FormValue("Desc")
			videoID = r.FormValue("VideoID")
			i       int64
		)

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		// get content index
		err = tx.QueryRow(`
			select i from course_contents where course_id = $1 order by i desc limit 1
		`, id).Scan(&i)
		if err == sql.ErrNoRows {
			i = -1
		}
		_, err = tx.Exec(`
			insert into course_contents
				(course_id, i, title, long_desc, video_id, video_type)
			values
				($1, $2, $3, $4, $5, $6)
		`, id, i+1, title, desc, videoID, model.Youtube)
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

	course, err := model.GetCourse(ctx, db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	view.EditorContentCreate(w, r, course)
}

func editorContentEdit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// course content id
	id := intID(r.FormValue("id"))

	content, err := model.GetCourseContent(ctx, db, id)
	if err == sql.ErrNoRows {
		view.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	course, err := model.GetCourse(ctx, db, content.CourseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := appctx.GetUser(r.Context())
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
		http.Redirect(w, r, "/editor/content?id="+strID(course.ID), http.StatusSeeOther)
		return
	}

	view.EditorContentEdit(w, r, course, content)
}
