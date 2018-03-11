package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/acoshift/header"
	"github.com/acoshift/hime"
	"github.com/acoshift/httprouter"
	"github.com/acoshift/pgsql"
	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/acoshift/acourse/appctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/repository"
)

func courseView(ctx hime.Context) hime.Result {
	user := appctx.GetUser(ctx)
	link := httprouter.GetParam(ctx, "courseURL")

	// if id can parse to uuid get course from id
	id := link
	_, err := uuid.Parse(link)
	if err != nil {
		// link can not parse to uuid get course id from url
		id, err = repository.GetCourseIDFromURL(db, link)
		if err == entity.ErrNotFound {
			return notFound(ctx)
		}
		must(err)
	}
	x, err := repository.GetCourse(db, id)
	if err == entity.ErrNotFound {
		return notFound(ctx)
	}
	must(err)

	// if course has url, redirect to course url
	if x.URL.Valid && x.URL.String != link {
		return ctx.RedirectTo("course", x.URL.String)
	}

	enrolled := false
	pendingEnroll := false
	if user != nil {
		enrolled, err = repository.IsEnrolled(db, user.ID, x.ID)
		must(err)

		if !enrolled {
			pendingEnroll, err = repository.HasPendingPayment(db, user.ID, x.ID)
			must(err)
		}
	}

	var owned bool
	if user != nil {
		owned = user.ID == x.UserID
	}

	// if user enrolled or user is owner fetch course contents
	if enrolled || owned {
		x.Contents, err = repository.GetCourseContents(db, x.ID)
		must(err)
	}

	if owned {
		x.Owner = user
	} else {
		x.Owner, err = repository.GetUser(db, x.UserID)
		must(err)
	}

	page := newPage(ctx)
	page["Title"] = x.Title + " | " + page["Title"].(string)
	page["Desc"] = x.ShortDesc
	page["Image"] = x.Image
	page["URL"] = baseURL + "/course/" + url.PathEscape(x.Link())
	page["Course"] = x
	page["Enrolled"] = enrolled
	page["Owned"] = owned
	page["pendingEnroll"] = pendingEnroll
	return ctx.View("course", page)
}

func courseContent(ctx hime.Context) hime.Result {
	user := appctx.GetUser(ctx)
	link := httprouter.GetParam(ctx, "courseURL")

	// if id can parse to uuid get course from id
	id := link
	_, err := uuid.Parse(link)
	if err != nil {
		// link can not parse to uuid get course id from url
		id, err = repository.GetCourseIDFromURL(db, link)
		if err == entity.ErrNotFound {
			return notFound(ctx)
		}
		must(err)
	}
	x, err := repository.GetCourse(db, id)
	if err == entity.ErrNotFound {
		return notFound(ctx)
	}
	must(err)

	// if course has url, redirect to course url
	if x.URL.Valid && x.URL.String != link {
		return ctx.RedirectTo("course", x.URL.String, "content")
	}

	enrolled, err := repository.IsEnrolled(db, user.ID, x.ID)
	must(err)

	if !enrolled && user.ID != x.UserID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	x.Contents, err = repository.GetCourseContents(db, x.ID)
	must(err)

	x.Owner, err = repository.GetUser(db, x.UserID)
	must(err)

	var content *entity.CourseContent
	p, _ := strconv.Atoi(ctx.FormValue("p"))
	if p < 0 {
		p = 0
	}
	if p > len(x.Contents)-1 {
		p = len(x.Contents) - 1
	}
	if p >= 0 {
		content = x.Contents[p]
	}

	page := newPage(ctx)
	page["Title"] = x.Title + " | " + page["Title"].(string)
	page["Desc"] = x.ShortDesc
	page["Image"] = x.Image
	page["Course"] = x
	page["Content"] = content
	return ctx.View("course.content", page)
}

func editorCreate(ctx hime.Context) hime.Result {
	return ctx.View("editor.create", newPage(ctx))
}

