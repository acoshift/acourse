package app

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/acoshift/acourse/pkg/appctx"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/view"
	"github.com/acoshift/flash"
	"github.com/acoshift/header"
	"github.com/acoshift/httprouter"
	"github.com/lib/pq"
)

func getCourse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := appctx.GetUser(ctx)
	link := httprouter.GetParam(ctx, "courseID")

	// if id can parse to int64 get course from id
	id, err := strconv.ParseInt(link, 10, 64)
	if err != nil {
		// link can not parse to int64 get course id from url
		id, err = model.GetCourseIDFromURL(link)
		if err == model.ErrNotFound {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	x, err := model.GetCourse(id)
	if err == model.ErrNotFound {
		http.NotFound(w, r)
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
		enrolled, err = model.IsEnrolled(user.ID, x.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !enrolled {
			pendingEnroll, err = model.HasPendingPayment(user.ID, x.ID)
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
		x.Contents, err = model.GetCourseContents(x.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if owned {
		x.Owner = user
	} else {
		x.Owner, err = model.GetUser(x.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	view.Course(w, r, x, enrolled, owned, pendingEnroll)
}

func getCourseContent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := appctx.GetUser(ctx)
	link := httprouter.GetParam(ctx, "courseID")

	// if id can parse to int64 get course from id
	id, err := strconv.ParseInt(link, 10, 64)
	if err != nil {
		// link can not parse to int64 get course id from url
		id, err = model.GetCourseIDFromURL(link)
		if err == model.ErrNotFound {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	x, err := model.GetCourse(id)
	if err == model.ErrNotFound {
		http.NotFound(w, r)
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

	enrolled, err := model.IsEnrolled(user.ID, x.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !enrolled && user.ID != x.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	x.Contents, err = model.GetCourseContents(x.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	x.Owner, err = model.GetUser(x.UserID)
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

func getEditorCreate(w http.ResponseWriter, r *http.Request) {
	view.EditorCreate(w, r)
}

func postEditorCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f := flash.Get(ctx)
	user := appctx.GetUser(ctx)

	var (
		title         = r.FormValue("Title")
		shortDesc     = r.FormValue("ShortDesc")
		desc          = r.FormValue("Desc")
		imageURL      string
		start         pq.NullTime
		assignment, _ = strconv.ParseBool(r.FormValue("Assignment"))
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

		imageURL, err = UploadProfileImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}
	}

	tx, err := db.Begin()
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
			(course_id, assignment)
		values
			($1, $2)
	`, id, assignment)
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
	db.QueryRow(`select url from courses where id = $1`, id).Scan(&link)
	if !link.Valid {
		http.Redirect(w, r, "/course/"+strconv.FormatInt(id, 10), http.StatusFound)
		return
	}
	http.Redirect(w, r, "/course/"+link.String, http.StatusFound)
}

func getEditorCourse(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	course, err := model.GetCourse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.EditorCourse(w, r, course)
}

func postEditorCourse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)

	f := flash.Get(ctx)

	var (
		title         = r.FormValue("Title")
		shortDesc     = r.FormValue("ShortDesc")
		desc          = r.FormValue("Desc")
		imageURL      string
		start         pq.NullTime
		assignment, _ = strconv.ParseBool(r.FormValue("Assignment"))
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

		imageURL, err = UploadProfileImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}
	}

	tx, err := db.Begin()
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
	_, err = tx.Exec(`
		upsert into course_options
			(course_id, assignment)
		values
			($1, $2)
	`, id, assignment)
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
	db.QueryRow(`select url from courses where id = $1`, id).Scan(&link)
	if !link.Valid {
		http.Redirect(w, r, "/course/"+strconv.FormatInt(id, 10), http.StatusFound)
		return
	}
	http.Redirect(w, r, "/course/"+link.String, http.StatusSeeOther)
}

func getEditorContent(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	course, err := model.GetCourse(id)
	if err == model.ErrNotFound {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	course.Contents, err = model.GetCourseContents(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.EditorContent(w, r, course)
}

func getCourseEnroll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := appctx.GetUser(ctx)

	link := httprouter.GetParam(ctx, "courseID")

	id, err := strconv.ParseInt(link, 10, 64)
	if err != nil {
		id, err = model.GetCourseIDFromURL(link)
		if err == model.ErrNotFound {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	x, err := model.GetCourse(id)
	if err == model.ErrNotFound {
		http.NotFound(w, r)
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
	enrolled, err := model.IsEnrolled(user.ID, id)
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
	f := flash.Get(ctx)

	link := httprouter.GetParam(ctx, "courseID")

	id, err := strconv.ParseInt(link, 10, 64)
	if err != nil {
		id, err = model.GetCourseIDFromURL(link)
		if err == model.ErrNotFound {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	x, err := model.GetCourse(id)
	if err == model.ErrNotFound {
		http.NotFound(w, r)
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
	enrolled, err := model.IsEnrolled(user.ID, id)
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

	priceStr := r.FormValue("Price")
	if len(priceStr) == 0 && originalPrice != 0 {
		f.Add("Errors", "price can not be empty")
		back(w, r)
		return
	}
	price, _ := strconv.ParseFloat(priceStr, 64)

	if price < 0 {
		f.Add("Errors", "price can not be negative")
		back(w, r)
		return
	}

	var imageURL string
	if x.Price != 0 {
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

			imageURL, err = UploadPaymentImage(ctx, image)
			if err != nil {
				f.Add("Errors", err.Error())
				back(w, r)
				return
			}
		}
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	if x.Price == 0 {
		err = model.Enroll(tx, user.ID, x.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		paymentID, err := model.CreatePayment(tx, &model.Payment{
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
		// TODO: send email to user
		_ = paymentID
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/course/"+link, http.StatusFound)
}

func getEditorContentCreate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	course, err := model.GetCourse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	view.EditorContentCreate(w, r, course)
}

func getEditorContentEdit(w http.ResponseWriter, r *http.Request) {
	// course content id
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)

	content, err := model.GetCourseContent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	course, err := model.GetCourse(content.CourseID)
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

	view.EditorContentEdit(w, r, course, content)
}
