package actions

import "github.com/gobuffalo/buffalo"

// HomeHandler is for home page
func HomeHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("home.html"))
}

// IndexHandler is for routs page
func IndexHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("index.html"))
}

func DescHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("description.html"))
}

func ProgressHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("progress.html"))
}

func LinksxHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("links.html"))
}

func ResultHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("result.html"))
}

func FormHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("form.html"))
}
