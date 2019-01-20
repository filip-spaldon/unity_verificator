package grifts

import (
	"github.com/Filip/diplomovka_UNITY/unity_verificator/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
