package middleware

import (
	"github.com/labstack/echo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.MustLoadMessageFile("en.json")
}

// Middleware untuk setting bahasa dan localizer
func I18nMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		lang := c.Request().Header.Get("Accept-Language")
		if lang == "" {
			lang = "en"
		}
		bundle = i18n.NewBundle(language.English)
		bundle.MustLoadMessageFile("en.json")
		localizer := i18n.NewLocalizer(bundle, lang)
		c.Set("localizer", localizer)
		return next(c)
	}
}

func GetErrorMessage(c echo.Context, code string) (string, error) {
	localizer := c.Get("localizer").(*i18n.Localizer)
	translatedMessage, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: code,
	})
	return translatedMessage, err
}
