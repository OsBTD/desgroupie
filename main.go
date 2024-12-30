package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

type Artists struct {
	Image          string   `json:"image"`
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	Members        []string `json:"members"`
	CreationDate   int      `json:"creationDate"`
	FirstAlbum     string   `json:"firstAlbum"`
	RelationsURL   string   `json:"relations"`
	DatesLocations Relations
}

type Relations struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func fetchData(url string, target interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}

func main() {
	artists := []Artists{}
	err := fetchData("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err != nil {
		fmt.Println("Error fetching artists")
		return
	}

	for i := 0; i < len(artists); i++ {
		var relation Relations
		err := fetchData(artists[i].RelationsURL, &relation)
		if err != nil {
			fmt.Println("Error fetching relation for artist", artists[i].ID)
		} else {
			artists[i].DatesLocations = relation
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/index3.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, artists)
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("templates"))))
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
//handle all errors by separating the home handler from the main, controlling methods, 404, 500 etc..
//forbbid access to static assets and all folders like that
//try putting the images on another server 