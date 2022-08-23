# intl

Internationalization library based on different DB storages

[![Go Report Card](https://goreportcard.com/badge/github.com/arthurkushman/intl)](https://goreportcard.com/report/github.com/arthurkushman/intl)
[![Build and test](https://github.com/arthurkushman/intl/actions/workflows/go.yml/badge.svg)](https://github.com/arthurkushman/intl/actions/workflows/go.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![GoDoc](https://github.com/golang/gddo/blob/c782c79e0a3c3282dacdaaebeff9e6fd99cb2919/gddo-server/assets/status.svg)](https://godoc.org/github.com/arthurkushman/intl)

### Installation

```go
import "github.com/arthurkushman/intl"
```

### Migrations
You'll need a migration for particular database engine: postgresql, mysql, mssql, sqlite, oci which u can find in db/migrations folder.
After running one u'll be set up, then u need to insert data e.g.:

`source_message` table
```
1,delivery.datetime.price,"Hi, your delivery date is {{.Date}} and a price is {{.Price}}"
```
`message` table (in this example the last field localize_config u need only if using plurals)
```
1,fr-FR,"Bonjour, votre date de livraison est le {{.Date}} et le prix est le {{.Price}}","{\"DefaultMessage\": {\"ID\": \"Delivery\", \"One\": \"Bonjour, votre date de livraison est le {{.Date}} et le prix est le {{.Price}}\", \"Other\": \"Bonjour, votre date de livraison est le {{.Date}} many et le prix est le {{.Price}} many\"}, \"TemplateData\": {\"Date\": \"demain\", \"Price\": 123}, \"PluralCount\": 2}"
```

### Message translation example

```go 
// assuming that in db there is a source message in English: "Hi, your delivery date is {{.Date}} and a price is {{.Price}}"
// and for French (in message table) is: "Bonjour, votre date de livraison est le {{.Date}} et le prix est le {{.Price}}"

msg, err := intl.Translate("some.actual.key", "fr-FR", &Delivery{Date: "demain", Price: 123})
// msg: "Bonjour, votre date de livraison est le demain et le prix est le 123"
```

### Plural message translation example

```go
// to set up plurals correctly see an example above + intl_test.go
msg, err := intl.TranslatePlurals("some.actual.key", "fr-FR")
// msg: "Bonjour, votre date de livraison est le demain many et le prix est le 123 many"
```