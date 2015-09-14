package middleware

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
)

type Logger struct {
	AppName   string
	ErrHeader string
}

func NewLogger(appName, errHeader string) *Logger {
	return &Logger{
		AppName:   appName,
		ErrHeader: errHeader,
	}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)

	took := time.Since(start)
	ngrw := rw.(negroni.ResponseWriter)

	log := logrus.WithFields(logrus.Fields{
		"app":     l.AppName,
		"request": r.RequestURI,
		"action":  r.Method,
		"remote":  r.RemoteAddr,
		"status":  ngrw.Status(),
		"took":    took,
	})
	if ngrw.Status() != 200 {
		log.Errorln(ngrw.Header().Get(l.ErrHeader))
		return
	}
	log.Infoln("Request Complete")
}

func (l *Logger) WriteErrHeader(rw *http.ResponseWriter, err *error) {
	(*rw).Header().Add(l.ErrHeader, (*err).Error())
	(*rw).WriteHeader(http.StatusInternalServerError)
}

func (l *Logger) LogFatal(appmsg string, err error) {
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"app":    l.AppName,
			"appmsg": appmsg,
		}).Fatalln(err)
	}
}

func (l *Logger) LogError(appmsg, msgBody string, err error) {
	logrus.WithFields(logrus.Fields{
		"app":     l.AppName,
		"msgbody": msgBody,
		"appmsg":  appmsg,
	}).Errorln(err)
}
