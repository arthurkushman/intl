package intl

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

type Delivery struct {
	Date  string
	Price int64
}

func TestIntl_Translate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var tests = map[string]struct {
		before   func(i *Intl)
		key      string
		lang     string
		delivery *Delivery
		msg      string
		err      error
	}{
		"ok": {
			before: func(i *Intl) {
				mock.ExpectQuery("SELECT .*").
					WithArgs("delivery.datetime.price", "fr-FR").
					WillReturnRows(sqlmock.NewRows([]string{"message", "translation", "localize_config", "localize_config"}).
						AddRow("Hi, your delivery date is {{.Date}} and a price is {{.Price}}",
							"Bonjour, votre date de livraison est le {{.Date}} et le prix est le {{.Price}}", "", ""))
			},
			key:      "delivery.datetime.price",
			lang:     "fr-FR",
			delivery: &Delivery{Date: "demain", Price: 123},
			msg:      "Bonjour, votre date de livraison est le demain et le prix est le 123",
			err:      nil,
		},
		"ok empty translation": {
			before: func(i *Intl) {
				mock.ExpectQuery("SELECT .*").
					WithArgs("delivery.datetime.price", "fr-FR").
					WillReturnRows(sqlmock.NewRows([]string{"message", "translation", "localize_config", "localize_config"}).
						AddRow("Hi, your delivery date is {{.Date}} and a price is {{.Price}}",
							"", "", ""))
			},
			key:      "delivery.datetime.price",
			lang:     "fr-FR",
			delivery: &Delivery{Date: "tomorrow", Price: 123},
			msg:      "Hi, your delivery date is tomorrow and a price is 123",
			err:      nil,
		},
		"ok with changed table names": {
			before: func(i *Intl) {
				i.SetTblNames("tbl1", "tbl2")
				mock.ExpectQuery("SELECT .*").
					WithArgs("delivery.datetime.price", "fr-FR").
					WillReturnRows(sqlmock.NewRows([]string{"message", "translation", "localize_config", "localize_config"}).
						AddRow("Hi, your delivery date is {{.Date}} and a price is {{.Price}}",
							"Bonjour, votre date de livraison est le {{.Date}} et le prix est le {{.Price}}", "", ""))
			},
			key:      "delivery.datetime.price",
			lang:     "fr-FR",
			delivery: &Delivery{Date: "demain", Price: 123},
			msg:      "Bonjour, votre date de livraison est le demain et le prix est le 123",
			err:      nil,
		},
		"template: translate:1: unexpected ": {
			before: func(i *Intl) {
				mock.ExpectQuery("SELECT .*").
					WithArgs("delivery.datetime.price", "fr-FR").
					WillReturnRows(sqlmock.NewRows([]string{"message", "translation", "localize_config", "localize_config"}).
						AddRow("Hi, your delivery date is {{.Date}} and a price is {{.Price}}",
							"Bonjour, votre date de livraison est le {{{.Date}} et le prix est le {{.Price}}", "", ""))
			},
			key:      "delivery.datetime.price",
			lang:     "fr-FR",
			delivery: &Delivery{Date: "demain", Price: 123},
			msg:      "Bonjour, votre date de livraison est le demain et le prix est le 123",
			err:      errors.New("template: translate:1: unexpected \"{\" in command"),
		},
		"err sql: no rows in result set": {
			before: func(i *Intl) {
				mock.ExpectQuery("SELECT .*").
					WithArgs("delivery.datetime.price", "fr-FR").
					WillReturnRows(sqlmock.NewRows([]string{"message", "translation", "localize_config", "localize_config"}))
			},
			key:      "delivery.datetime.price",
			lang:     "fr-FR",
			delivery: &Delivery{Date: "завтра же", Price: 123},
			msg:      "",
			err:      errors.New("sql: no rows in result set"),
		},
	}

	intl := NewIntl(db)
	for n, tt := range tests {
		t.Run(n, func(t *testing.T) {
			tt.before(intl)
			msg, err := intl.Translate(tt.key, tt.lang, tt.delivery)
			if tt.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.msg, msg)
			} else {
				assert.Error(t, err, tt.err)
			}
		})
	}
}

