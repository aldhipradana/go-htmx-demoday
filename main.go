package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

var tmpl *template.Template

type PokemonResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func main() {
	// Parse all templates in the 'templates' directory
	tmpl = template.Must(template.ParseGlob("templates/*.html"))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/process", processHandler)
	http.HandleFunc("/success", successHandler)
	http.HandleFunc("/pokemon", listpokemonHandler)

	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func listpokemonHandler(w http.ResponseWriter, r *http.Request) {

	url := "https://pokeapi.co/api/v2/pokemon?limit=100&offset=0"

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return
	}

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Received status code", resp.StatusCode)
		return
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return
	}

	// Parse the JSON response
	var pokemonResponse PokemonResponse
	if err := json.Unmarshal(body, &pokemonResponse); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "pokemonList.html", pokemonResponse)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
	}
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("name")
	email := r.FormValue("email")
	studentId := r.FormValue("studentId")

	// Redirect to success page with registration details
	http.Redirect(w, r, fmt.Sprintf("/success?name=%s&email=%s&studentId=%s", name, email, studentId), http.StatusSeeOther)
}

func successHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")
	studentId := r.URL.Query().Get("studentId")

	data := struct {
		Name      string
		Email     string
		StudentID string
	}{
		Name:      name,
		Email:     email,
		StudentID: studentId,
	}

	err := tmpl.ExecuteTemplate(w, "success.html", data)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
	}
}
