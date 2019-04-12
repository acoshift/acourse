package course

import (
	"bytes"
	"context"
	"database/sql"
	"io"

	"github.com/acoshift/pgsql/pgctx"

	"github.com/acoshift/acourse/internal/pkg/app"
	"github.com/acoshift/acourse/internal/pkg/file"
	"github.com/acoshift/acourse/internal/pkg/image"
)

// Content type
type Content struct {
	ID          string
	CourseID    string
	Title       string
	Desc        string
	VideoID     string
	VideoType   int
	DownloadURL string
}

type CreateContentArgs struct {
	ID        string
	Title     string
	LongDesc  string
	VideoID   string
	VideoType int
}

// CreateContent creates new course content
func CreateContent(ctx context.Context, m *CreateContentArgs) (string, error) {
	// TODO: validate instructor

	var contentID string
	err := pgctx.QueryRow(ctx, `
		insert into course_contents
			(
				course_id,
				i,
				title, long_desc, video_id, video_type
			)
		values
			(
				$1,
				(select coalesce(max(i)+1, 0) from course_contents where course_id = $1),
				$2, $3, $4, $5
			)
		returning id
	`,
		m.ID,
		m.Title, m.LongDesc, m.VideoID, m.VideoType,
	).Scan(&contentID)
	return contentID, err
}

type UpdateContentArgs struct {
	ContentID string
	Title     string
	Desc      string
	VideoID   string
}

// UpdateContent updates a course content
func UpdateContent(ctx context.Context, m *UpdateContentArgs) error {
	// TODO: validate ownership

	_, err := pgctx.Exec(ctx, `
		update course_contents
		set
			title = $2,
			long_desc = $3,
			video_id = $4,
			updated_at = now()
		where id = $1
	`, m.ContentID, m.Title, m.Desc, m.VideoID)
	return err
}

// GetContent gets a course's content
func GetContent(ctx context.Context, contentID string) (*Content, error) {
	// TODO: validate ownership

	var x Content
	err := pgctx.QueryRow(ctx, `
		select
			id, course_id, title, long_desc, video_id, video_type, download_url
		from course_contents
		where id = $1
	`, contentID).Scan(
		&x.ID, &x.CourseID, &x.Title, &x.Desc, &x.VideoID, &x.VideoType, &x.DownloadURL,
	)
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// DeleteContent deletes a course's content
func DeleteContent(ctx context.Context, contentID string) error {
	// TODO: validate ownership

	_, err := pgctx.Exec(ctx, `delete from course_contents where id = $1`, contentID)
	return err
}

// GetContents gets course's contents
func GetContents(ctx context.Context, id string) ([]*Content, error) {
	// TODO: validate ownership

	rows, err := pgctx.Query(ctx, `
		select
			id, course_id, title, long_desc, video_id, video_type, download_url
		from course_contents
		where course_id = $1
		order by i
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var xs []*Content
	for rows.Next() {
		var x Content
		err = rows.Scan(
			&x.ID, &x.CourseID, &x.Title, &x.Desc, &x.VideoID, &x.VideoType, &x.DownloadURL,
		)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return xs, nil
}

func uploadCourseCoverImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := image.JPEG(buf, r, 1200, 0, 90, false)
	if err != nil {
		return "", err
	}

	filename := file.GenerateFilename() + ".jpg"

	downloadURL, err := file.Store(ctx, buf, filename, false)
	if err != nil {
		return "", err
	}
	return downloadURL, nil
}
