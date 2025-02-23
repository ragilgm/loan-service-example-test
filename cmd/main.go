package main

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"github.com/test/loan-service/internal"
	"github.com/typical-go/typical-go/pkg/typapp"
	// Important to enable dependency injectino
	_ "github.com/test/loan-service/internal/infra"
)

var bundle *i18n.Bundle

func main() {

	if err := typapp.StartApp(internal.Start, internal.Shutdown); err != nil {
		logrus.Fatal(err.Error())
	}
}