func postEditorCreate(ctx hime.Context) hime.Result {
	f := appctx.GetSession(ctx).Flash()
	user := appctx.GetUser(ctx)

	var (
		title     = ctx.FormValue("Title")
		shortDesc = ctx.FormValue("ShortDesc")
		desc      = ctx.FormValue("Desc")
		imageURL  string
		start     pq.NullTime
		// assignment, _ = strconv.ParseBool(ctx.FormValue("Assignment"))
	)
	if len(title) == 0 {
		f.Add("Errors", "title required")
		return ctx.RedirectToGet()
	}

	if v := ctx.FormValue("Start"); len(v) > 0 {
		t, _ := time.Parse("2006-01-02", v)
		if !t.IsZero() {
			start.Time = t
			start.Valid = true
		}
	}

	if image, info, err := ctx.FormFileNotEmpty("Image"); err != http.ErrMissingFile {
		if err != nil {
			f.Add("Errors", err.Error())
			return ctx.RedirectToGet()
		}

		if !strings.Contains(info.Header.Get(header.ContentType), "image") {
			f.Add("Errors", "file is not an image")
			return ctx.RedirectToGet()
		}

		imageURL, err = uploadCourseCoverImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			return ctx.RedirectToGet()
		}
	}

	var id string
	err := pgsql.RunInTx(db, nil, func(tx *sql.Tx) error {
		err := db.QueryRow(`
			insert into courses
				(user_id, title, short_desc, long_desc, image, start)
			values
				($1, $2, $3, $4, $5, $6)
			returning id
		`, user.ID, title, shortDesc, desc, imageURL, start).Scan(&id)
		if err != nil {
			return err
		}
		_, err = db.Exec(`
			insert into course_options
				(course_id)
			values
				($1)
		`, id)
		return err
	})
	must(err)

	var link sql.NullString
	db.QueryRow(`select url from courses where id = $1`, id).Scan(&link)
	if !link.Valid {
		return ctx.RedirectTo("course", id)
	}
	return ctx.RedirectTo("course", link.String)
}

func editorCourse(ctx hime.Context) hime.Result {
	id := ctx.FormValue("id")
	course, err := repository.GetCourse(db, id)
	must(err)

	page := newPage(ctx)
	page["Course"] = course
	return ctx.View("editor.course", page)
}

func postEditorCourse(ctx hime.Context) hime.Result {
	id := ctx.FormValue("id")

	f := appctx.GetSession(ctx).Flash()

	var (
		title     = ctx.FormValue("Title")
		shortDesc = ctx.FormValue("ShortDesc")
		desc      = ctx.FormValue("Desc")
		imageURL  string
		start     pq.NullTime
		// assignment, _ = strconv.ParseBool(ctx.FormValue("Assignment"))
	)
	if len(title) == 0 {
		f.Add("Errors", "title required")
		return ctx.RedirectToGet()
	}

	if v := ctx.FormValue("Start"); len(v) > 0 {
		t, _ := time.Parse("2006-01-02", v)
		if !t.IsZero() {
			start.Time = t
			start.Valid = true
		}
	}

	if image, info, err := ctx.FormFileNotEmpty("Image"); err != http.ErrMissingFile {
		if err != nil {
			f.Add("Errors", err.Error())
			return ctx.RedirectToGet()
		}

		if !strings.Contains(info.Header.Get(header.ContentType), "image") {
			f.Add("Errors", "file is not an image")
			return ctx.RedirectToGet()
		}

		imageURL, err = uploadCourseCoverImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			return ctx.RedirectToGet()
		}
	}

	err := pgsql.RunInTx(db, nil, func(tx *sql.Tx) error {
		_, err := tx.Exec(`
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
			return err
		}

		if len(imageURL) > 0 {
			_, err = tx.Exec(`
				update courses
				set
					image = $2
				where id = $1
			`, id, imageURL)
			if err != nil {
				return err
			}
		}
		// _, err = tx.Exec(`
		// 	upsert into course_options
		// 		(course_id, assignment)
		// 	values
		// 		($1, $2)
		// `, id, assignment)
		// must(err)
		return nil
	})
	must(err)

	var link sql.NullString
	db.QueryRow(`select url from courses where id = $1`, id).Scan(&link)
	if !link.Valid {
		return ctx.RedirectTo("course", id)
	}
	return ctx.RedirectTo("course", link.String)
}

func editorContent(ctx hime.Context) hime.Result {
	id := ctx.FormValue("id")

	course, err := repository.GetCourse(db, id)
	if err == entity.ErrNotFound {
		return notFound(ctx)
	}
	must(err)
	course.Contents, err = repository.GetCourseContents(db, id)
	must(err)

	page := newPage(ctx)
	page["Course"] = course
	return ctx.View("editor.content", page)
}

