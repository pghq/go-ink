package ink

import (
	"context"
	"time"

	"github.com/pghq/go-ark"
	"github.com/pghq/go-tea"

	"github.com/pghq/go-ink/lang"
)

const (
	// TranslateTTL is the cache ttl for translations
	TranslateTTL = 24 * time.Hour
)

// Translate text
func (l *Linguist) Translate(ctx context.Context, text string, targetLanguage lang.Language, opts ...lang.TranslateOption) (*lang.Text, error) {
	var response *lang.Text
	return response, l.db.Do(ctx, func(tx ark.Txn) error {
		var translation lang.Text
		err := tx.Get("", text, &translation)
		if err == nil {
			response = &translation
			return nil
		}

		response, err = l.translator.Translate(ctx, text, targetLanguage, opts...)
		if err != nil {
			return tea.Stack(err)
		}

		return tx.InsertTTL("", text, response, TranslateTTL)
	})
}
