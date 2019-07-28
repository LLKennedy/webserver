package app

// Server starts and only exits if there's an error
type Server interface {
	Start() error
}
