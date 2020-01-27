package ray

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	SYNC    = "/api/sync/"
	TRAFFIC = "/api/traffic/"
)

func Sync(c *ServiceClient, users map[string]string, url string, token string) {
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
			c.AddUser("proxy", key+"@ssbear", 0, value, 64)
			fmt.Println("add user key:" + value)
		}
		if ok && v != value {
			c.RemoveUser("proxy", key+"@ssbear")
			c.AddUser("proxy", key+"@ssbear", 0, value, 64)
		}
	}

	for k, v := range users {
		if _, ok := nu[k]; !ok {
			c.RemoveUser("proxy", k+"@ssbear")
			fmt.Println("remove user key:" + v)
		}
	}

	for key := range users {
		delete(users, key)
	}

	for k, v := range nu {
		users[k] = nu[v]
	}
}

func Traffic(c *ServiceClient, url string, token string) {

	result := c.Traffic("user>>>", true)
	if len(result) == 0 {
		return
	}

	var data = make(map[string]int64)
	for key, value := range result {
		newkey := strings.Split(key, ">>>")
		id := strings.Split(newkey[1], "@")
		data[id[0]+"@"+newkey[3]] = value
	}
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
