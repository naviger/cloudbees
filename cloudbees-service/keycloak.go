package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

var baseAdminUrl string = "https://cloudbees.dev:8443/admin/realms/cloudbees/"
var baseAuthUrl string = "https://cloudbees.dev:8443/realms/cloudbees/"
var adminUrl string = "protocol/openid-connect/token"
var userPostUrl string = "users"
var userGetUrl string = "users?username=%s"
var clientUrl string = "clients?clientId=%s"
var clientRoleUrl = "users/%s/role-mappings/clients/%s"
var grantType string = "client_credentials"
var clientId string = "admin-cli"
var clientSecret string = "WlkZGEFIY8OREsky9UgDiWi20mARa4yz"
var customerRole string = "861d0bc9-aca8-493f-bea5-069081c1f076"

func GetAdminToken() string {

	type Token struct{ *string }

	type TokenDoc struct {
		AccessToken string `json:"access_token"`
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   time.Duration(6000) * time.Second,
		Transport: tr,
	}

	form := url.Values{}
	form.Add("grant_type", grantType)
	form.Add("client_id", clientId)
	form.Add("client_secret", clientSecret)

	req, _ := http.NewRequest("POST", baseAuthUrl+adminUrl, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	//authHdr := res.Header.Get("Authorization")
	body, _ := io.ReadAll(res.Body)
	rawToken := TokenDoc{AccessToken: ""}

	err = json.Unmarshal(body, &rawToken)
	if err != nil {
		panic(err)
	}

	//fmt.Println("Token: ", rawToken, res)
	//log.Println("RESULT:", res, string(body[:]), err, authHdr)

	res.Header.Get("authorization")
	return rawToken.AccessToken
}

func GetClientId(token string, clientId string) string {
	if len(token) == 0 {
		token = GetAdminToken()
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   time.Duration(6000) * time.Second,
		Transport: tr,
	}

	req, _ := http.NewRequest("GET", baseAdminUrl+fmt.Sprintf(clientUrl, clientId), nil)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, _ := client.Do(req)
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		body, _ := io.ReadAll(res.Body)
		//log.Println("GET CLIENT ID SUCCESS: ", string(body[:]), res, req.URL)
		id := gjson.Get(string(body[:]), "0.id")
		//log.Println("ClientID: ", token)
		return id.String()

	} else {
		//log.Println("GET CLIENT ID FAILURE: ", res, req.URL)
		return ""
	}
}

func GetUser(token string, username string) User {
	if len(token) == 0 {
		token = GetAdminToken()
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   time.Duration(6000) * time.Second,
		Transport: tr,
	}

	req, _ := http.NewRequest("GET", baseAdminUrl+fmt.Sprintf(userGetUrl, username), nil)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, _ := client.Do(req)
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		body, _ := io.ReadAll(res.Body)
		id := gjson.Get(string(body[:]), "0.id").String()
		un := gjson.Get(string(body[:]), "0.username").String()
		fn := gjson.Get(string(body[:]), "0.firstName").String()
		ln := gjson.Get(string(body[:]), "0.lastName").String()
		en := gjson.Get(string(body[:]), "0.enabled").Bool()
		em := gjson.Get(string(body[:]), "0.email").String()

		return User{
			Id:        id,
			FirstName: fn,
			LastName:  ln,
			Username:  un,
			Enabled:   en,
			Email:     em,
		}
	} else {
		return User{}
	}
}

func DeleteUser(token string, username string) bool {
	if len(token) == 0 {
		token = GetAdminToken()
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   time.Duration(6000) * time.Second,
		Transport: tr,
	}

	user := GetUser(token, username)

	sUrl := baseAdminUrl + userPostUrl + "/" + user.Id
	req, _ := http.NewRequest("DELETE", sUrl, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, _ := client.Do(req)
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return true
	} else {
		return false
	}
}

func CreateUser(token string, firstname string, lastname string, email string) string {

	if len(token) == 0 {
		token = GetAdminToken()
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   time.Duration(6000) * time.Second,
		Transport: tr,
	}

	creds := Credentials{
		Type:      "password",
		Value:     "Test123!",
		Temporary: false,
	}

	credArray := make([]Credentials, 0)
	credArray = append(credArray, creds)

	var user User = User{
		FirstName:   firstname,
		LastName:    lastname,
		Email:       email,
		Username:    firstname + "." + lastname,
		Enabled:     true,
		Credentials: credArray,
	}

	u, _ := json.Marshal(user)
	sUrl := baseAdminUrl + userPostUrl
	req, _ := http.NewRequest("POST", sUrl, bytes.NewBuffer(u))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, _ := client.Do(req)

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		location, _ := res.Location()
		usrUrl := strings.Split(location.String(), "/")
		userId := usrUrl[len(usrUrl)-1]
		return userId
	} else {
		return ""
	}
}

func AssignUserAsClientRole(token string, userId string, role string) bool {
	if len(token) == 0 {
		token = GetAdminToken()
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   time.Duration(6000) * time.Second,
		Transport: tr,
	}

	clientId = GetClientId(token, "cloudbees-client")

	fullUrl := baseAdminUrl + fmt.Sprintf(clientRoleUrl, userId, clientId)
	clientRole := ClientRole{
		Id:          customerRole,
		Name:        role,
		Composite:   false,
		ClientRole:  true,
		ContainerId: clientId,
	}

	rolesArray := make([]ClientRole, 0)
	rolesArray = append(rolesArray, clientRole)

	c, _ := json.Marshal(rolesArray)

	req, _ := http.NewRequest("POST", fullUrl, bytes.NewBuffer(c))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, _ := client.Do(req)

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return true
	} else {
		return false
	}
}

type Credentials struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

type User struct {
	Id          string        `json:"id"`
	Username    string        `json:"username"`
	FirstName   string        `json:"firstName"`
	LastName    string        `json:"lastName"`
	Email       string        `json:"email"`
	Enabled     bool          `json:"enabled"`
	Credentials []Credentials `json:"credentials"`
}

type ClientRole struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Composite   bool   `json:"composite"`
	ClientRole  bool   `json:"clientRole"`
	ContainerId string `json:"containerId"`
}
