package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type UsersRepo struct {
	createUser            func(context.Context, *sql.DB, string, string) (sql.Result, error)
	updatePassword        func(context.Context, *sql.DB, string, string) (sql.Result, error)
	checkPassword         func(context.Context, *sql.DB, string, string) (*CheckPasswordResult, error)
	getUser               func(context.Context, *sql.DB, string) (*GetUserResult, error)
	getAllUsers           func(context.Context, *sql.DB) ([]*GetAllUsersResult, error)
	lockUser              func(context.Context, *sql.DB, string) (sql.Result, error)
	deleteUser            func(context.Context, *sql.DB, string) (sql.Result, error)
	elevateToAdmin        func(context.Context, *sql.DB, string) (sql.Result, error)
	insertAuditLog        func(context.Context, *sql.DB, string, string) (sql.Result, error)
	createSession         func(context.Context, *sql.DB, string) (*CreateSessionResult, error)
	updateSessionLiveness func(context.Context, *sql.DB, string) (sql.Result, error)
	getSessionUser        func(context.Context, *sql.DB, string) (*GetSessionUserResult, error)
	invalidateSession     func(context.Context, *sql.DB, string) (sql.Result, error)
	getLatestLogEntries   func(context.Context, *sql.DB, int) ([]*GetLatestLogEntriesResult, error)
	grantAuth             func(context.Context, *sql.DB, string, string) (sql.Result, error)
	revokeAuth            func(context.Context, *sql.DB, string, string) (sql.Result, error)
	getAuthorizations     func(context.Context, *sql.DB) ([]*GetAuthorizationsResult, error)
	userAuth              func(context.Context, *sql.DB, uint64) ([]*UserAuthResult, error)
	userAuthNotGranted    func(context.Context, *sql.DB, string) ([]*UserAuthNotGrantedResult, error)
}

func (repo *UsersRepo) RedirectCreateUser(delegate func(context.Context, *sql.DB, string, string) (sql.Result, error)) {
	repo.createUser = delegate
}

func (repo *UsersRepo) CreateUser(ctx context.Context, conn *sql.DB, username string, password string) (sql.Result, error) {
	if repo.createUser != nil {
		return repo.createUser(ctx, conn, username, password)
	}
	return CreateUser(ctx, conn, username, password)
}

func (repo *UsersRepo) RedirectUpdatePassword(delegate func(context.Context, *sql.DB, string, string) (sql.Result, error)) {
	repo.updatePassword = delegate
}

func (repo *UsersRepo) UpdatePassword(ctx context.Context, conn *sql.DB, username string, password string) (sql.Result, error) {
	if repo.updatePassword != nil {
		return repo.updatePassword(ctx, conn, username, password)
	}
	return UpdatePassword(ctx, conn, username, password)
}

func (repo *UsersRepo) RedirectCheckPassword(delegate func(context.Context, *sql.DB, string, string) (*CheckPasswordResult, error)) {
	repo.checkPassword = delegate
}

func (repo *UsersRepo) CheckPassword(ctx context.Context, conn *sql.DB, username string, password string) (*CheckPasswordResult, error) {
	if repo.checkPassword != nil {
		return repo.checkPassword(ctx, conn, username, password)
	}
	return CheckPassword(ctx, conn, username, password)
}

func (repo *UsersRepo) RedirectGetUser(delegate func(context.Context, *sql.DB, string) (*GetUserResult, error)) {
	repo.getUser = delegate
}

func (repo *UsersRepo) GetUser(ctx context.Context, conn *sql.DB, username string) (*GetUserResult, error) {
	if repo.getUser != nil {
		return repo.getUser(ctx, conn, username)
	}
	return GetUser(ctx, conn, username)
}

func (repo *UsersRepo) RedirectGetAllUsers(delegate func(context.Context, *sql.DB) ([]*GetAllUsersResult, error)) {
	repo.getAllUsers = delegate
}

