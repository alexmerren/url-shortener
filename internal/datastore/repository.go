package datastore

type Datastorer interface {
	GetURL(id string) (url string, err error)
	InsertURL(url string) (id string, err error)
}
