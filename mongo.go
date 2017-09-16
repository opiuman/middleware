package middleware

import (
	"context"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
)

type contextKey string

//MongoDB represents the current mgo session and database.
type MongoDB struct {
	Session *mgo.Session
	DB      string
}

//NewDB initialize new MongoDB from config
func NewDB(addrs []string, db, user, pwd string, log *Logger) *MongoDB {
	mdb := MongoDB{
		DB: db,
	}
	info := &mgo.DialInfo{
		Addrs:    addrs,
		Timeout:  10 * time.Second,
		Database: db,
		Username: user,
		Password: pwd,
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		log.Fatalf("failed to connect to mongodb at %s -- %s", addrs, err)
	}
	mdb.Session = session

	return &mdb
}

//ServeHTTP adds database to context
func (mdb *MongoDB) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	reqSession := mdb.Session.Clone()
	defer reqSession.Close()
	db := reqSession.DB(mdb.DB)
	ctx := context.WithValue(req.Context(), contextKey(mdb.DB), db)
	next(rw, req.WithContext(ctx))
}

//GetDB returns the database from context
func (mdb *MongoDB) GetDB(req *http.Request) *mgo.Database {
	if reqv := req.Context().Value(contextKey(mdb.DB)); reqv != nil {
		return reqv.(*mgo.Database)
	}
	return nil
}

//Close closes the active mgo session
func (mdb *MongoDB) Close() {
	mdb.Session.Close()
}
