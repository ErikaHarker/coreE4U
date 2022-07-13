package controllers

import (
	"coreapp/app/controllers/authentication"
	"coreapp/app/controllers/authentication/facebookOauth2/models"
	settings "coreapp/conf"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	"github.com/revel/revel"
)

type AuthApp struct {
	*revel.Controller
}

func (c AuthApp) connected() *models.User {
	return c.ViewArgs["user"].(*models.User)
}

func (c AuthApp) IndexTest() revel.Result {
	a := "FACEBOOK"
	b := "saurio"
	mapa := make(map[string]interface{})
	mapa["a"] = a
	mapa["b"] = b
	c.ViewArgs["mapa"] = mapa
	return c.RenderTemplate("App/Index.html")
}

func (c AuthApp) IndexFB() revel.Result {
	fmt.Println("-------")
	u := c.connected()
	me := map[string]interface{}{}
	if u != nil && u.AccessToken != "" {
		resp, _ := http.Get("https://graph.facebook.com/me?access_token=" +
			url.QueryEscape(u.AccessToken))
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&me); err != nil {
			c.Log.Error("json decode error", "error", err)
		}
		c.Log.Info("Data fetched", "data", me)
	}

	authUrl := settings.FACEBOOK.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return c.Render(me, authUrl)
}

func (c AuthApp) LoginSocial() revel.Result {
	login_url, err := authentication.LoginSocialUrl("saurio", "facebook")
	if err != nil {
		return c.Redirect(AuthApp.IndexTest)
	}
	return c.Redirect(login_url)
}

func (c AuthApp) AuthSocial() revel.Result {
	code := c.Params.Get("code")
	user_id := c.Params.Get("state")
	token, user_id, err := authentication.TokenSocial(code, user_id, "facebook")
	if err != nil {
		return c.Redirect(AuthApp.IndexTest)
	}
	token_user, err := authentication.TokenUserBySocial(user_id, "facebook")
	mapa := make(map[string]interface{})
	mapa["token_social"] = token
	mapa["token_user"] = token_user.Token
	mapa["user_id"] = user_id
	return c.RenderJSON(mapa)
}

func (c AuthApp) LoginUser() revel.Result {
	var jsonData map[string]interface{}
	c.Params.BindJSON(&jsonData)

	username := jsonData["username"].(string)
	password := jsonData["password"].(string)

	fmt.Println(username, password)

	token, err := authentication.Login(username, password)
	if err != nil {
		return c.Redirect(AuthApp.IndexTest)
	}
	mapa_r := make(map[string]interface{})
	mapa_r["token_social"] = token.Token
	mapa_r["username"] = username
	return c.RenderJSON(mapa_r)
}
