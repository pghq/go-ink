package lang

import (
	"context"
	"net/url"
)

// TranslateOption for custom translate requests
type TranslateOption func(v url.Values)

// Translator of text
type Translator interface {
	Translate(ctx context.Context, text string, target Language, opts ...TranslateOption) (*Text, error)
}
