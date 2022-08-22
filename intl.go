package intl

import (
	"bytes"
	"database/sql"
	"html/template"
)

type Translator interface {
	Translate(key, lang string, params any) (string, error)
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

	t, err := template.New("todos").Parse(trans.TranslatedMsg)
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