func TestIntl_TranslatePlural(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var v interface{}
	jsonErr := json.Unmarshal([]byte(`{{\"DefaultMessage\": {\"ID\": \"Delivery\", \"One\": \"Bonjour, votre date de livraison est le "+
								"{{.Date}} et le prix est le {{.Price}}\", \"Other\": \"Bonjour, votre date de livraison est le "+
								"{{.Date}} many et le prix est le {{.Price}} many\"}, \"TemplateData\": {\"Date\": \"demain\", \"Price\": 123}, "+
								"\"PluralCount\": 2}`), &v)

	var tests = map[string]struct {
		before func(i *Intl)
		key    string
		lang   string
		msg    string
		err    error
	}{
		"ok": {
			before: func(i *Intl) {
				mock.ExpectQuery("SELECT .*").
					WithArgs("delivery.datetime.price", "fr-FR").
					WillReturnRows(sqlmock.NewRows([]string{"message", "translation", "localize_config", "localize_config"}).
						AddRow("Hi, your delivery date is {{.Date}} and a price is {{.Price}}",
							"Bonjour, votre date de livraison est le {{.Date}} et le prix est le {{.Price}}",
							"{\"DefaultMessage\": {\"ID\": \"Delivery\", \"One\": \"Bonjour, votre date de livraison est le "+
								"{{.Date}} et le prix est le {{.Price}}\", \"Other\": \"Bonjour, votre date de livraison est le "+
								"{{.Date}} many et le prix est le {{.Price}} many\"}, \"TemplateData\": {\"Date\": \"demain\", \"Price\": 123}, "+
								"\"PluralCount\": 2}",
							"{\"DefaultMessage\": {\"ID\": \"Delivery\", \"One\": \"Bonjour, votre date de livraison est le "+
								"{{.Date}} et le prix est le {{.Price}}\", \"Other\": \"Bonjour, votre date de livraison est le "+
								"{{.Date}} many et le prix est le {{.Price}} many\"}, \"TemplateData\": {\"Date\": \"demain\", \"Price\": 123}, "+
								"\"PluralCount\": 2}"))
			},
			key:  "delivery.datetime.price",
			lang: "fr-FR",
			msg:  "Bonjour, votre date de livraison est le demain many et le prix est le 123 many",
			err:  nil,
		},
		"err": {
			before: func(i *Intl) {
				mock.ExpectQuery("SELECT .*").
					WithArgs("delivery.datetime.price", "fr-FR").
					WillReturnRows(sqlmock.NewRows([]string{"message", "translation", "localize_config", "localize_config"}).
						AddRow("Hi, your delivery date is {{.Date}} and a price is {{.Price}}",
							"Bonjour, votre date de livraison est le {{.Date}} et le prix est le {{.Price}}",
							"{\"DefaultMessage\": {\"ID\": \"Delivery\", \"One\": \"Bonjour, votre date de livraison est le "+
								"{{.Date}} et le prix est le {{.Price}}\", \"Other\": \"Bonjour, votre date de livraison est le "+
								"{{.Date}} many et le prix est le {{.Price}} many\"}, \"TemplateData\": {\"Date\": \"demain\", \"Price\": 123}, "+
								"\"PluralCount\": 2}",
							"{{\"DefaultMessage\": {\"ID\": \"Delivery\", \"One\": \"Bonjour, votre date de livraison est le "+
								"{{.Date}} et le prix est le {{.Price}}\", \"Other\": \"Bonjour, votre date de livraison est le "+
								"{{.Date}} many et le prix est le {{.Price}} many\"}, \"TemplateData\": {\"Date\": \"demain\", \"Price\": 123}, "+
								"\"PluralCount\": 2}"))
			},
			key:  "delivery.datetime.price",
			lang: "fr-FR",
			msg:  "Bonjour, votre date de livraison est le demain many et le prix est le 123 many",
			err:  jsonErr,
		},
		"ok empty 2nd localize_config": {
			before: func(i *Intl) {
				mock.ExpectQuery("SELECT .*").
					WithArgs("delivery.datetime.price", "fr-FR").
					WillReturnRows(sqlmock.NewRows([]string{"message", "translation", "localize_config", "localize_config"}).
						AddRow("Hi, your delivery date is {{.Date}} and a price is {{.Price}}",
							"Bonjour, votre date de livraison est le {{.Date}} et le prix est le {{.Price}}",
							"{\"DefaultMessage\": {\"ID\": \"Delivery\", \"One\": \"Bonjour, votre date de livraison est le "+
								"{{.Date}} et le prix est le {{.Price}}\", \"Other\": \"Bonjour, votre date de livraison est le "+
								"{{.Date}} many et le prix est le {{.Price}} many\"}, \"TemplateData\": {\"Date\": \"demain\", \"Price\": 123}, "+
								"\"PluralCount\": 2}", ""))
			},
			key:  "delivery.datetime.price",
			lang: "fr-FR",
			msg:  "Bonjour, votre date de livraison est le demain many et le prix est le 123 many",
			err:  nil,
		},
		"err sql: no rows in result set": {
			before: func(i *Intl) {
				mock.ExpectQuery("SELECT .*").
					WithArgs("delivery.datetime.price", "fr-FR").
					WillReturnRows(sqlmock.NewRows([]string{"message", "translation", "localize_config", "localize_config"}))
			},
			key:  "delivery.datetime.price",
			lang: "fr-FR",
			msg:  "",
			err:  errors.New("sql: no rows in result set"),
		},
		"language: tag is not well-formed": {
			before: func(i *Intl) {
				mock.ExpectQuery("SELECT .*").
					WithArgs("delivery.datetime.price", "fr-FR").
					WillReturnRows(sqlmock.NewRows([]string{"message", "translation", "localize_config", "localize_config"}).
						AddRow("Hi, your delivery date is {{.Date}} and a price is {{.Price}}",
							"Bonjour, votre date de livraison est le {{.Date}} et le prix est le {{.Price}}",
							"{\"DefaultMessage\": {\"ID\": \"Delivery\", \"One\": \"Bonjour, votre date de livraison est le "+
								"{{.Date}} et le prix est le {{.Price}}\", \"Other\": \"Bonjour, votre date de livraison est le "+
								"{{.Date}} many et le prix est le {{.Price}} many\"}, \"TemplateData\": {\"Date\": \"demain\", \"Price\": 123}, "+
								"\"PluralCount\": 2}", ""))
			},
			key:  "delivery.datetime.price",
			lang: "___fake-Lang___",
			msg:  "",
			err:  errors.New("language: tag is not well-formed"),
		},
	}

	intl := NewIntl(db)
	for n, tt := range tests {
		t.Run(n, func(t *testing.T) {
			tt.before(intl)
			msg, err := intl.TranslatePlurals(tt.key, tt.lang)
			if tt.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.msg, msg)
			} else {
				assert.Error(t, err, tt.err)
				assert.Equal(t, tt.err, err)
			}
		})
	}
}
