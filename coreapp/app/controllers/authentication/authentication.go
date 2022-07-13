package authentication

import (
	"coreapp/app/controllers/authentication/models"
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"coreapp/app/controllers/authentication/facebookOauth2"

	jwt "github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/request"
)

// se firman los token con llave privada
// se verifican con una llave publica
var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func init() {
	privateBytes, err := ioutil.ReadFile("./conf/keys/private.rsa")
	if err != nil {
		panic(err)
		//log.Fatal("No se pudeo leer el archivo privado")
	}
	publicBytes, err := ioutil.ReadFile("./conf/keys/public.rsa.pub")
	if err != nil {
		panic(err)
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		panic(err)
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		panic(err)
	}
}

func GenerateJWT(user models.User) string {
	claims := models.Claim{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	result, err := token.SignedString(privateKey)
	if err != nil {
		log.Fatal("No se pudo firmar el token -- ", err)
	}
	return result
}

func LoginSocialUrl(user_id string, social_login string) (string, error) {
	social_login = strings.ToLower(social_login)
	switch social_login {
	case "facebook":
		return facebookOauth2.LoginUrl(user_id), nil
	case "google":
		//TODO: Realizar api Google
		return "--", nil
	default:
		err := errors.New("Red social no valida")
		return "", err
	}
}

func TokenSocial(code string, user_id string, social_login string) (string, string, error) {

	social_login = strings.ToLower(social_login)
	switch social_login {
	case "facebook":
		token, user_id := facebookOauth2.TokenUser(code, user_id)
		return token, user_id, nil
	case "google":
		//TODO: Realizar api Google
		return "--", "--", nil
	default:
		err := errors.New("Red social no valida")
		return "", "", err
	}
}

func TokenUserBySocial(user_id string, social_login string) (models.ResponseToken, error) {
	var user models.User
	user.Name = user_id
	user.Password = ""
	user.Role = "admin"
	user.IsSocial = false
	token := GenerateJWT(user)

	result := models.ResponseToken{token}

	return result, nil
}

func Login(username string, password string) (models.ResponseToken, error) {

	var user models.User
	user.Name = username
	if username == "pepe" && password == "123" {
		user.Password = ""
		user.Role = "admin"
		user.IsSocial = false
		token := GenerateJWT(user)

		result := models.ResponseToken{token}

		return result, nil

	} else {
		Err := errors.New("Usuario o contraseña no validos")
		result := models.ResponseToken{}
		return result, Err

	}
}

func ValidateToken(w http.ResponseWriter, r *http.Request) {

	token, err := request.ParseFromRequestWithClaims(
		r, request.OAuth2Extractor, &models.Claim{},
		func(t *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

	//token_r := r.Header.Get("Authorization")
	//fmt.Println(token_r)

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				fmt.Fprintln(w, "Su token ha expirado")
				return
			case jwt.ValidationErrorSignatureInvalid:
				fmt.Fprintln(w, "La firma del token no coincide")
				return
			default:
				fmt.Fprintln(w, "Su token no es valido")
				return
			}
		default:
			fmt.Fprintln(w, "Su token no es valido")
			return
		}
	} else {
		if token.Valid {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprintln(w, "Bienvenido al sistema")
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Su token no es valido")
		}
	}
}
