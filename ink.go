package ink

import (
	"github.com/pghq/go-ark"

	"github.com/pghq/go-ink/deepl"
	"github.com/pghq/go-ink/lang"
)

// Linguist translates text to various languages
type Linguist struct {
	translator lang.Translator
	mapper     *ark.Mapper
}

// New creates a new linguist instance
func New(authKey string, opts ...LinguistOption) *Linguist {
	l := Linguist{
		translator: deepl.NewClient(authKey),
		mapper:     ark.New(),
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

// Mapper sets a custom data mapper
func Mapper(o *ark.Mapper) LinguistOption {
	return func(l *Linguist) {
		l.mapper = o
	}
}
