package internal

type FileEvent struct {
	Bucket  string
	Object  string
	Status  string
	Version int
}