# intl

Internationalization library based on different DB storages

[![Go Report Card](https://goreportcard.com/badge/github.com/arthurkushman/intl)](https://goreportcard.com/report/github.com/arthurkushman/intl)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![GoDoc](https://github.com/golang/gddo/blob/c782c79e0a3c3282dacdaaebeff9e6fd99cb2919/gddo-server/assets/status.svg)](https://godoc.org/github.com/arthurkushman/intl)

### Installation

```go
import "github.com/arthurkushman/intl"
```

### Message translation example

```go 
// assuming that in db there is a source message in English: "Hi, your delivery date is {{.Date}} and a price is {{.Price}}"
// and for French (in message table) is: "Bonjour, votre date de livraison est le {{.Date}} et le prix est le {{.Price}}"

msg, err := intl.Translate("some.actual.key", "fr-FR", &Delivery{Date: "demain", Price: 123})
// msg: "Bonjour, votre date de livraison est le demain et le prix est le 123"
```