func postEditorContent(ctx hime.Context) hime.Result {
	id := ctx.FormValue("id")

	if ctx.FormValue("action") == "delete" {
		contentID := ctx.FormValue("contentId")
		_, err := db.Exec(`delete from course_contents where id = $1 and course_id = $2`, contentID, id)
		must(err)
	}
	return ctx.RedirectToGet()
}

func courseEnroll(ctx hime.Context) hime.Result {
	user := appctx.GetUser(ctx)

	link := httprouter.GetParam(ctx, "courseURL")

	id := link
	_, err := uuid.Parse(link)
	if err != nil {
		id, err = repository.GetCourseIDFromURL(db, link)
		if err == entity.ErrNotFound {
			return notFound(ctx)
		}
		must(err)
	}

	x, err := repository.GetCourse(db, id)
	if err == entity.ErrNotFound {
		return notFound(ctx)
	}
	must(err)

	// if user is course owner redirect back to course page
	if user.ID == x.UserID {
		return ctx.RedirectTo("course", link)
	}

	// redirect enrolled user back to course page
	enrolled, err := repository.IsEnrolled(db, user.ID, id)
	must(err)
	if enrolled {
		return ctx.RedirectTo("course", link)
	}

	// check is user has pending enroll
	pendingPayment, err := repository.HasPendingPayment(db, user.ID, id)
	must(err)
	if pendingPayment {
		return ctx.RedirectTo("course", link)
	}

	page := newPage(ctx)
	page["Title"] = x.Title + " | " + page["Title"].(string)
	page["Desc"] = x.ShortDesc
	page["Image"] = x.Image
	page["URL"] = baseURL + "/course/" + url.PathEscape(x.Link())
	page["Course"] = x
	return ctx.View("course.enroll", page)
}

func postCourseEnroll(ctx hime.Context) hime.Result {
	user := appctx.GetUser(ctx)
	f := appctx.GetSession(ctx).Flash()

	link := httprouter.GetParam(ctx, "courseURL")

	id := link
	_, err := uuid.Parse(link)
	if err != nil {
		id, err = repository.GetCourseIDFromURL(db, link)
		if err == entity.ErrNotFound {
			return notFound(ctx)
		}
		must(err)
	}

	x, err := repository.GetCourse(db, id)
	if err == entity.ErrNotFound {
		return notFound(ctx)
	}
	must(err)

	// if user is course owner redirect back to course page
	if user.ID == x.UserID {
		return ctx.RedirectTo("course", link)
	}

	// redirect enrolled user back to course page
	enrolled, err := repository.IsEnrolled(db, user.ID, id)
	must(err)
	if enrolled {
		return ctx.RedirectTo("course", link)
	}

	// check is user has pending enroll
	pendingPayment, err := repository.HasPendingPayment(db, user.ID, id)
	must(err)
	if pendingPayment {
		return ctx.RedirectTo("course", link)
	}

	originalPrice := x.Price
	if x.Option.Discount {
		originalPrice = x.Discount
	}

	price, _ := strconv.ParseFloat(ctx.FormValue("Price"), 64)

	if price < 0 {
		f.Add("Errors", "price can not be negative")
		return ctx.RedirectToGet()
	}

	var imageURL string
	if originalPrice != 0 {
		image, info, err := ctx.FormFileNotEmpty("Image")
		if err == http.ErrMissingFile {
			f.Add("Errors", "image required")
			return ctx.RedirectToGet()
		}
		if err != nil {
			f.Add("Errors", err.Error())
			return ctx.RedirectToGet()
		}
		if !strings.Contains(info.Header.Get(header.ContentType), "image") {
			f.Add("Errors", "file is not an image")
			return ctx.RedirectToGet()
		}

		imageURL, err = uploadPaymentImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			return ctx.RedirectToGet()
		}
	}

	newPayment := false

	err = pgsql.RunInTx(db, nil, func(tx *sql.Tx) error {
		if x.Price == 0 {
			err := repository.Enroll(tx, user.ID, x.ID)
			if err != nil {
				return err
			}
		} else {
			err = repository.CreatePayment(tx, &entity.Payment{
				CourseID:      x.ID,
				UserID:        user.ID,
				Image:         imageURL,
				Price:         price,
				OriginalPrice: originalPrice,
			})
			if err != nil {
				return err
			}
			newPayment = true
		}
		return nil
	})
	must(err)

	if newPayment {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			sendSlackMessage(ctx, fmt.Sprintf("New payment for course %s, price %.2f", x.Title, price))
		}()
	}

	return ctx.RedirectTo("course", link)
}

