package middleware

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
)

type MongoDB struct {
	Session *mgo.Session
	DB      string
	Log     *Logger
}

func NewDB(mgoSrv, db string, log *Logger) *MongoDB {
	mdb := MongoDB{
		DB: db,
	}
	session, err := mgo.Dial(mgoSrv)
	log.LogFatal(fmt.Sprintf("failed to connect to mongodb at %s", mgoSrv), err)
	mdb.Session = session

	return &mdb
}

func (mdb *MongoDB) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	reqSession := mdb.Session.Clone()
	defer reqSession.Close()
	db := reqSession.DB(mdb.DB)
	context.Set(req, mdb.DB, db)
	next(rw, req)
}

