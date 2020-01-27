package ss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	SYNC    = "/api/sync_ss/"
	TRAFFIC = "/api/traffic_ss/"
)

func Sync(c *Conn, users map[string]string, url string, token string) {
	client := &http.Client{Timeout: 3 * time.Second}
	req, _ := http.NewRequest("GET", "https://"+url+SYNC, nil)
	req.Header.Set("token", token)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Println("Non-OK HTTP status:", res.StatusCode)
		return
	}
	data, _ := ioutil.ReadAll(res.Body)
	var nu = make(map[string]string)

	json.Unmarshal([]byte(data), &nu)

	//fmt.Println(nu)
	for key, value := range nu {
		v, ok := users[key]
		if !ok {
			resp, _ := c.AddUser(key, value)
			fmt.Println("add new user" + resp)
		}
		if ok && v != value {
			resp, _ := c.RemoveUser(key)
			fmt.Println("remove user" + resp)
			resp1, _ := c.AddUser(key, value)
			fmt.Println("add new user" + resp1)
		}
	}

	for k := range users {
		if _, ok := nu[k]; !ok {
			resp, _ := c.RemoveUser(k)
			fmt.Println("remove user" + resp)
		}
	}

	for key := range users {
		delete(users, key)
	}
	//TODO  maybe need sync.map ??
	for k, v := range nu {
		users[k] = nu[v]
	}
}

func SyncTraffic(c *Conn, url string, token string) {

	result, _ := c.Traffic()
	fmt.Println(result)

	var data = make(map[string]string)
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(data)

	client := &http.Client{Timeout: 3 * time.Second}
	req, _ := http.NewRequest("POST", "https://"+url+TRAFFIC, buf)
	req.Header.Set("token", token)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("Non-OK HTTP status:", res.StatusCode)
		return
	}
}
