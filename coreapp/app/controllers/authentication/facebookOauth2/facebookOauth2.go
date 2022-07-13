package facebookOauth2

import (
	models "coreapp/app/controllers/authentication/facebookOauth2/models"
	settings "coreapp/conf"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/revel/revel"
)

type Application struct {
	*revel.Controller
}

// The following keys correspond to a test application
// registered on Facebook, and associated with the loisant.org domain.
// You need to bind loisant.org to your machine with /etc/hosts to
// test the application locally.

func LoginUrl(user_id string) string {
	fb_login_url := settings.FacebookAuth.AuthURL + "?"
	fb_login_url += "redirect_uri=" + settings.FacebookAuth.RedirectURL + "&"
	fb_login_url += "client_id=" + settings.FacebookAuth.ClientID + "&"
	fb_login_url += "state=" + user_id + "&"
	return fb_login_url
}

func TokenUser(code string, user_id string) (string, string) {
	fmt.Println("CODE: ", code)
	fmt.Println("user_id: ", user_id)

	fb_auth_url := settings.FacebookAuth.TokenURL + "?"
	fb_auth_url += "redirect_uri=" + settings.FacebookAuth.RedirectURL + "&"
	fb_auth_url += "client_id=" + settings.FacebookAuth.ClientID + "&"
	fb_auth_url += "client_secret=" + settings.FacebookAuth.ClientSecret + "&"
	fb_auth_url += "code=" + code + "&grant_type=authorization_code"

	resp, err := http.Get(fb_auth_url)
	if err != nil {
		fmt.Println(err)
		fmt.Errorf("err")
	}
	data_fb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		fmt.Print(data_fb)
	}

	jsonStr := string(data_fb)
	fmt.Println(jsonStr)
	mapa := make(map[string]interface{})
	if err := json.Unmarshal([]byte(data_fb), &mapa); err != nil {
		panic(err)
	}

	token := mapa["access_token"].(string)

	return token, user_id
}

func setuser(c *revel.Controller) revel.Result {
	var user *models.User
	if _, ok := c.Session["uid"]; ok {
		uid, _ := strconv.ParseInt(c.Session["uid"].(string), 10, 0)
		user = models.GetUser(int(uid))
	}
	if user == nil {
		user = models.NewUser()
		c.Session["uid"] = fmt.Sprintf("%d", user.Uid)
	}
	c.ViewArgs["user"] = user
	return nil
}

func init() {
	revel.InterceptFunc(setuser, revel.BEFORE, &Application{})
}

func (c Application) connected() *models.User {
	return c.ViewArgs["user"].(*models.User)
}
