package i18nx

import "github.com/nicksnyder/go-i18n/v2/i18n"

type OptionFunc func(config *i18n.LocalizeConfig) *i18n.LocalizeConfig

func WithPluralCount(pluralCount interface{}) OptionFunc {
	return func(c *i18n.LocalizeConfig) *i18n.LocalizeConfig {
		c.PluralCount = pluralCount
		return c
	}
}

func WithData(data interface{}) OptionFunc {
	return func(c *i18n.LocalizeConfig) *i18n.LocalizeConfig {
		c.TemplateData = data
		return c
	}
}