func (repo *UsersRepo) GetAllUsers(ctx context.Context, conn *sql.DB) ([]*GetAllUsersResult, error) {
	if repo.getAllUsers != nil {
		return repo.getAllUsers(ctx, conn)
	}
	return GetAllUsers(ctx, conn)
}

func (repo *UsersRepo) RedirectLockUser(delegate func(context.Context, *sql.DB, string) (sql.Result, error)) {
	repo.lockUser = delegate
}

func (repo *UsersRepo) LockUser(ctx context.Context, conn *sql.DB, username string) (sql.Result, error) {
	if repo.lockUser != nil {
		return repo.lockUser(ctx, conn, username)
	}
	return LockUser(ctx, conn, username)
}

func (repo *UsersRepo) RedirectDeleteUser(delegate func(context.Context, *sql.DB, string) (sql.Result, error)) {
	repo.deleteUser = delegate
}

func (repo *UsersRepo) DeleteUser(ctx context.Context, conn *sql.DB, username string) (sql.Result, error) {
	if repo.deleteUser != nil {
		return repo.deleteUser(ctx, conn, username)
	}
	return DeleteUser(ctx, conn, username)
}

func (repo *UsersRepo) RedirectElevateToAdmin(delegate func(context.Context, *sql.DB, string) (sql.Result, error)) {
	repo.elevateToAdmin = delegate
}

func (repo *UsersRepo) ElevateToAdmin(ctx context.Context, conn *sql.DB, username string) (sql.Result, error) {
	if repo.elevateToAdmin != nil {
		return repo.elevateToAdmin(ctx, conn, username)
	}
	return ElevateToAdmin(ctx, conn, username)
}

func (repo *UsersRepo) RedirectInsertAuditLog(delegate func(context.Context, *sql.DB, string, string) (sql.Result, error)) {
	repo.insertAuditLog = delegate
}

func (repo *UsersRepo) InsertAuditLog(ctx context.Context, conn *sql.DB, username string, message string) (sql.Result, error) {
	if repo.insertAuditLog != nil {
		return repo.insertAuditLog(ctx, conn, username, message)
	}
	return InsertAuditLog(ctx, conn, username, message)
}

func (repo *UsersRepo) RedirectCreateSession(delegate func(context.Context, *sql.DB, string) (*CreateSessionResult, error)) {
	repo.createSession = delegate
}

func (repo *UsersRepo) CreateSession(ctx context.Context, conn *sql.DB, username string) (*CreateSessionResult, error) {
	if repo.createSession != nil {
		return repo.createSession(ctx, conn, username)
	}
	return CreateSession(ctx, conn, username)
}

func (repo *UsersRepo) RedirectUpdateSessionLiveness(delegate func(context.Context, *sql.DB, string) (sql.Result, error)) {
	repo.updateSessionLiveness = delegate
}

func (repo *UsersRepo) UpdateSessionLiveness(ctx context.Context, conn *sql.DB, sessionKey string) (sql.Result, error) {
	if repo.updateSessionLiveness != nil {
		return repo.updateSessionLiveness(ctx, conn, sessionKey)
	}
	return UpdateSessionLiveness(ctx, conn, sessionKey)
}

func (repo *UsersRepo) RedirectGetSessionUser(delegate func(context.Context, *sql.DB, string) (*GetSessionUserResult, error)) {
	repo.getSessionUser = delegate
}

func (repo *UsersRepo) GetSessionUser(ctx context.Context, conn *sql.DB, sessionKey string) (*GetSessionUserResult, error) {
	if repo.getSessionUser != nil {
		return repo.getSessionUser(ctx, conn, sessionKey)
	}
	return GetSessionUser(ctx, conn, sessionKey)
}

func (repo *UsersRepo) RedirectInvalidateSession(delegate func(context.Context, *sql.DB, string) (sql.Result, error)) {
	repo.invalidateSession = delegate
}

