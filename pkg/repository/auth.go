package repository

import "context"

func (repo) StoreMagicLink(ctx context.Context, linkID string, userID string) error {
	return nil
}

func (repo) FindMagicLink(ctx context.Context, linkID string) (string, error) {
	return "", nil
}
