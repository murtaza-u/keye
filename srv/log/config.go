package log

// Config holds configuration options passed to the logger.
type Config struct {
	// Debug enables debug logs.
	Debug bool
	// UseJSON configures logger to use json instead of plain text.
	UseJSON bool
}
