package mocklog

import (
	"bytes"
	"fmt"
)

// New returns a new mock logger
func New() *MockLogger {
	m := new(MockLogger)
	m.buffer = bytes.NewBuffer(nil)
	return m
}

// MockLogger is a mock logger
type MockLogger struct {
	buffer     *bytes.Buffer
	fatalCount int
}

// Println prints a line
func (m *MockLogger) Println(v ...interface{}) {
	fmt.Fprintln(m.buffer, v...)
}

// Printf prints formatted text
func (m *MockLogger) Printf(format string, v ...interface{}) {
	fmt.Fprintf(m.buffer, format, v...)
}

// Fatalf calls Printf then os.Exit(1)
func (m *MockLogger) Fatalf(format string, v ...interface{}) {
	m.Printf(format, v...)
	m.fatalCount++
}

// GetBuffer returns the internal buffer
func (m *MockLogger) GetBuffer() *bytes.Buffer {
	if m == nil {
		return nil
	}
	return m.buffer
}

// GetFatalCount returns the internal fatal count
func (m *MockLogger) GetFatalCount() int {
	if m == nil {
		return 0
	}
	return m.fatalCount
}

// SetFatalCount changes the internal fatal count
func (m *MockLogger) SetFatalCount(newVal int) {
	if m != nil {
		m.fatalCount = newVal
	}
}

// GetContents returns the buffer contents as a string
func (m *MockLogger) GetContents() string {
	if m.GetBuffer() == nil {
		return ""
	}
	return string(m.GetBuffer().Bytes())
}
