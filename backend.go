package proxy

import (
	"crypto/tls"

	"github.com/emersion/go-imap/backend"
	"github.com/emersion/go-imap/client"
)

type Security int

const (
	SecurityNone Security = iota
	SecuritySTARTTLS
	SecurityTLS
)

type Backend struct {
	Addr string
	Security Security
	TLSConfig *tls.Config
}

func New(addr string) *Backend {
	return &Backend{
		Addr: addr,
		Security: SecuritySTARTTLS,
	}
}

func (be *Backend) login(username, password string) (*client.Client, error) {
	var c *client.Client
	var err error
	if be.Security == SecurityTLS {
		if c, err = client.DialTLS(be.Addr, be.TLSConfig); err != nil {
			return nil, err
		}
	} else {
		if c, err = client.Dial(be.Addr); err != nil {
			return nil, err
		}

		if be.Security == SecuritySTARTTLS {
			if err := c.StartTLS(be.TLSConfig); err != nil {
				return nil, err
			}
		}
	}

	if err := c.Login(username, password); err != nil {
		return nil, err
	}

	return c, nil
}

func (be *Backend) Login(username, password string) (backend.User, error) {
	c, err := be.login(username, password)
	if err != nil {
		return nil, err
	}

	u := &user{
		be: be,
		c: c,
		username: username,
	}
	return u, nil
}
