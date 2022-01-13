package mux

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	sessionIdLength = 16
	sessionDuration = time.Minute * 60
)

type Session struct {
	Username string
	timer    *time.Timer
}

type sessions struct {
	list map[string]*Session
	lock sync.Mutex
}

func (s *sessions) add(session *Session) (id string) {
	id = newSessionId()
	s.lock.Lock()
	defer s.lock.Unlock()
	s.list[id] = session
	s.setExpiry(id)
	return
}

func (s *sessions) get(id string) (session *Session, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	var ok bool
	if session, ok = s.list[id]; ok {
		s.setExpiry(id)
		return
	}
	err = fmt.Errorf("no such session")
	return
}

func (s *sessions) remove(id string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if session, ok := s.list[id]; ok {
		if session.timer != nil {
			session.timer.Stop()
		}
	}
	delete(s.list, id)
}

func (s *sessions) setExpiry(id string) {
	session, ok := s.list[id]
	if !ok {
		return
	}
	if session.timer != nil {
		session.timer.Stop()
	}
	session.timer = time.AfterFunc(sessionDuration, func() {
		s.lock.Lock()
		defer s.lock.Unlock()
		if _, ok := s.list[id]; ok {
			delete(s.list, id)
		}
	})
}

func newSessionId() (id string) {
	bytes := make([]byte, sessionIdLength)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Printf("error creating sessionId: %s", err.Error())
		return
	}
	id = base64.RawURLEncoding.EncodeToString(bytes)
	return
}
