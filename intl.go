package intl

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"html/template"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const DefaultTemplateName = "translate"

type Translator interface {
	Translate(key, lang string, params any) (string, error)
	TranslatePlurals(key, lang string) (string, error)
}

type Intl struct {
	db *sql.DB
}

func NewIntl(db *sql.DB) *Intl {
	return &Intl{
		db: db,
	}
}

// Translate translates message by key and language, replacing with params from custom struct
// there should be templated message in db (for particular language) to perform such transformation ex.:
// Hello {{.Name}}, it's {{.Time}}.
// Basket: {{.Product}} price is {{.Price}}.
func (i *Intl) Translate(key, lang string, params any) (string, error) {
	trans, err := i.GetMessage(key, lang)
	if err != nil {
		return "", err
	}

	var t *template.Template
	if trans.TranslatedMsg != "" {
		t, err = template.New(DefaultTemplateName).Parse(trans.TranslatedMsg)
	} else {
		t, err = template.New(DefaultTemplateName).Parse(trans.Message)
	}

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	err = t.Execute(&tpl, params)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil
}

// TranslatePlurals translates any message with plurals based on LocalizeConfig object
// see: github.com/nicksnyder/go-i18n/v2/i18n
// it tries to find translation language and makes fallback if not found
/*
ex. of LocalizeConfig:
localizer.Localize(&i18n.LocalizeConfig{
    DefaultMessage: &i18n.Message{
        ID: "PersonCats",
        One: "{{.Name}} has {{.Count}} cat.",
        Other: "{{.Name}} has {{.Count}} cats.",
    },
    TemplateData: map[string]interface{}{
        "Name": "Nick",
        "Count": 2,
    },
    PluralCount: 2,
}) // Nick has 2 cats.
*/
func (i *Intl) TranslatePlurals(key, lang string) (string, error) {
	t, err := language.Parse(lang)
	if err != nil {
		return "", err
	}

	bundle := i18n.NewBundle(t)
	localizer := i18n.NewLocalizer(bundle, lang)

	trans, err := i.GetMessage(key, lang)
	if err != nil {
		return "", err
	}

	localizeConfig := &i18n.LocalizeConfig{}
	if trans.TranslatedLocalizeConfig != "" {
		err = json.Unmarshal([]byte(trans.TranslatedLocalizeConfig), localizeConfig)
	} else {
		err = json.Unmarshal([]byte(trans.LocalizeConfig), localizeConfig)
	}

	if err != nil {
		return "", err
	}

	localizeConfig.PluralCount = int(localizeConfig.PluralCount.(float64))
	msg, err := localizer.Localize(localizeConfig)
	if err != nil {
		return "", err
	}

	return msg, nil
}