func (repo *UsersRepo) InvalidateSession(ctx context.Context, conn *sql.DB, sessionKey string) (sql.Result, error) {
	if repo.invalidateSession != nil {
		return repo.invalidateSession(ctx, conn, sessionKey)
	}
	return InvalidateSession(ctx, conn, sessionKey)
}

func (repo *UsersRepo) RedirectGetLatestLogEntries(delegate func(context.Context, *sql.DB, int) ([]*GetLatestLogEntriesResult, error)) {
	repo.getLatestLogEntries = delegate
}

func (repo *UsersRepo) GetLatestLogEntries(ctx context.Context, conn *sql.DB, limit int) ([]*GetLatestLogEntriesResult, error) {
	if repo.getLatestLogEntries != nil {
		return repo.getLatestLogEntries(ctx, conn, limit)
	}
	return GetLatestLogEntries(ctx, conn, limit)
}

func (repo *UsersRepo) RedirectGrantAuth(delegate func(context.Context, *sql.DB, string, string) (sql.Result, error)) {
	repo.grantAuth = delegate
}

func (repo *UsersRepo) GrantAuth(ctx context.Context, conn *sql.DB, userID string, authID string) (sql.Result, error) {
	if repo.grantAuth != nil {
		return repo.grantAuth(ctx, conn, userID, authID)
	}
	return GrantAuth(ctx, conn, userID, authID)
}

func (repo *UsersRepo) RedirectRevokeAuth(delegate func(context.Context, *sql.DB, string, string) (sql.Result, error)) {
	repo.revokeAuth = delegate
}

func (repo *UsersRepo) RevokeAuth(ctx context.Context, conn *sql.DB, userID string, authID string) (sql.Result, error) {
	if repo.revokeAuth != nil {
		return repo.revokeAuth(ctx, conn, userID, authID)
	}
	return RevokeAuth(ctx, conn, userID, authID)
}

func (repo *UsersRepo) RedirectGetAuthorizations(delegate func(context.Context, *sql.DB) ([]*GetAuthorizationsResult, error)) {
	repo.getAuthorizations = delegate
}

func (repo *UsersRepo) GetAuthorizations(ctx context.Context, conn *sql.DB) ([]*GetAuthorizationsResult, error) {
	if repo.getAuthorizations != nil {
		return repo.getAuthorizations(ctx, conn)
	}
	return GetAuthorizations(ctx, conn)
}

func (repo *UsersRepo) RedirectUserAuth(delegate func(context.Context, *sql.DB, uint64) ([]*UserAuthResult, error)) {
	repo.userAuth = delegate
}

func (repo *UsersRepo) UserAuth(ctx context.Context, conn *sql.DB, userID uint64) ([]*UserAuthResult, error) {
	if repo.userAuth != nil {
		return repo.userAuth(ctx, conn, userID)
	}
	return UserAuth(ctx, conn, userID)
}

func (repo *UsersRepo) RedirectUserAuthNotGranted(delegate func(context.Context, *sql.DB, string) ([]*UserAuthNotGrantedResult, error)) {
	repo.userAuthNotGranted = delegate
}

func (repo *UsersRepo) UserAuthNotGranted(ctx context.Context, conn *sql.DB, username string) ([]*UserAuthNotGrantedResult, error) {
	if repo.userAuthNotGranted != nil {
		return repo.userAuthNotGranted(ctx, conn, username)
	}
	return UserAuthNotGranted(ctx, conn, username)
}

func CreateUser(ctx context.Context, conn *sql.DB, username string, password string) (sql.Result, error) {
	const query = `
insert into users (username, pass_hash) values ($1, gen_passwd($2));
`
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction in CreateUser: %w", err)
	}

	result, err := tx.Exec(query, username, password)
	if err != nil {
		rerr := fmt.Errorf("failed to run CreateUser: %w", err)
		return nil, errors.Join(rerr, tx.Rollback())
	}
	return result, tx.Commit()
}

