package i18nx

import (
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func initBundleMessages(t *testing.T, bundle *Bundle) {
	err := bundle.Instance().AddMessages(language.English, &i18n.Message{
		ID:    "cat",
		Other: "{{.Name}} has {{.Count}} cat.",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = bundle.Instance().AddMessages(language.Chinese, &i18n.Message{
		ID:    "cat",
		Other: "{{.Name}} 有 {{.Count}} 只猫。",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestLocalize(t *testing.T) {
	bundle := NewBundle()
	initBundleMessages(t, bundle)

	msg, err := bundle.Localize(language.English.String(), "cat", WithData(map[string]interface{}{
		"Name":  "Nick",
		"Count": 1,
	}))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Nick has 1 cat.", msg)

	msg, err = bundle.Localize(language.Chinese.String(), "cat", WithData(map[string]interface{}{
		"Name":  "尼克",
		"Count": 1,
	}))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "尼克 有 1 只猫。", msg)
}
