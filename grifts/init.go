package grifts

import (
	"github.com/Filip/unity_verificator/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
