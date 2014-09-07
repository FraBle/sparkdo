package main

import (
	// standard library
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	// external packages
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/hoisie/mustache"
)

var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

// homeHandler returns the index.html.
func homeHandler(rw http.ResponseWriter, req *http.Request) {
	var body []byte
	session, err := store.Get(req, CONFIG.SessionName)
	if err != nil {
		log.Printf("Error reading session: %v", err)
		if cookie, err := req.Cookie(CONFIG.SessionName); err != nil {
			log.Printf("Error reading cookie: %v", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		} else {
			cookie.MaxAge = -1
			http.SetCookie(rw, cookie)
			http.Redirect(rw, req, "/", http.StatusFound)
		}
		return
	}
	if accesstoken, ok := session.Values["accesstoken"]; session.IsNew || !ok {
		body = []byte(mustache.RenderFileInLayout("../src/github.com/FraBle/sparkdo/templates/welcome.html", "../src/github.com/FraBle/sparkdo/templates/layouts/index.html", nil))
	} else {
		client := &http.Client{}
		dropletReq, err := http.NewRequest("GET", "https://api.digitalocean.com/v2/droplets", nil)
		if err != nil {
			log.Printf("Error creating droplet request: %v", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		dropletReq.Header.Add("Authorization", "Bearer "+accesstoken.(string))
		dropletResp, err := client.Do(dropletReq)
		if err != nil {
			log.Printf("Error getting droplets response: %v", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		defer dropletResp.Body.Close()
		data, err := ioutil.ReadAll(dropletResp.Body)
		if err != nil {
			log.Printf("Error reading droplets response: %v", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		var droplets Droplets
		err = json.Unmarshal(data, &droplets)
		if err != nil {
			log.Printf("Error reading droplets: %v", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		name, _ := session.Values["name"]
		dropletMap := make(map[string][]interface{})
		dropletMap["droplets"] = []interface{}{}
		dropletMap["meta"] = []interface{}{}
		dropletMap["meta"] = append(dropletMap["meta"], map[string]interface{}{"dropletSize": len(droplets.Droplets), "name": name})
		for _, droplet := range droplets.Droplets {
			dropletMap["droplets"] = append(dropletMap["droplets"], map[string]string{"name": droplet.Name, "status": droplet.Status, "id": strconv.Itoa(droplet.Id)})
		}
		body = []byte(mustache.RenderFileInLayout("../src/github.com/FraBle/sparkdo/templates/adminpanel.html", "../src/github.com/FraBle/sparkdo/templates/layouts/index.html", dropletMap))
	}
	rw.Write(body)
}

// loginHandler handles Digital Ocean login callback.
func loginHandler(rw http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, CONFIG.SessionName)
	if err != nil {
		log.Printf("Error reading session: %v", err)
		if cookie, err := req.Cookie(CONFIG.SessionName); err != nil {
			log.Printf("Error reading cookie: %v", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		} else {
			cookie.MaxAge = -1
			http.SetCookie(rw, cookie)
			http.Redirect(rw, req, "/", http.StatusFound)
		}
		return
	}
	code := req.FormValue("code")
	client := &http.Client{}
	resp, err := client.PostForm("https://cloud.digitalocean.com/v1/oauth/token",
		url.Values{
			"client_id":     {CONFIG.ClientId},
			"client_secret": {CONFIG.ClientSecret},
			"code":          {code},
			"grant_type":    {"authorization_code"},
			"redirect_uri":  {CONFIG.CallbackUrl},
		})
	if err != nil {
		log.Printf("Error reading cookie: %v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading cookie: %v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	var credentials DigitalOceanResponse
	err = json.Unmarshal(body, &credentials)
	session.Values["accesstoken"] = credentials.AccessToken
	session.Values["name"] = credentials.Info.Name
	session.Save(req, rw)
	http.Redirect(rw, req, "/", http.StatusFound)
}

// sparkHandler handles updates for spark credentials.
func sparkHandler(rw http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, CONFIG.SessionName)
	if err != nil {
		log.Printf("Error reading session: %v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["spark_accesstoken"] = req.FormValue("accesstoken")
	session.Values["spark_deviceid"] = req.FormValue("deviceid")
	session.Save(req, rw)
}

// monitorHandler handles request for monitoring a DigitalOcean Droplet with a Spark core.
func monitorHandler(rw http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, CONFIG.SessionName)
	if err != nil {
		log.Printf("Error reading session: %v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	dropletId := req.FormValue("dropletid")
	sparkDeviceId := session.Values["spark_deviceid"]
	sparkAccessToken := session.Values["spark_accesstoken"]
	digitalOceanAccessToken := session.Values["accesstoken"]
	go monitorDroplet(dropletId, sparkDeviceId.(string), sparkAccessToken.(string), digitalOceanAccessToken.(string), &http.Client{})
}
