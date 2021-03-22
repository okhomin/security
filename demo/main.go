package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	url = "http://localhost:8888"
)

type Group struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Read  bool     `json:"read"`
	Write bool     `json:"write"`
	Users []string `json:"users"`
}

type File struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Content string   `json:"content"`
	Groups  []string `json:"groups"`
	Acls    []string `json:"acls"`
}

func main() {
	var firstUserToken string
	var secondUserToken string

	req, err := http.NewRequest(http.MethodPost, url+"/signup", bytes.NewReader([]byte(`{"login": "root", "password": "root"}`)))
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	firstUserToken = resp.Header.Get("Authorization")
	log.Println("---Singup root user")
	log.Println("root token = ", firstUserToken)
	log.Println("---")
	fmt.Println()

	req, err = http.NewRequest(http.MethodGet, url+"/groups", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", firstUserToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var g []Group
	if err := json.Unmarshal(b, &g); err != nil {
		panic(err)
	}
	log.Println("---Get all groups")
	log.Printf("[root]groups = %+v\n", g)
	log.Println("---")
	fmt.Println()

	req, err = http.NewRequest(http.MethodPost, url+"/file", bytes.NewReader([]byte(fmt.Sprintf(`{"name": "test", "content": "testtesttestestestestes", "groups":  ["%v"], "acls": []}`, g[0].ID))))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", firstUserToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var f File
	if err := json.Unmarshal(b, &f); err != nil {
		panic(err)
	}
	log.Println("---Create file test")
	log.Printf("[root]files = %+v\n", f)
	log.Println("---")
	fmt.Println()

	req, err = http.NewRequest(http.MethodPost, url+"/signup", bytes.NewReader([]byte(`{"login": "user", "password": "user"}`)))
	if err != nil {
		panic(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	secondUserToken = resp.Header.Get("Authorization")
	log.Println("---Singup user")
	log.Println("user token = ", secondUserToken)
	log.Println("---")
	fmt.Println()

	req, err = http.NewRequest(http.MethodGet, url+"/files", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", secondUserToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var fs []File
	if err := json.Unmarshal(b, &fs); err != nil {
		panic(err)
	}
	log.Println("---Get all files")
	log.Printf("[user]files = %+v\n", fs)
	log.Println("---")
	fmt.Println()

	req, err = http.NewRequest(http.MethodGet, url+"/file/"+f.ID, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", firstUserToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b, &f); err != nil {
		panic(err)
	}
	log.Println("---Get file " + f.ID)
	log.Printf("[root]file = %+v\n", f)
	log.Println("---")
	fmt.Println()

	req, err = http.NewRequest(http.MethodGet, url+"/file/"+f.ID, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", secondUserToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Println("---Get file " + f.ID)
	log.Printf("[user]file = %v\n", string(b))
	log.Println("---")
	fmt.Println()
}
