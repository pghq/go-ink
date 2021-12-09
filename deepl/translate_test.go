package deepl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pghq/go-ink/lang"
)

func TestClient_Translate(t *testing.T) {
	t.Run("nil timeout", func(t *testing.T) {
		c := NewClient("")
		_, err := c.Translate(nil, "Hello, world!", lang.German)
		assert.NotNil(t, err)
	})

	t.Run("context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), 0)
		defer cancel()
		s := serve("does-not-exit")
		c := NewClient("", BaseURL(s.URL))
		_, err := c.Translate(ctx, "Hello, world!", lang.German)
		assert.NotNil(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		c := NewClient("", BaseURL("https://does-not-exist"))
		_, err := c.Translate(context.TODO(), "Hello, world!", lang.German)
		assert.NotNil(t, err)
	})

	t.Run("no translations", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{}`))
		}))

		c := NewClient("", BaseURL(s.URL))
		_, err := c.Translate(context.TODO(), "Hello, world!", lang.German)
		assert.NotNil(t, err)
	})

	t.Run("bad request", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))

		c := NewClient("", BaseURL(s.URL))
		_, err := c.Translate(context.TODO(), "Hello, world!", lang.German)
		assert.NotNil(t, err)
	})

	t.Run("forbidden", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))

		c := NewClient("", BaseURL(s.URL))
		_, err := c.Translate(context.TODO(), "Hello, world!", lang.German)
		assert.NotNil(t, err)
	})

	t.Run("request entity too large", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
		}))

		c := NewClient("", BaseURL(s.URL))
		_, err := c.Translate(context.TODO(), "Hello, world!", lang.German)
		assert.NotNil(t, err)
	})

	t.Run("too many requests", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
		}))

		c := NewClient("", BaseURL(s.URL))
		_, err := c.Translate(context.TODO(), "Hello, world!", lang.German)
		assert.NotNil(t, err)
	})

	t.Run("quota exceeded", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(456)
		}))

		c := NewClient("", BaseURL(s.URL))
		_, err := c.Translate(context.TODO(), "Hello, world!", lang.German)
		assert.NotNil(t, err)
	})

	t.Run("service unavailable", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))

		c := NewClient("", BaseURL(s.URL))
		_, err := c.Translate(context.TODO(), "Hello, world!", lang.German)
		assert.NotNil(t, err)
	})

	t.Run("too many requests internal", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(529)
		}))

		c := NewClient("", BaseURL(s.URL))
		_, err := c.Translate(context.TODO(), "Hello, world!", lang.German)
		assert.NotNil(t, err)
	})

	t.Run("internal error", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		c := NewClient("", BaseURL(s.URL))
		_, err := c.Translate(context.TODO(), "Hello, world!", lang.German)
		assert.NotNil(t, err)
	})

	t.Run("unexpected body", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`[]`))
		}))

		c := NewClient("", BaseURL(s.URL))
		_, err := c.Translate(context.TODO(), "Hello, world!", lang.German)
		assert.NotNil(t, err)
	})

	t.Run("bad request", func(t *testing.T) {
		s := serve("does-not-exit")
		c := NewClient("", BaseURL(s.URL))
		_, err := c.Translate(context.TODO(), "Hello, world!", lang.German)
		assert.NotNil(t, err)
	})

	t.Run("english to german", func(t *testing.T) {
		s := serve("../testdata/hello-world-de.json")
		c := NewClient("", BaseURL(s.URL), HttpClient(http.DefaultClient))

		opts := []lang.TranslateOption{
			SourceLanguage(lang.AmericanEnglish),
			SplitSentences("1"),
			PreserveFormatting(),
			Formal(),
			Informal(),
			GlossaryId("id"),
			TagHandling(),
			NonSplittingTags(),
			SplittingTags(),
			IgnoreTags(),
			DisableOutlineDetection(),
		}

		text, err := c.Translate(context.TODO(), "Hello, world!", lang.German, opts...)
		assert.Nil(t, err)
		assert.Equal(t, lang.English, text.SourceLanguage)
		assert.Equal(t, "Hallo, Welt!", text.Translation)
	})
}

func serve(path string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}))
}
