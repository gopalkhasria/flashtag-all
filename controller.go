package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
)

type PageData struct {
	Name     string `json:"name"`
	Cover    string `json:"cover"`
	Avatar   string `json:"avatar"`
	Links    []Link `json:"links"`
	Public   bool   `json:"public"`
	Visits   int    `json:"visits"`
	Email    string `json:"email"`
	Number   string `json:"number"`
	Number2  string `json:"number2"`
	Bio      string `json:"bio"`
	Show     string `json:"show"`
	Color    string `json:"color"`
	Username string `json:"username"`
}

type Link struct {
	Type   string `json:"type"`
	Icon   string `json:"icon"`
	Value  string `json:"value"`
	Prefix string `json:"prefix"`
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var data PageData
	name := mux.Vars(r)["name"]
	data.Name = name
	data.Cover = "https://ftp.flashtag.it/background.jpeg"
	data.Avatar = "https://ftp.flashtag.it/avatar.png"
	docID := ""
	data.Username = name
	iter := Client.Collection("users").Where("name", "==", name).Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			return
		}
		tmp := doc.Data()["public"]
		tmpStr := fmt.Sprintf("%v", tmp)
		if tmpStr == "false" {
			http.ServeFile(w, r, "private.html")
			return
		}
		docID = doc.Ref.ID
		tmp = doc.Data()["links"]
		tmpStr = fmt.Sprintf("%v", tmp)
		err = json.Unmarshal([]byte(tmpStr), &data.Links)
		if err != nil {
			log.Println(err)
			return
		}
		tmp = doc.Data()["bio"]
		data.Bio = fmt.Sprintf("%v", tmp)
		tmpStr = fmt.Sprintf("%v", doc.Data()["cover"])
		if tmpStr != "" {
			data.Cover = tmpStr
		}
		tmpStr = fmt.Sprintf("%v", doc.Data()["avatar"])
		if tmpStr != "" && doc.Data()["avatar"] != nil {
			data.Avatar = tmpStr
		}
		tmpStr = fmt.Sprintf("%v", doc.Data()["visits"])
		visit, _ := strconv.Atoi(tmpStr)
		data.Visits = visit

		tmpStr = fmt.Sprintf("%v", doc.Data()["show"])
		data.Show = tmpStr

		tmpStr = fmt.Sprintf("%v", doc.Data()["color"])
		data.Color = tmpStr
		data.Color = data.Color[2:]
	}
	if len(data.Links) == 0 {
		http.ServeFile(w, r, "404.html")
		return
	}
	for _, l := range data.Links {
		if l.Type == "EMAIL" {
			data.Email = l.Value
		}
		if l.Type == "NAME" {
			data.Name = l.Value
		}
		if l.Type == "PHONE" {
			data.Number = l.Value
		}
		if l.Type == "PHONE2" {
			data.Number2 = l.Value
		}
	}
	_, err := r.Cookie("visited")
	if err != nil {
		visted := &http.Cookie{
			Name:   "visited",
			Value:  "true",
			MaxAge: 86400,
		}
		http.SetCookie(w, visted)
		data.Visits++
		var update []firestore.Update
		update = append(update, firestore.Update{
			Path:  "visits",
			Value: data.Visits,
		})
		Client.Collection("users").Doc(docID).Update(context.Background(), update)

	}
	if data.Show == "false" {
		http.Redirect(w, r, data.Links[0].Prefix+data.Links[0].Value, http.StatusTemporaryRedirect)
		return
	}
	t := template.Must(template.ParseFiles("index.html"))
	err = t.Execute(w, data)
	if err != nil {
		log.Println(err)
	}

}
