package ink

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/pghq/go-ark"
	"github.com/pghq/go-tea"
	"github.com/stretchr/testify/assert"

	"github.com/pghq/go-ink/deepl"
	"github.com/pghq/go-ink/lang"
)

func TestMain(m *testing.M) {
	tea.Testing()
	os.Exit(m.Run())
}

func TestLinguist_Translate(t *testing.T) {
	s := serve("testdata/hello-world-de.json")
	dc := deepl.NewClient("", deepl.BaseURL(s.URL))
	l := New("", Translator(dc), Database(ark.New("memory://")))

	t.Run("bad translation", func(t *testing.T) {
		s := serve("does-not-exit")
		dc := deepl.NewClient("", deepl.BaseURL(s.URL))
		l := New("", Translator(dc))

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
