package database

import (
	"fmt"
	"net/url"

	"emperror.dev/errors"

	_ "github.com/lib/pq"
)

type Config struct {
	Port     int
	Host     string
	User     string
	Password string
	Name     string

	Params map[string]string
}

func (c Config) Validate() error {
	if c.Host == "" {
		return errors.New("database host is required ")
	}
	if c.Port == 0 {
		return errors.New("database port is required ")
	}
	if c.User == "" {
		return errors.New("database user is required ")
	}
	if c.Password == "" {
		return errors.New("database password is required ")
	}
	if c.Name == "" {
		return errors.New("database name is required ")
	}

	return nil
}

func (c Config) DSN() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:   c.Name,
	}

	q := url.Values{}
	for k, v := range c.Params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String()
}
