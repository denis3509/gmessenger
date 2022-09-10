package auth

import (
	"errors"
	"fmt"
	"messenger/internal/entity"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	UserId    int
	ttl       time.Duration
	createdAt time.Time
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.createdAt.Add(s.ttl))
}

const (
	cost = 10
)

var (
	sessionTtl int64  = 8
	ErrUnauthorized     = errors.New("auth: unauthorized")
	ErrAlreadyExists     = errors.New("auth: ")
	ErrSessionExpired   = errors.New("auth: session expired")
	ErrWrongCredentials = errors.New("auth: wrong username or password")
)

type Service interface {
	Login(username, password string) (string, error)
	Register(username, password, email string) (int, error)
	// returns authenticated user id by session key
	GetUserId(sessionKey string) (int, error)
}

type service struct {
	// TODO mutex
	sessions   map[string]*Session
	repository Repository
}

func NewService(repository Repository) Service {
	return service{
		sessions:   make(map[string]*Session, 100),
		repository: repository,
	}
}

func (s service) NewSession(id int) string {
	session := Session{
		id,
		time.Duration( time.Duration(sessionTtl)  * time.Hour),
		time.Now(),
	}
	uuid := uuid.New().String()
	s.sessions[string(uuid)] = &session
	return uuid
}

func (s service) Login(username string, password string) (string, error) {
	user, err := s.repository.GetByName(username)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return "", err
	}
	sessionKey := s.NewSession(user.Id)
	return sessionKey, nil

}

func (s service) Register(username, password, email string) (int, error) {
	user, _ := s.repository.GetByName(username)
	if user.Id > 0 {
		msg := fmt.Sprintf("auth: user with username '%s' already exists", username)
		return 0, errors.New(msg)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return 0, err
	}
	user = entity.User{
		Username:       username,
		HashedPassword: string(hashedPassword),
		Email:          email,
	}
	id, err := s.repository.Create(user)
	return id, err
}

func (s service) removeSession(sessionKey string) error {
	 delete(s.sessions,sessionKey)
	 return nil
}

func (s service) GetUserId(sessionKey string) (int, error) {
	session, ok := s.sessions[sessionKey]
	if !ok {
		return 0, ErrUnauthorized
	}
	if session.IsExpired() {
		s.removeSession(sessionKey)
		return 0, ErrSessionExpired
	}
	
	return session.UserId, nil
}

 
