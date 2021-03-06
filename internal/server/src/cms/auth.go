package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/julienschmidt/httprouter"
)

const monthSeconds = 60 * 60 * 24 * 30
const cookieName = "authtkt"

// these vars will drive the auth system
var secretKey []byte
var nonceCtr uint64
var noncePrefix []byte
var verificationNonce []byte

var lastNonceGenerationTime time.Time

var cipherInterface cipher.AEAD

func initCipher() {
	secretKey = make([]byte, 32)
	noncePrefix = make([]byte, 4)
	verificationNonce = make([]byte, 32)
	rand.Read(secretKey)
	rand.Read(noncePrefix)
	rand.Read(verificationNonce)
	lastNonceGenerationTime = time.Now()
	block, _ := aes.NewCipher(secretKey)
	var err error
	cipherInterface, err = cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}
}

func generateNewGCMNonce() []byte {
	ctr := make([]byte, 8)
	binary.BigEndian.PutUint64(ctr, atomic.AddUint64(&nonceCtr, 1))
	return append(noncePrefix, ctr...)
}

func generateSignedAuthToken() []byte {
	//we will do a simple signed random nonce to check against
	nonce := generateNewGCMNonce()
	token := append(nonce, cipherInterface.Seal(nil, nonce, verificationNonce, nil)...)
	return token
}

func isTokenValid(token []byte) bool {
	//if it has been past one month since the last nonce, force return false
	//after generating a new nonce
	if time.Since(lastNonceGenerationTime).Seconds() > monthSeconds {
		verificationNonce = make([]byte, 32)
		rand.Read(verificationNonce)
		lastNonceGenerationTime = time.Now()
		return false
	}

	nonce := token[:12]
	data := token[12:]
	decodedVal, err := cipherInterface.Open(nil, nonce, data, nil)
	if err != nil || string(decodedVal) != string(verificationNonce) {
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
			authorized = isTokenValid(tokenBytes)
		}

		if authorized {
			f(_w, _r, _p)
		} else {
			http.Redirect(_w, _r, "/login", http.StatusFound)
		}
	}
}

func loginView(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := BaseData(r, "CorvusCrypto.com - Login")
	err := globalTemplate.ExecuteTemplate(w, "login", data)
	if err != nil {
		fmt.Println(err)
	}
}

func loginViewWithError(w http.ResponseWriter, r *http.Request) {
	data := BaseData(r, "CorvusCrypto.com - Login")
	data["Error"] = true
	err := globalTemplate.ExecuteTemplate(w, "login", data)
	if err != nil {
		fmt.Println(err)
	}
}

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	username := strings.TrimSpace(strings.ToLower(r.PostFormValue("username")))
	password := []byte(r.PostFormValue("password"))

	// first check for username consistency
	if username != Config.OwnerUsername {
		loginViewWithError(w, r)
		return
	}

	//check password matching
	err := bcrypt.CompareHashAndPassword(Config.passHash, password)
	if err != nil {
		loginViewWithError(w, r)
		return
	}

	//create the token and base64 encode it
	token := generateSignedAuthToken()
	encodedToken := base64.StdEncoding.EncodeToString(token)
	//create the cookie to hold the signed token
	tokenCookie := new(http.Cookie)
	tokenCookie.Name = cookieName
	tokenCookie.HttpOnly = true
	tokenCookie.MaxAge = 60 * 60 * 24 * 30 //one month expiration
	tokenCookie.Value = encodedToken
	tokenCookie.Domain = ""

	http.SetCookie(w, tokenCookie)
	//generate data and do redirect and get
	http.Redirect(w, r, "/", http.StatusFound)
}
