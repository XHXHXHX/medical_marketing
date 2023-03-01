package core

type logData struct {
	namespace string
	project   string
	logger    string
	linenum   string
	ts        int64
	level     Level
	msg       string
	args      map[string]interface{}
}
