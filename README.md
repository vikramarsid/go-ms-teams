[![Build Status][travis-image]][travis-url]
[![Github Tag][githubtag-image]][githubtag-url]

[![Maintainability][codeclimate-image]][codeclimate-url]
[![codecov][codecov-image]][codecov-url]

[![Go Report Card][goreport-image]][goreport-url]
[![GoDoc][godoc-image]][godoc-url]
[![License][license-image]][license-url]

***

# go-ms-teams

> A package to send messages to Microsoft Teams (channels)

...

# Usage

To get the package, execute:

```
go get https://github.com/vikramarsid/go-ms-teams
```

To import this package, add the following line to your code:

```
import "github.com/vikramarsid/go-ms-teams"
```

And this is an example of a simple implementation ...

```
import (
	"github.com/vikramarsid/go-ms-teams"
)

func main() {
	_ = sendTheMessage()
}

func sendTheMessage() error {
	// init the client
	pts := Options{
		Timeout: 60 * time.Second,
		Verbose: true,
	}
	mstClient := NewClient(opts)
	
	mstClient := gomsteams.NewClient()

	// setup webhook url
	webhookUrl := "https://outlook.office.com/webhook/YOUR_WEBHOOK_URL_OF_TEAMS_CHANNEL"

	// setup message card
	msgCard := gomsteams.NewMessageCard()
	msgCard.Title = "Hello world"
	msgCard.Text = "Here are some examples of formatted stuff like <br> * this list itself  <br> * **bold** <br> * *italic* <br> * ***bolditalic***"
	msgCard.ThemeColor = "#DF813D"

	// send
	return mstClient.Send(webhookUrl, msgCard)
}
```

# <a id="links"></a>some useful links

* [Inspiration - Credits to](https://github.com/dasrick/go-teams-notify)
* [MS Teams - adaptive cards](https://docs.microsoft.com/de-de/outlook/actionable-messages/adaptive-card)
* [MS Teams - send via connectors](https://docs.microsoft.com/de-de/outlook/actionable-messages/send-via-connectors)
* [adaptivecards.io](https://adaptivecards.io/designer)

***

[travis-image]: https://travis-ci.org/vikramarsid/go-ms-teams.svg?branch=master
[travis-url]: https://travis-ci.org/vikramarsid/go-ms-teams

[githubtag-image]: https://img.shields.io/github/tag/vikramarsid/go-ms-teams.svg?style=flat
[githubtag-url]: https://github.com/vikramarsid/go-ms-teams

[codeclimate-image]: https://api.codeclimate.com/v1/badges/fe69cc992370b3f97d94/maintainability
[codeclimate-url]: https://codeclimate.com/github/vikramarsid/go-ms-teams/maintainability

[codecov-image]: https://codecov.io/gh/vikramarsid/go-ms-teams/branch/master/graph/badge.svg
[codecov-url]: https://codecov.io/gh/vikramarsid/go-ms-teams

[goreport-image]: https://goreportcard.com/badge/github.com/vikramarsid/go-ms-teams
[goreport-url]: https://goreportcard.com/report/github.com/vikramarsid/go-ms-teams

[godoc-image]: https://godoc.org/github.com/vikramarsid/go-ms-teams?status.svg
[godoc-url]: https://godoc.org/github.com/vikramarsid/go-ms-teams

[license-image]: https://img.shields.io/github/license/vikramarsid/go-ms-teams.svg?style=flat
[license-url]: https://github.com/vikramarsid/go-ms-teams/blob/master/LICENSE
