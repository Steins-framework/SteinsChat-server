package pool

import (
	"chat/lib"
	"chat/models"
	"errors"
)

var UserConnPool = UserConn{
	pool: map[string]*lib.Client{},
}

type UserConn struct {
	pool map[string]*lib.Client
}

func (p *UserConn) Get(key string) (*lib.Client, error) {
	if client, exist := p.pool[key]; exist {
		return client, nil
	}

	return nil, errors.New("user is offline")
}

func (p *UserConn) Add(key string, client *lib.Client) {
	if c, err := p.Get(key); err == nil {
		if client.Conn != c.Conn {
			_ = c.Close()
		}
	}

	p.pool[key] = client
}

func (p *UserConn) IsOnline(user *models.User) bool {
	client, err := p.Get(user.Id)

	if err != nil {
		return false
	}

	return client.IsConnect()
}

func (p *UserConn) Delete(key string) {
	delete(p.pool, key)
}

func (p *UserConn) OfUser(user *models.User) (*lib.Client, error) {
	return p.Get(user.Id)
}

func (p *UserConn) Offline(user *models.User) {
	client, _ := p.Get(user.Id)

	_ = client.Close()

	p.Delete(user.Id)
}