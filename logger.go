package middleware

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/negroni"
)

//Logger is a wrapper of logrus entry
type Logger struct {
	*logrus.Entry
	ErrHeader string
}

//NewLogger returns Logger
func NewLogger(appName string) *Logger {
	return &Logger{logrus.WithField("application", appName), appName + "-error"}
}

//ServeHTTP logs the request info
func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)

	took := time.Since(start)
	ngrw := rw.(negroni.ResponseWriter)

	status := ngrw.Status()
	log := l.WithFields(logrus.Fields{
		"request": r.RequestURI,
		"action":  r.Method,
		"remote":  r.RemoteAddr,
		"status":  status,
		"took":    took,
	})
	if status >= 400 {
		log.Errorln(ngrw.Header().Get(l.ErrHeader))
		return
	}
	if ngrw.Header().Get("info") != "" {
		log.Info(ngrw.Header().Get("info"))
		return
	}
	log.Infof("%d OK", status)
}

//WriteErrHeader inserts the error to repsonsewriter header
func (l *Logger) WriteErrHeader(rw *http.ResponseWriter, err *error, status int) {
	(*rw).Header().Add(l.ErrHeader, (*err).Error())
	(*rw).WriteHeader(status)
}

//WriteInfoHeader inserts the info to repsonsewriter header
func (l *Logger) WriteInfoHeader(rw *http.ResponseWriter, info string) {
	(*rw).Header().Add("info", info)
	(*rw).WriteHeader(http.StatusOK)
}
