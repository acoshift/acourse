package user

import (
	"context"
	"log"
	"time"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/ds"
	"github.com/acoshift/gotcha"
)

var cacheRole = gotcha.New()

func (s *service) startCacheRole() {
	ctx := context.Background()
	for {
		var xs []*role
		err := s.client.Query(ctx, kindRole, &xs)
		err = ds.IgnoreFieldMismatch(err)
		if err != nil {
			time.Sleep(time.Minute * 10)
			continue
		}
		cacheRole.Purge()
		for _, x := range xs {
			cacheRole.Set(x.ID(), x)
		}
		log.Println("Cached Roles")
		time.Sleep(time.Hour)
	}
}

func (s *service) GetRole(ctx context.Context, req *acourse.UserIDRequest) (*acourse.Role, error) {
	if req.GetUserId() == "" {
		return &acourse.Role{}, nil
	}

	x, _ := cacheRole.Get(req.GetUserId()).(*role)
	if x == nil {
		x = &role{}
	}
	return &acourse.Role{
		Admin:      x.Admin,
		Instructor: x.Instructor,
	}, nil
}
