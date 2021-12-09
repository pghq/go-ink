package ink

import (
	"context"

	"github.com/pghq/go-ark"

	"github.com/pghq/go-ink/deepl"
	"github.com/pghq/go-ink/lang"
)

// Linguist translates text to various languages
type Linguist struct {
	translator lang.Translator
	conn       *ark.KVSConn
}

// NewLinguist creates a new linguist instance
func NewLinguist(authKey string, opts ...LinguistOption) *Linguist {
	conn, _ := ark.Open().ConnectKVS(context.Background(), "inmem")
	l := Linguist{
		translator: deepl.NewClient(authKey),
		conn:       conn,
	}

	for _, opt := range opts {
		opt(&l)
	}

	return &l
}

// LinguistOption to configure custom translator
type LinguistOption func(l *Linguist)

// Translator sets a custom translator
func Translator(o lang.Translator) LinguistOption {
	return func(l *Linguist) {
		l.translator = o
	}
}

// KVSConn sets a custom KVS client for the translator
func KVSConn(o *ark.KVSConn) LinguistOption {
	return func(l *Linguist) {
		l.conn = o
	}
}
