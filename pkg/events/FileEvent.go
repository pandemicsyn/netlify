package events

//StatusCreated indicates the file was just created
const StatusCreated = "created"

//FileEvent is the message we distribute to down stream subscribers
type FileEvent struct {
	Bucket  string
	Object  string
	Status  string
	Version int
}
