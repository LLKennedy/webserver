package logs

// Logger write logs to some location
type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}
