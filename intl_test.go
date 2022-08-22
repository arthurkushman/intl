package intl

import (
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
					WillReturnRows(sqlmock.NewRows([]string{"message", "translation"}).
						AddRow("Hi, your delivery date is {{.Date}} and a price is {{.Price}}",
							"Bonjour, votre date de livraison est le {{.Date}} et le prix est le {{.Price}}"))
			},
			key:      "delivery.datetime.price",
			lang:     "fr-FR",
			delivery: &Delivery{Date: "demain", Price: 123},
			msg:      "Bonjour, votre date de livraison est le demain et le prix est le 123",
		},
		"ok with changed table names": {
			before: func(i *Intl) {
				i.SetTblNames("tbl1", "tbl2")
				mock.ExpectQuery("SELECT .*").
					WithArgs("delivery.datetime.price", "fr-FR").
					WillReturnRows(sqlmock.NewRows([]string{"message", "translation"}).
						AddRow("Hi, your delivery date is {{.Date}} and a price is {{.Price}}",
							"Bonjour, votre date de livraison est le {{.Date}} et le prix est le {{.Price}}"))
			},
			key:      "delivery.datetime.price",
			lang:     "fr-FR",
			delivery: &Delivery{Date: "demain", Price: 123},
			msg:      "Bonjour, votre date de livraison est le demain et le prix est le 123",
		},
		"err sql: no rows in result set": {
			before: func(i *Intl) {
				mock.ExpectQuery("SELECT .*").
					WithArgs("delivery.datetime.price", "fr-FR").
					WillReturnRows(sqlmock.NewRows([]string{"message", "translation"}))
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
