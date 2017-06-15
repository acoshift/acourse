package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	time.Local = time.UTC

	db1, err := sql.Open("postgres", "postgresql://acourse_dev@db-0.cluster.acoshift.com:26257/acourse_dev?sslmode=verify-full&sslcert=private%2Fclient.crt&sslkey=private%2Fclient.key&sslrootcert=private%2Fca.crt")
	must(err)
	db2, err := sql.Open("postgres", "")
	must(err)

	_, err = db2.Exec("delete from payments")
	must(err)
	_, err = db2.Exec("delete from user_assignments")
	must(err)
	_, err = db2.Exec("delete from assignments")
	must(err)
	_, err = db2.Exec("delete from enrolls")
	must(err)
	_, err = db2.Exec("delete from course_contents")
	must(err)
	_, err = db2.Exec("delete from course_options")
	must(err)
	_, err = db2.Exec("delete from attends")
	must(err)
	_, err = db2.Exec("delete from courses")
	must(err)
	_, err = db2.Exec("delete from roles")
	must(err)
	_, err = db2.Exec("delete from users")
	must(err)

	log.Println("migrate users")
	stmt, err := db2.Prepare(`
		insert into users
			(id, username, name, email, about_me, image, created_at, updated_at)
		values
			($1, $2, $3, $4, $5, $6, $7, $8)
	`)
	must(err)
	rows, err := db1.Query(`
		select id, username, name, email, about_me, image, created_at, updated_at
		from users
	`)
	must(err)
	for rows.Next() {
		var id, username, name, email, aboutMe, image, createdAt, updatedAt interface{}
		err = rows.Scan(&id, &username, &name, &email, &aboutMe, &image, &createdAt, &updatedAt)
		must(err)
		_, err = stmt.Exec(id, username, name, email, aboutMe, image, createdAt, updatedAt)
		must(err)
	}
	rows.Close()
	stmt.Close()

	log.Println("migrate role")
	stmt, err = db2.Prepare(`
		insert into roles
			(user_id, admin, instructor, created_at, updated_at)
		values
			($1, $2, $3, $4, $5);
	`)
	must(err)
	rows, err = db1.Query(`
		select user_id, admin, instructor, created_at, updated_at
		from roles
	`)
	for rows.Next() {
		var userID, admin, instructor, createdAt, updatedAt interface{}
		err = rows.Scan(&userID, &admin, &instructor, &createdAt, &updatedAt)
		must(err)
		_, err = stmt.Exec(userID, admin, instructor, createdAt, updatedAt)
		must(err)
	}
	rows.Close()
	stmt.Close()

	log.Println("migrate courses")
	stmt, err = db2.Prepare(`
		insert into courses
			(user_id, title, short_desc, long_desc, image, start, url, type, price, discount, enroll_detail, created_at, updated_at)
		values
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		returning id
	`)
	must(err)
	rows, err = db1.Query(`
		select id, user_id, title, short_desc, long_desc, image, start, url, type, price, discount, enroll_detail, created_at, updated_at
		from courses
	`)
	must(err)
	mapCourseID := make(map[int64]string)
	for rows.Next() {
		var id int64
		var userID, title, shortDesc, longDesc, image, start, url, typ, price, discount, enrollDetail, createdAt, updatedAt interface{}
		err = rows.Scan(&id, &userID, &title, &shortDesc, &longDesc, &image, &start, &url, &typ, &price, &discount, &enrollDetail, &createdAt, &updatedAt)
		must(err)
		var newID string
		err = stmt.QueryRow(userID, title, shortDesc, longDesc, image, start, url, typ, price, discount, enrollDetail, createdAt, updatedAt).Scan(&newID)
		must(err)
		mapCourseID[id] = newID
	}
	rows.Close()
	stmt.Close()

	log.Println("migrate course_options")
	stmt, err = db2.Prepare(`
		insert into course_options
			(course_id, public, enroll, attend, assignment, discount)
		values
			($1, $2, $3, $4, $5, $6);
	`)
	must(err)
	rows, err = db1.Query(`
		select course_id, public, enroll, attend, assignment, discount
		from course_options
	`)
	must(err)
	for rows.Next() {
		var courseID int64
		var public, enroll, attend, assignment, discount interface{}
		err = rows.Scan(&courseID, &public, &enroll, &attend, &assignment, &discount)
		must(err)
		_, err = stmt.Exec(mapCourseID[courseID], public, enroll, attend, assignment, discount)
		must(err)
	}
	rows.Close()
	stmt.Close()

	log.Println("migrate course_contents")
	stmt, err = db2.Prepare(`
		insert into course_contents
			(course_id, i, title, long_desc, video_id, video_type, download_url)
		values
			($1, $2, $3, $4, $5, $6, $7);
	`)
	must(err)
	rows, err = db1.Query(`
		select course_id, i, title, long_desc, video_id, video_type, download_url
		from course_contents
	`)
	for rows.Next() {
		var courseID int64
		var i, title, longDesc, videoID, videoType, downloadURL interface{}
		err = rows.Scan(&courseID, &i, &title, &longDesc, &videoID, &videoType, &downloadURL)
		must(err)
		_, err = stmt.Exec(mapCourseID[courseID], i, title, longDesc, videoID, videoType, downloadURL)
		must(err)
	}
	rows.Close()
	stmt.Close()

	log.Println("migrate payments")
	stmt, err = db2.Prepare(`
		insert into payments
			(user_id, course_id, image, price, original_price, code, status, created_at, updated_at, at)
		values
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`)
	must(err)
	rows, err = db1.Query(`
		select user_id, course_id, image, price, original_price, code, status, created_at, updated_at, at
		from payments
	`)
	must(err)
	for rows.Next() {
		var courseID int64
		var userID, image, price, originalPrice, code, status, createdAt, updatedAt, at interface{}
		err = rows.Scan(&userID, &courseID, &image, &price, &originalPrice, &code, &status, &createdAt, &updatedAt, &at)
		must(err)
		_, err = stmt.Exec(userID, mapCourseID[courseID], image, price, originalPrice, code, status, createdAt, updatedAt, at)
		must(err)
	}
	rows.Close()
	stmt.Close()

	log.Println("migrate assignments")
	stmt, err = db2.Prepare(`
		insert into assignments
			(course_id, i, title, long_desc, open, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7)
		returning id;
	`)
	must(err)
	rows, err = db1.Query(`
		select id, course_id, i, title, long_desc, open, created_at, updated_at
		from assignments
	`)
	must(err)
	mapAssignmentID := make(map[int64]string)
	for rows.Next() {
		var id int64
		var courseID int64
		var i, title, longDesc, open, createdAt, updatedAt interface{}
		err = rows.Scan(&id, &courseID, &i, &title, &longDesc, &open, &createdAt, &updatedAt)
		must(err)
		var newID string
		err = stmt.QueryRow(mapCourseID[courseID], i, title, longDesc, open, createdAt, updatedAt).Scan(&newID)
		must(err)
		mapAssignmentID[id] = newID
	}
	rows.Close()
	stmt.Close()

	log.Println("migrate enroll")
	stmt, err = db2.Prepare(`
		insert into enrolls
			(user_id, course_id, created_at)
		values
			($1, $2, $3);
	`)
	must(err)
	rows, err = db1.Query(`
		select user_id, course_id, created_at
		from enrolls
	`)
	must(err)
	for rows.Next() {
		var courseID int64
		var userID, createdAt interface{}
		err = rows.Scan(&userID, &courseID, &createdAt)
		must(err)
		_, err = stmt.Exec(userID, mapCourseID[courseID], createdAt)
		must(err)
	}
	rows.Close()
	stmt.Close()

	log.Println("migrate attend")
	stmt, err = db2.Prepare(`
		insert into attends
			(user_id, course_id, created_at)
		values
			($1, $2, $3);
	`)
	must(err)
	rows, err = db1.Query(`
		select user_id, course_id, created_at
		from attends
	`)
	must(err)
	for rows.Next() {
		var courseID int64
		var userID, createdAt interface{}
		err = rows.Scan(&userID, &courseID, &createdAt)
		must(err)
		_, err = stmt.Exec(userID, mapCourseID[courseID], createdAt)
		must(err)
	}
	rows.Close()
	stmt.Close()

	log.Println("migrate user assignments")
	stmt, err = db2.Prepare(`
		insert into user_assignments
			(user_id, assignment_id, download_url, created_at)
		values
			($1, $2, $3, $4);
	`)
	must(err)
	rows, err = db1.Query(`
		select user_id, assignment_id, download_url, created_at
		from user_assignments
	`)
	must(err)
	for rows.Next() {
		var id int64
		var userID, downloadURL, createdAt interface{}
		err = rows.Scan(&userID, &id, &downloadURL, &createdAt)
		must(err)
		_, err = stmt.Exec(userID, mapAssignmentID[id], downloadURL, createdAt)
		must(err)
	}
	rows.Close()
	stmt.Close()
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
