package session

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type SessionsManager struct {
	data *sql.DB
	mu   *sync.RWMutex
}

func NewSessionsManager(db *sql.DB) *SessionsManager {
	return &SessionsManager{
		data: db,
		mu:   &sync.RWMutex{},
	}
}

func (sm *SessionsManager) Check(r *http.Request) (*Session, error) {
	sessionCookie, err := r.Cookie("session_id")
	fmt.Println("ERRERER ", err)
	if err == http.ErrNoCookie {
		return nil, ErrNoAuth
	}

	s := &Session{}
	sm.mu.RLock()
	row := sm.data.QueryRow("SELECT id, userid FROM sessions WHERE id = ? LIMIT 1", sessionCookie.Value)
	err = row.Scan(&s.ID, &s.UserID)
	sm.mu.RUnlock()

	if err != nil {
		return nil, ErrNoAuth
	}

	return s, nil
}

func (sm *SessionsManager) Create(w http.ResponseWriter, userID uint32, path string) (*Session, error) {
	sess := NewSession(userID)

	sm.mu.Lock()
	_, err := sm.data.Exec(
		"INSERT INTO sessions (`id`, `userid`) VALUES (?, ?)",
		sess.ID,
		sess.UserID,
	)
	if err != nil {
		return nil, fmt.Errorf("no user")
	}

	sm.mu.Unlock()

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	return sess, nil
}

func (sm *SessionsManager) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	sess, err := SessionFromContext(r.Context())
	if err != nil {
		return err
	}

	sm.mu.Lock()
	_, err = sm.data.Exec(
		"DELETE FROM sessions WHERE id = ?",
		sess.ID,
	)
	if err != nil {
		return err
	}

	sm.mu.Unlock()

	cookie := http.Cookie{
		Name:    "session_id",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	return nil
}