func UpdatePassword(ctx context.Context, conn *sql.DB, username string, password string) (sql.Result, error) {
	const query = `
update users set pass_hash = gen_passwd($2) where username = $1;
`
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction in UpdatePassword: %w", err)
	}

	result, err := tx.Exec(query, username, password)
	if err != nil {
		rerr := fmt.Errorf("failed to run UpdatePassword: %w", err)
		return nil, errors.Join(rerr, tx.Rollback())
	}
	return result, tx.Commit()
}

type CheckPasswordResult struct {
	Matches bool `json:"matches"`
}

func CheckPassword(ctx context.Context, conn *sql.DB, username string, password string) (*CheckPasswordResult, error) {
	const query = `
select check_passwd($1, $2);
`
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction in CheckPassword: %w", err)
	}

	var result CheckPasswordResult
	err = tx.QueryRow(query, username, password).Scan(&result.Matches)
	if err != nil {
		rerr := fmt.Errorf("failed to run CheckPassword: %w", err)
		return nil, errors.Join(rerr, tx.Rollback())
	}
	return &result, tx.Commit()
}

type GetUserResult struct {
	UserID   uint64 `json:"userID"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
}

func GetUser(ctx context.Context, conn *sql.DB, username string) (*GetUserResult, error) {
	const query = `
select id, username, admin from users where username = $1;
`
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction in GetUser: %w", err)
	}

	var result GetUserResult
	err = tx.QueryRow(query, username).Scan(&result.UserID, &result.Username, &result.Admin)
	if err != nil {
		rerr := fmt.Errorf("failed to run GetUser: %w", err)
		return nil, errors.Join(rerr, tx.Rollback())
	}
	return &result, tx.Commit()
}

type GetAllUsersResult struct {
	UserID   uint64 `json:"userID"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
}

func GetAllUsers(ctx context.Context, conn *sql.DB) ([]*GetAllUsersResult, error) {
	const query = `
select id, username, admin from users;
`
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction in GetAllUsers: %w", err)
	}

	var results []*GetAllUsersResult
	rows, err := tx.Query(query)
	if err != nil {
		rerr := fmt.Errorf("failed to run GetAllUsers: %w", err)
		return nil, errors.Join(rerr, tx.Rollback())
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		result := new(GetAllUsersResult)
		if err := rows.Scan(&result.UserID, &result.Username, &result.Admin); err != nil {
			rerr := fmt.Errorf("failed to scan row in GetAllUsers: %w", err)
			return nil, errors.Join(rerr, tx.Rollback())
		}
		results = append(results, result)
	}
	return results, tx.Commit()
}

func LockUser(ctx context.Context, conn *sql.DB, username string) (sql.Result, error) {
	const query = `
--- @write-tx
update users set pass_hash = 'locked' where username = $1;
`
	result, err := conn.ExecContext(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("failed to run LockUser: %w", err)
	}
	return result, nil
}

func DeleteUser(ctx context.Context, conn *sql.DB, username string) (sql.Result, error) {
	const query = `
delete from users where username = $1;
`
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction in DeleteUser: %w", err)
	}

	result, err := tx.Exec(query, username)
	if err != nil {
		rerr := fmt.Errorf("failed to run DeleteUser: %w", err)
		return nil, errors.Join(rerr, tx.Rollback())
	}
	return result, tx.Commit()
}

func ElevateToAdmin(ctx context.Context, conn *sql.DB, username string) (sql.Result, error) {
	const query = `
update users set admin = true where username = $1;
`
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction in ElevateToAdmin: %w", err)
	}

	result, err := tx.Exec(query, username)
	if err != nil {
		rerr := fmt.Errorf("failed to run ElevateToAdmin: %w", err)
		return nil, errors.Join(rerr, tx.Rollback())
	}
	return result, tx.Commit()
}

func InsertAuditLog(ctx context.Context, conn *sql.DB, username string, message string) (sql.Result, error) {
	const query = `
insert into user_audit (username, action) values ($1, $2);
`
	result, err := conn.ExecContext(ctx, query, username, message)
	if err != nil {
		return nil, fmt.Errorf("failed to run InsertAuditLog: %w", err)
	}
	return result, nil
}

