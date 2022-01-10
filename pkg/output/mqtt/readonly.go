package mqtt

import (
	auth "github.com/mochi-co/mqtt/server/listeners/auth"
)

var ReadOnlyAuth ReadOnly

type ReadOnly struct{}

var _ auth.Controller = (*ReadOnly)(nil)

func (ReadOnly) Authenticate(_user, _password []byte) bool {
	return true
}

func (ReadOnly) ACL(_user []byte, _topic string, write bool) bool {
	return !write
}
