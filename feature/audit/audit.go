package audit

import (
	"context"
	"database/sql"
	"fmt"
	"yourapp/feature/model"
)

const (
	AnonymousUser = "<anonymous>"
)

type Logger struct {
	delegate LogDelegate
	pool     *sql.DB
	UserRepo model.UsersRepo
}

func NewLogger(pool *sql.DB, delegate LogDelegate) *Logger {
	return &Logger{delegate: delegate, pool: pool}
}

func (l *Logger) Post(ctx context.Context, username, action string) {
	l.delegate.Debug(action)
	_, err := l.UserRepo.InsertAuditLog(ctx, l.pool, username, action)
	if err != nil {
		l.delegate.Error("Failed to insert into audit log: %w", err)
	}
}

func (l *Logger) Postf(ctx context.Context, username, msg string, args ...any) {
	l.Post(ctx, username, fmt.Sprintf(msg, args...))
}
