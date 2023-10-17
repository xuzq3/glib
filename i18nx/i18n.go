package i18nx

import (
	"embed"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Bundle struct {
	bundle *i18n.Bundle
}

func NewBundle() *Bundle {
	bundle := i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	return &Bundle{
		bundle: bundle,
	}
}

func (b *Bundle) Instance() *i18n.Bundle {
	return b.bundle
}

func (b *Bundle) LoadEmbedFS(fs embed.FS, dir string) error {
	entries, err := fs.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		entryPath := filepath.Join(dir, entry.Name())
		_, err := b.bundle.LoadMessageFileFS(fs, entryPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bundle) Localize(lang string, id string, opts ...OptionFunc) (string, error) {
	if id == "" {
		return "", nil
	}

	c := &i18n.LocalizeConfig{
		MessageID: id,
	}
	for _, opt := range opts {
		c = opt(c)
	}

	localizer := i18n.NewLocalizer(b.bundle, lang)
	return localizer.Localize(c)
}
