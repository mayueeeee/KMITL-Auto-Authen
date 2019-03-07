package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

type UserResponse struct {
	Data struct {
		AccessStatus           int         `json:"accessStatus"`
		Account                string      `json:"account"`
		AccountMac             string      `json:"accountMac"`
		AccountRemainDays      int         `json:"accountRemainDays"`
		AccountexpiredTime     interface{} `json:"accountexpiredTime"`
		AssignedFlow           bool        `json:"assignedFlow"`
		AssignedTime           bool        `json:"assignedTime"`
		CanRegisterApplication bool        `json:"canRegisterApplication"`
		CanVisitorApplication  bool        `json:"canVisitorApplication"`
		FailCount              int         `json:"failCount"`
		FirstFlowLogin         bool        `json:"firstFlowLogin"`
		FirstTimeLogin         bool        `json:"firstTimeLogin"`
		HasBindTelFlag         interface{} `json:"hasBindTelFlag"`
		IncludeCharLT          bool        `json:"includeCharLT"`
		IncludeNumber          bool        `json:"includeNumber"`
		IncludeSpecialChar     bool        `json:"includeSpecialChar"`
		IP                     string      `json:"ip"`
		IsNeedBindTel          interface{} `json:"isNeedBindTel"`
		LoginDate              string      `json:"loginDate"`
		LoginType              int         `json:"loginType"`
		Message                interface{} `json:"message"`
		OpenID                 interface{} `json:"openId"`
		PermitUpdatePwd        bool        `json:"permitUpdatePwd"`
		PortalAuth             bool        `json:"portalAuth"`
		PortalAuthStatus       int         `json:"portalAuthStatus"`
		PortalErrorCode        int         `json:"portalErrorCode"`
		PwdMaxLen              int         `json:"pwdMaxLen"`
		PwdMinLen              int         `json:"pwdMinLen"`
		PwdRemainDays          int         `json:"pwdRemainDays"`
		PwdexpiredTime         interface{} `json:"pwdexpiredTime"`
		RedirectURL            interface{} `json:"redirectUrl"`
		//ResidualFlow            int         `json:"residualFlow"`
		//ResidualTime            int         `json:"residualTime"`
		SessionID               string      `json:"sessionId"`
		StatusCode              int         `json:"statusCode"`
		TimeOutPeriod           int         `json:"timeOutPeriod"`
		Token                   interface{} `json:"token"`
		TotalTime               int         `json:"totalTime"`
		UserName                interface{} `json:"userName"`
		WebHeatbeatPeriod       int         `json:"webHeatbeatPeriod"`
		WebPortalOvertimePeriod int         `json:"webPortalOvertimePeriod"`
	} `json:"data"`
	Message        interface{} `json:"message"`
	SessionTimeOut bool        `json:"sessionTimeOut"`
	Success        bool        `json:"success"`
	Token          string      `json:"token"`
	Total          int         `json:"total"`
}

type HeartBeatResponse struct {
	Data           string      `json:"data"`
	Message        interface{} `json:"message"`
	SessionTimeOut bool        `json:"sessionTimeOut"`
	Success        bool        `json:"success"`
	Token          interface{} `json:"token"`
	Total          int         `json:"total"`
}

var (
	Client            http.Client
	CookieJar         *cookiejar.Jar
	heartbeatInterval int
	heartbeatState    = true
	heartbeatTimer    time.Time
	loginCount        int
	connectionState   bool
)

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

func sendLoginRequest(username string, password string) bool {
	formData := url.Values{"userName": {username}, "password": {password}, "browserFlag": {"en"}}
	resp, err := Client.PostForm("http://10.252.23.101:8080/PortalServer/Webauth/webAuthAction!login.action", formData)
	errorHandler(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	respData := UserResponse{}
	jsonErr := json.Unmarshal(body, &respData)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	if respData.Success {
		heartbeatInterval = respData.Data.WebHeatbeatPeriod
		heartbeatTimer = time.Now().Add(time.Millisecond * time.Duration(heartbeatInterval))
		log.Println("Successful login at", respData.Data.LoginDate)
		return true
	} else {
		log.Println("Server: ", respData.Message)
		panic("Unexpected error from server, Bye!")
		return false
	}
}

func checkConnection() bool {
	resp, err := http.Get("http://detectportal.firefox.com/")
	errorHandler(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return strings.TrimRight(string(body), "\n") == "success"
}

func sendHeartBeat(username string) bool {
	formData := url.Values{"userName": {username}}
	resp, err := Client.PostForm("http://10.252.23.101:8080/PortalServer/Webauth/webAuthAction!hearbeat.action", formData)
	errorHandler(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	respData := HeartBeatResponse{}
	jsonErr := json.Unmarshal(body, &respData)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	if (respData.Data == "ONLINE") {
		heartbeatTimer = time.Now().Add(time.Millisecond * time.Duration(heartbeatInterval))
	}
	return respData.Data == "ONLINE"
}

func syncState() {
	formData := url.Values{}
	resp, err := Client.PostForm("http://10.252.23.101:8080/PortalServer/Webauth/webAuthAction!syncPortalAuthResult.action", formData)
	errorHandler(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	respData := UserResponse{}
	jsonErr := json.Unmarshal(body, &respData)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}

func doLogin(username string, password string) {
	sendLoginRequest(username, password)
	syncState()
}

func main() {
	fmt.Println("Auto Authen By Mayueeeee: CE KMITL")
	err := godotenv.Load()
	errorHandler(err)
	CookieJar, _ = cookiejar.New(nil)
	Client = http.Client{
		Jar: CookieJar,
	}
	username := os.Getenv("KMITL_USERNAME")
	password := os.Getenv("KMITL_PASSWORD")
	log.Println("Try to login as", username)
	resetTimer := time.Now().Add(time.Hour * time.Duration(8))
	doLogin(username, password)
	time.Sleep(30 * time.Second)

	for {
		if resetTimer.Before(time.Now()) {
			log.Println("Session timeout, Login again")
			doLogin(username, password)
			goto Sleep

		}
		if heartbeatTimer.Before(time.Now()) {
			log.Println("Send heartbeat ❤")
			heartbeatState := sendHeartBeat(username)
			if !heartbeatState {
				log.Println("Send heartbeat failed. Re-login again")
				doLogin(username, password)
			}
		}
		connectionState = checkConnection()
		if connectionState {
			log.Println("Connection OK, Yay!")
		} else {
			log.Println("Require auth, Try to login")
			doLogin(username, password)
		}
	Sleep:
		time.Sleep(30 * time.Second)
	}
}