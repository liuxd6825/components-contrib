package repository

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"go.mongodb.org/mongo-driver/mongo"
)

type session struct {
	client *mongo.Client
}

func StartSession(ctx context.Context, session eventstorage.Session, dbFunc eventstorage.SessionFunc) error {
	return session.UseTransaction(ctx, dbFunc)
}

func NewSession(client *mongo.Client) eventstorage.Session {
	return &session{client: client}
}

func (r *session) UseTransaction(ctx context.Context, dbFunc eventstorage.SessionFunc) error {
	return r.client.UseSession(ctx, func(sCtx mongo.SessionContext) (err error) {
		defer func() {
			if e := recover(); e != nil {
				if err1, ok := e.(error); ok {
					err = err1
				}
			}
		}()
		if err = sCtx.StartTransaction(); err != nil {
			return err
		}
		err = dbFunc(sCtx)
		if err != nil {
			if e1 := sCtx.AbortTransaction(sCtx); e1 != nil {
				err = e1
			}
		} else {
			err = sCtx.CommitTransaction(sCtx)
		}
		return err
	})
}
