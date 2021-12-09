package ink

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pghq/go-ark"
	"github.com/stretchr/testify/assert"

	"github.com/pghq/go-ink/deepl"
	"github.com/pghq/go-ink/lang"
)

func TestLinguist_Translate(t *testing.T) {
	s := serve("testdata/hello-world-de.json")
	dc := deepl.NewClient("", deepl.BaseURL(s.URL))
	kvs, _ := ark.Open().ConnectKVS(context.TODO(), "inmem")

	l := NewLinguist("", Translator(dc), KVSConn(kvs))

	t.Run("bad translation", func(t *testing.T) {
		s := serve("does-not-exit")
		dc := deepl.NewClient("", deepl.BaseURL(s.URL))
		l := NewLinguist("", Translator(dc))

		_, err := l.Translate(context.TODO(), "Hello, world!", lang.German)
		assert.NotNil(t, err)
	})

	t.Run("uncached response", func(t *testing.T) {
		text, err := l.Translate(context.TODO(), "Hello, world!", lang.German)
		assert.Nil(t, err)
		assert.NotNil(t, text)
	})

	t.Run("cached response", func(t *testing.T) {
		text, err := l.Translate(context.TODO(), "Hello, world!", lang.German)
		assert.Nil(t, err)
		assert.NotNil(t, text)
	})
}

func serve(path string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}))
}
