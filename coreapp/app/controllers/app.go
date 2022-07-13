package controllers

import (
	auth "coreapp/app/controllers/authentication"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	a := "pepe"
	b := "saurio"
	mapa := make(map[string]interface{})
	mapa["a"] = a
	mapa["b"] = b
	return c.Render(mapa)
}

func (c App) Index2() revel.Result {
	a := "pepe2"
	b := "saurio2"
	mapa := make(map[string]interface{})
	mapa["a"] = a
	mapa["b"] = b
	c.ViewArgs["mapa"] = mapa
	return c.RenderTemplate("App/Index.html")
}

type Stuff struct {
	Foo string ` json:"foo" xml:"foo" `
	Bar int    ` json:"bar" xml:"bar" `
}

func (c App) MyWork() revel.Result {
	data := make(map[string]interface{})
	data["error"] = nil
	stuff := Stuff{Foo: "xyz", Bar: 999}
	data["stuff"] = stuff
	return c.RenderJSON(data)
	// or alternately
	// return c.RenderXML(data)
}

func (c App) FacebookOauth() revel.Result {
	data := make(map[string]interface{})
	data["error"] = nil
	stuff := Stuff{Foo: "xyz", Bar: 999}
	data["stuff"] = stuff
	return c.RenderJSON(data)
	// or alternately
	// return c.RenderXML(data)
}

func (c App) Login() revel.Result {
	auth.Login("pepe", "123")
	data := make(map[string]interface{})
	data["error"] = nil
	stuff := Stuff{Foo: "xyz", Bar: 999}
	data["stuff"] = stuff
	return c.RenderJSON(data)
	// or alternately
	// return c.RenderXML(data)
}