func courseAssignment(ctx hime.Context) hime.Result {
	user := appctx.GetUser(ctx)
	link := httprouter.GetParam(ctx, "courseURL")

	// if id can parse to int64 get course from id
	id := link
	_, err := uuid.Parse(link)
	if err != nil {
		// link can not parse to int64 get course id from url
		id, err = repository.GetCourseIDFromURL(db, link)
		if err == entity.ErrNotFound {
			return notFound(ctx)
		}
		must(err)
	}
	x, err := repository.GetCourse(db, id)
	if err == entity.ErrNotFound {
		return notFound(ctx)
	}
	must(err)

	// if course has url, redirect to course url
	if x.URL.Valid && x.URL.String != link {
		return ctx.RedirectTo("course", x.URL.String, "assignment")
	}

	enrolled, err := repository.IsEnrolled(db, user.ID, x.ID)
	must(err)

	if !enrolled && user.ID != x.UserID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	assignments, err := repository.GetAssignments(db, x.ID)
	must(err)

	page := newPage(ctx)
	page["Title"] = x.Title + " | " + page["Title"].(string)
	page["Desc"] = x.ShortDesc
	page["Image"] = x.Image
	page["URL"] = baseURL + "/course/" + url.PathEscape(x.Link())
	page["Course"] = x
	page["Assignments"] = assignments
	return ctx.View("assignment", page)
}

func editorContentCreate(ctx hime.Context) hime.Result {
	id := ctx.FormValue("id")

	course, err := repository.GetCourse(db, id)
	must(err)

	page := newPage(ctx)
	page["Course"] = course
	return ctx.View("editor.content.create", page)
}

func postEditorContentCreate(ctx hime.Context) hime.Result {
	id := ctx.FormValue("id")

	var (
		title   = ctx.FormValue("Title")
		desc    = ctx.FormValue("Desc")
		videoID = ctx.FormValue("VideoID")
		i       int64
	)

	err := pgsql.RunInTx(db, nil, func(tx *sql.Tx) error {
		// get content index
		err := tx.QueryRow(`
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
		`, id, i+1, title, desc, videoID, entity.Youtube)
		if err != nil {
			return err
		}
		return nil
	})
	must(err)

	return ctx.RedirectTo("editor.content", ctx.Param("id", ctx.FormValue("id")))
}

func editorContentEdit(ctx hime.Context) hime.Result {
	// course content id
	id := ctx.FormValue("id")

	content, err := repository.GetCourseContent(db, id)
	if err == sql.ErrNoRows {
		return notFound(ctx)
	}
	must(err)

	course, err := repository.GetCourse(db, content.CourseID)
	must(err)

	user := appctx.GetUser(ctx)
	// user is not course owner
	if user.ID != course.UserID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	page := newPage(ctx)
	page["Course"] = course
	page["Content"] = content
	return ctx.View("editor.content.edit", page)
}

func postEditorContentEdit(ctx hime.Context) hime.Result {
	// course content id
	id := ctx.FormValue("id")

	content, err := repository.GetCourseContent(db, id)
	if err == sql.ErrNoRows {
		return notFound(ctx)
	}
	must(err)

	course, err := repository.GetCourse(db, content.CourseID)
	must(err)

	user := appctx.GetUser(ctx)
	// user is not course owner
	if user.ID != course.UserID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	var (
		title   = ctx.FormValue("Title")
		desc    = ctx.FormValue("Desc")
		videoID = ctx.FormValue("VideoID")
	)

	_, err = db.Exec(`
		update course_contents
		set
			title = $3,
			long_desc = $4,
			video_id = $5,
			updated_at = now()
		where id = $1 and course_id = $2
	`, id, course.ID, title, desc, videoID)
	must(err)

	return ctx.RedirectTo("editor.content", ctx.Param("id", course.ID))
}
