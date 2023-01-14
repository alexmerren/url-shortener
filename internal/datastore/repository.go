package datastore

import "context"

type UrlStorer interface {
	Start() error
	Stop() error
	GetUrl(id string) (url string, err error)
	InsertUrl(url string) (id string, err error)
	GetUrlWithContext(ctx context.Context, id string) (url string, err error)
	InsertUrlWithContext(ctx context.Context, url string) (id string, err error)
}
