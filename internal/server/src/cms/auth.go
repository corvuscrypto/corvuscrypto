package main

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/julienschmidt/httprouter"
)

const monthSeconds = 60 * 60 * 24 * 30
const cookieName = "authtkt"

// these vars will drive the auth system
var secretKey []byte
var secretNonce []byte
var verificationNonce []byte
var lastNonceGenerationTime time.Time

var cipherInterface cipher.AEAD

func generateSignedAuthToken() []byte {
	var token []byte
	//we will do a simple signed random nonce to check against
	cipherInterface.Seal(token, secretNonce, verificationNonce, nil)
	return token
}

func tokenValid(token []byte) bool {
	//if it has been past one month since the last nonce, force return false
	//after generating a new nonce
	if time.Since(lastNonceGenerationTime).Seconds() > monthSeconds {
		verificationNonce = make([]byte, 32)
		rand.Read(verificationNonce)
		return false
	}
	var decodedVal []byte
	_, err := cipherInterface.Open(decodedVal, secretNonce, token, nil)
	if err != nil {
		return false
	}
	return true
}

func checkAuth(f httprouter.Handle) httprouter.Handle {
	return func(_w http.ResponseWriter, _r *http.Request, _p httprouter.Params) {
		var authorized bool
		//get the token from the auth Cookie
		tokenCookie, err := _r.Cookie(cookieName)
		if err != nil {
			authorized = false
		} else {
			tokenBytes, err := base64.StdEncoding.DecodeString(tokenCookie.Value)
			if err != nil {
				authorized = false
			}
			//check the token
			authorized = tokenValid(tokenBytes)
		}

		if authorized {
			f(_w, _r, _p)
		} else {
			_w.WriteHeader(401)
		}
	}
}

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	username := strings.ToLower(r.PostFormValue("username"))
	password := []byte(r.PostFormValue("password"))

	// first check for username consistency
	if username != Config.OwnerUsername {
		//Send Failure data
		return
	}

	//check password matching
	err := bcrypt.CompareHashAndPassword(Config.OwnerPassword, password)
	if err != nil {
		//Send Failure data
		return
	}

	//create the token and base64 encode it
	token := generateSignedAuthToken()
	encodedToken := base64.StdEncoding.EncodeToString(token)

	//create the cookie to hold the signed token
	tokenCookie := new(http.Cookie)
	tokenCookie.Secure = true
	tokenCookie.Name = cookieName
	tokenCookie.HttpOnly = true
	tokenCookie.MaxAge = 60 * 60 * 24 * 30 //one month expiration
	tokenCookie.Value = encodedToken

	http.SetCookie(w, tokenCookie)
	//generate data and do redirect and get
}
