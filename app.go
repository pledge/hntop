package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var templates = template.Must(template.ParseFiles("templates/index.html"))

func render(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl + ".html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "index")
}

func StoryHandler(w http.ResponseWriter, r *http.Request) {
	resp, _ := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	w.Header().Set("Content-type", "application/json")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var keys []int
	json.Unmarshal(body, &keys)

	//sleep to avoid rate limiting
	time.Sleep(10 * time.Millisecond)

	//select random item
	rand.Seed(time.Now().Unix())
	key := rand.Intn(100)

	storyResp, _ := http.Get(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", key))
	defer storyResp.Body.Close()
	body, _ = ioutil.ReadAll(storyResp.Body)
	fmt.Fprintf(w, "%s", body)
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resp, _ := http.Get(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/user/%s.json", id))
	w.Header().Set("Content-type", "application/json")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "%s", body)
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/story", StoryHandler)
	r.HandleFunc("/user/{id:.+}", UserHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Listening...")
	http.ListenAndServe(":3000", r)
}
