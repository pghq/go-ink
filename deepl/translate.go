package deepl

import (
	"context"
	"net/url"
	"strings"

	"github.com/pghq/go-tea"

	"github.com/pghq/go-ink/lang"
)

// TranslationResponse List of the translations in the order the text parameters have been specified.
type TranslationResponse struct {
	Translations []Translation `json:"translations"`
}

// Translation A single instance of translated text
type Translation struct {
	// DetectedSourceLanguage The language detected in the source text. It reflects the value of the source_lang parameter, when specified.
	DetectedSourceLanguage string `json:"detected_source_language"`

	// Text The translated text.
	Text string `json:"text"`
}

// Translate is an API service that allows
// to translate texts and is available at https://api-free.deepl.com/v2/translate.
// Documentation: https://www.deepl.com/docs-api/translating-text/
func (c Client) Translate(ctx context.Context, text string, targetLanguage lang.Language, opts ...lang.TranslateOption) (*lang.Text, error) {
	form := url.Values{}
	for _, opt := range opts {
		opt(form)
	}

	form.Set("target_lang", string(targetLanguage))
	form.Set("text", text)

	req, err := c.newRequest(ctx, "POST", "/translate", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, tea.Error(err)
	}

	var resp TranslationResponse
	if err := c.do(req, &resp); err != nil {
		return nil, tea.Error(err)
	}

	if len(resp.Translations) == 0 {
		return nil, tea.NewError("no translations")
	}

	translation := lang.Text{
		SourceLanguage: lang.Language(resp.Translations[0].DetectedSourceLanguage),
		Translation:    resp.Translations[0].Text,
	}

	return &translation, nil
}

// SourceLanguage Language of the text to be translated
func SourceLanguage(o lang.Language) lang.TranslateOption {
	return func(v url.Values) {
		v.Set("source_language", string(o))
	}
}

// SplitSentences Sets whether the translation engine should first
// split the input into sentences. This is enabled by default.
func SplitSentences(o string) lang.TranslateOption {
	return func(v url.Values) {
		v.Set("split_sentences", o)
	}
}

// PreserveFormatting Sets whether the translation engine should respect the original formatting,
// even if it would usually correct some aspects.
func PreserveFormatting() lang.TranslateOption {
	return func(v url.Values) {
		v.Set("preserve_formatting", "1")
	}
}

// Formal Sets the formality to more formal
// Sets whether the translated text should lean towards formal or informal language.
// This feature currently only works for target languages "DE" (German),
// "FR" (French), "IT" (Italian), "ES" (Spanish), "NL" (Dutch), "PL" (Polish),
// "PT-PT", "PT-BR" (Portuguese) and "RU" (Russian).
func Formal() lang.TranslateOption {
	return func(v url.Values) {
		v.Set("formality", "more")
	}
}

// Informal Sets the formality to less formal
// Sets whether the translated text should lean towards formal or informal language.
// This feature currently only works for target languages "DE" (German),
// "FR" (French), "IT" (Italian), "ES" (Spanish), "NL" (Dutch), "PL" (Polish),
// "PT-PT", "PT-BR" (Portuguese) and "RU" (Russian).
func Informal() lang.TranslateOption {
	return func(v url.Values) {
		v.Set("formality", "less")
	}
}

// GlossaryId Specify the glossary to use for the translation.
// Important: This requires the source_lang parameter to be set and the
// language pair of the glossary has to match the language pair of the request.
func GlossaryId(o string) lang.TranslateOption {
	return func(v url.Values) {
		v.Set("glossary_id", o)
	}
}

// TagHandling Sets which kind of tags should be handled.
func TagHandling() lang.TranslateOption {
	return func(v url.Values) {
		v.Set("tag_handling", "xml")
	}
}

// NonSplittingTags Comma-separated list of XML tags which never split sentences.
func NonSplittingTags(tags ...string) lang.TranslateOption {
	return func(v url.Values) {
		v.Set("non_splitting_tags", strings.Join(tags, ","))
	}
}

//SplittingTags Comma-separated list of XML tags which always cause splits.
func SplittingTags(tags ...string) lang.TranslateOption {
	return func(v url.Values) {
		v.Set("splitting_tags", strings.Join(tags, ","))
	}
}

// IgnoreTags Comma-separated list of XML tags that indicate text not to be translated.
func IgnoreTags(tags ...string) lang.TranslateOption {
	return func(v url.Values) {
		v.Set("ignore_tags", strings.Join(tags, ","))
	}
}

// DisableOutlineDetection The automatic detection of the XML structure won't yield best results in all XML files.
// You can disable this automatic mechanism altogether by setting the outline_detection parameter to 0 and selecting
// the tags that should be considered structure tags. This will split sentences using the splitting_tags parameter.
func DisableOutlineDetection() lang.TranslateOption {
	return func(v url.Values) {
		v.Set("outline_detection", "0")
	}
}
