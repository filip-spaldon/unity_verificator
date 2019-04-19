package actions

import (
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr/v2"
)

var r *render.Engine
var assetsBox = packr.NewBox("../public")

func init() {
	r = render.New(render.Options{
		HTMLLayout:   "application.html",
		TemplatesBox: packr.New("../templates", "../templates"),
		AssetsBox:    assetsBox,
		Helpers:      render.Helpers{},
	})
}