type CreateSessionResult struct {
	SessionKey string `json:"sessionKey"`
}

func CreateSession(ctx context.Context, conn *sql.DB, username string) (*CreateSessionResult, error) {
	const query = `
select create_session($1);
`
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction in CreateSession: %w", err)
	}

	var result CreateSessionResult
	err = tx.QueryRow(query, username).Scan(&result.SessionKey)
	if err != nil {
		rerr := fmt.Errorf("failed to run CreateSession: %w", err)
		return nil, errors.Join(rerr, tx.Rollback())
	}
	return &result, tx.Commit()
}

func UpdateSessionLiveness(ctx context.Context, conn *sql.DB, sessionKey string) (sql.Result, error) {
	const query = `
call update_session_ttl($1);
`
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction in UpdateSessionLiveness: %w", err)
	}

	result, err := tx.Exec(query, sessionKey)
	if err != nil {
		rerr := fmt.Errorf("failed to run UpdateSessionLiveness: %w", err)
		return nil, errors.Join(rerr, tx.Rollback())
	}
	return result, tx.Commit()
}

type GetSessionUserResult struct {
	UserID   uint64 `json:"userID"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
}

func GetSessionUser(ctx context.Context, conn *sql.DB, sessionKey string) (*GetSessionUserResult, error) {
	const query = `
select u.id, u.username, u.admin
from session s
    join users u on s.user_id = u.id
where s.session_key = $1
    and s.revoked_at > current_timestamp;
`
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction in GetSessionUser: %w", err)
	}

	var result GetSessionUserResult
	err = tx.QueryRow(query, sessionKey).Scan(&result.UserID, &result.Username, &result.Admin)
	if err != nil {
		rerr := fmt.Errorf("failed to run GetSessionUser: %w", err)
		return nil, errors.Join(rerr, tx.Rollback())
	}
	return &result, tx.Commit()
}

func InvalidateSession(ctx context.Context, conn *sql.DB, sessionKey string) (sql.Result, error) {
	const query = `
delete from session where session_key = $1;
`
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction in InvalidateSession: %w", err)
	}

	result, err := tx.Exec(query, sessionKey)
	if err != nil {
		rerr := fmt.Errorf("failed to run InvalidateSession: %w", err)
		return nil, errors.Join(rerr, tx.Rollback())
	}
	return result, tx.Commit()
}

type GetLatestLogEntriesResult struct {
	Username  string    `json:"username"`
	EventTime time.Time `json:"eventTime"`
	Action    string    `json:"action"`
}

func GetLatestLogEntries(ctx context.Context, conn *sql.DB, limit int) ([]*GetLatestLogEntriesResult, error) {
	const query = `
select
    user_id,
    username,
    event_time,
    action
from (
    select
        u.id "user_id",
        u.username,
        event_time,
        action
    from user_audit ua
    left join users u on ua.username = u.username
    order by event_time desc
    limit $1
) segment
order by event_time;
`
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction in GetLatestLogEntries: %w", err)
	}

	var results []*GetLatestLogEntriesResult
	rows, err := tx.Query(query, limit)
	if err != nil {
		rerr := fmt.Errorf("failed to run GetLatestLogEntries: %w", err)
		return nil, errors.Join(rerr, tx.Rollback())
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		result := new(GetLatestLogEntriesResult)
		if err := rows.Scan(&result.Username, &result.EventTime, &result.Action); err != nil {
			rerr := fmt.Errorf("failed to scan row in GetLatestLogEntries: %w", err)
			return nil, errors.Join(rerr, tx.Rollback())
		}
		results = append(results, result)
	}
	return results, tx.Commit()
}

func GrantAuth(ctx context.Context, conn *sql.DB, userID string, authID string) (sql.Result, error) {
	const query = `
insert into user_authz (user_id, auth_id) values ($1, $2)
on conflict (user_id, auth_id) do update set revoked = null, granted = current_timestamp
;
`
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction in GrantAuth: %w", err)
	}

	result, err := tx.Exec(query, userID, authID)
	if err != nil {
		rerr := fmt.Errorf("failed to run GrantAuth: %w", err)
		return nil, errors.Join(rerr, tx.Rollback())
	}
	return result, tx.Commit()
}

func RevokeAuth(ctx context.Context, conn *sql.DB, userID string, authID string) (sql.Result, error) {
	const query = `
update user_authz set revoked = current_timestamp
where
    user_id = $1
    and auth_id = $2
;
`
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction in RevokeAuth: %w", err)
	}

	result, err := tx.Exec(query, userID, authID)
	if err != nil {
		rerr := fmt.Errorf("failed to run RevokeAuth: %w", err)
		return nil, errors.Join(rerr, tx.Rollback())
	}
	return result, tx.Commit()
}

type GetAuthorizationsResult struct {
	AuthID uint64 `json:"authID"`
	Name   string `json:"name"`
}

func GetAuthorizations(ctx context.Context, conn *sql.DB) ([]*GetAuthorizationsResult, error) {
	const query = `
select id, auth from authorizations;
`
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction in GetAuthorizations: %w", err)
	}

	var results []*GetAuthorizationsResult
	rows, err := tx.Query(query)
	if err != nil {
		rerr := fmt.Errorf("failed to run GetAuthorizations: %w", err)
		return nil, errors.Join(rerr, tx.Rollback())
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		result := new(GetAuthorizationsResult)
		if err := rows.Scan(&result.AuthID, &result.Name); err != nil {
			rerr := fmt.Errorf("failed to scan row in GetAuthorizations: %w", err)
			return nil, errors.Join(rerr, tx.Rollback())
		}
		results = append(results, result)
	}
	return results, tx.Commit()
}

type UserAuthResult struct {
	Id      uint64    `json:"id"`
	Auth    string    `json:"auth"`
	Granted time.Time `json:"granted"`
}

func UserAuth(ctx context.Context, conn *sql.DB, userID uint64) ([]*UserAuthResult, error) {
	const query = `
select auth_id, auth, granted from auth_grants where user_id = $1;
`
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction in UserAuth: %w", err)
	}

	var results []*UserAuthResult
	rows, err := tx.Query(query, userID)
	if err != nil {
		rerr := fmt.Errorf("failed to run UserAuth: %w", err)
		return nil, errors.Join(rerr, tx.Rollback())
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		result := new(UserAuthResult)
		if err := rows.Scan(&result.Id, &result.Auth, &result.Granted); err != nil {
			rerr := fmt.Errorf("failed to scan row in UserAuth: %w", err)
			return nil, errors.Join(rerr, tx.Rollback())
		}
		results = append(results, result)
	}
	return results, tx.Commit()
}

type UserAuthNotGrantedResult struct {
	Id   uint64 `json:"id"`
	Auth string `json:"auth"`
}

func UserAuthNotGranted(ctx context.Context, conn *sql.DB, username string) ([]*UserAuthNotGrantedResult, error) {
	const query = `
select id, auth
from authorizations
where id not in (select auth_id from auth_grants where username = $1)
;
`
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction in UserAuthNotGranted: %w", err)
	}

	var results []*UserAuthNotGrantedResult
	rows, err := tx.Query(query, username)
	if err != nil {
		rerr := fmt.Errorf("failed to run UserAuthNotGranted: %w", err)
		return nil, errors.Join(rerr, tx.Rollback())
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		result := new(UserAuthNotGrantedResult)
		if err := rows.Scan(&result.Id, &result.Auth); err != nil {
			rerr := fmt.Errorf("failed to scan row in UserAuthNotGranted: %w", err)
			return nil, errors.Join(rerr, tx.Rollback())
		}
		results = append(results, result)
	}
	return results, tx.Commit()
}
