package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Loca struct {
	ID       int      `json:"id"`
	Location []string `json:"locations"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Locations struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func main() {
	http.HandleFunc("/Info", Info)
	http.HandleFunc("/", Home)
	fmt.Println("Server starting on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "ERROR", http.StatusNotFound)
		return
	}

	// Parse the form values
	r.ParseForm()
	creationDateMin := r.FormValue("creationDateMin")
	creationDateMax := r.FormValue("creationDateMax")
	firstAlbumDateMin := r.FormValue("firstAlbumDateMin")
	firstAlbumDateMax := r.FormValue("firstAlbumDateMax")
	memberCounts := r.Form["memberCount"]
	locationSearch := r.FormValue("locationSearch")

	var DATA []Artist
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, "Failed to fetch artists data", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data, &DATA)
	if err != nil {
		http.Error(w, "Failed to unmarshal JSON", http.StatusInternalServerError)
		return
	}
	// Apply filters
	filteredData := []Artist{}
	for _, artist := range DATA {
		if firstAlbumDateMax != "" || firstAlbumDateMin != "" {
			FR := artist.FirstAlbum[len(artist.FirstAlbum)-4:]
			 if firstAlbumDateMax == ""{
				firstAlbumDateMax = "2024"
			 }
			 if firstAlbumDateMin == ""{
				firstAlbumDateMin = "1900"
			 }
			if firstAlbumDateMax < FR || firstAlbumDateMin > FR {
				continue
			}

		}
		// Filter for creation date
		if creationDateMin != "" || creationDateMax != "" {
			minDate, err := strconv.Atoi(creationDateMin)
			if err!= nil {
				minDate = 1900
			}
			maxDate, erre := strconv.Atoi(creationDateMax)
			if erre != nil {
				maxDate = 2024
			}
			if artist.CreationDate < minDate || artist.CreationDate > maxDate {
				continue
			}
		}

		// Filter by number of members
		if len(memberCounts) > 0 {
			found := false
			for _, count := range memberCounts {
				if count == "1" {
					if len(artist.Members) == 1 {
						found = true
					}
				}
				if count == "2" {
					if len(artist.Members) == 2 {
						found = true
					}
				}
				if count == "3" {
					if len(artist.Members) == 3 {
						found = true
					}
				}
				if count == "4+" {
					if len(artist.Members) >= 4 {
						found = true
					}
				}

			}
			if !found {
				continue
			}
		}

		// Filter location
		if locationSearch != "" {
			//  check
			if !contains(artist.Members, locationSearch) {
				continue
			}
		}

		filteredData = append(filteredData, artist)
	}

	tml, err := template.ParseFiles("Home.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	err = tml.Execute(w, map[string]interface{}{
		"DATA": filteredData,
	})
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}
}

func Info(w http.ResponseWriter, r *http.Request) {
	ID := r.FormValue("ID")

	// Fetch artist details
	var artist Artist
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists/" + ID)
	if err != nil {
		http.Error(w, "Failed to fetch artist data", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data, &artist)
	if err != nil {
		http.Error(w, "Failed to unmarshal JSON", http.StatusInternalServerError)
		return
	}

	// Fetch location data
	var locaData Loca
	responseLoca, err := http.Get("https://groupietrackers.herokuapp.com/api/locations/" + ID)
	if err != nil {
		http.Error(w, "Failed to fetch location data", http.StatusInternalServerError)
		return
	}
	defer responseLoca.Body.Close()

	dataLoca, err := io.ReadAll(responseLoca.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(dataLoca, &locaData)
	if err != nil {
		http.Error(w, "Failed to unmarshal JSON", http.StatusInternalServerError)
		return
	}

	// Fetch date data
	var dateData Date
	responseDat, err := http.Get("https://groupietrackers.herokuapp.com/api/dates/" + ID)
	if err != nil {
		http.Error(w, "Failed to fetch date data", http.StatusInternalServerError)
		return
	}
	defer responseDat.Body.Close()

	dataDat, err := io.ReadAll(responseDat.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(dataDat, &dateData)
	if err != nil {
		http.Error(w, "Failed to unmarshal JSON", http.StatusInternalServerError)
		return
	}

	// Fetch relation data
	var Relation Locations
	RelaData, err := http.Get("https://groupietrackers.herokuapp.com/api/relation/" + ID)
	if err != nil {
		http.Error(w, "Failed to fetch relation data", http.StatusInternalServerError)
		return
	}
	defer RelaData.Body.Close()

	Rdata, err := io.ReadAll(RelaData.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(Rdata, &Relation)
	if err != nil {
		http.Error(w, "Failed to unmarshal JSON", http.StatusInternalServerError)
		return
	}

	// Render the artist details page
	tmp, err := template.ParseFiles("Aristes.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	err = tmp.Execute(w, map[string]interface{}{
		"Location": locaData,
		"Artist":   artist,
		"Dates":    dateData,
		"Relation": Relation,
	})
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}
}

func contains(locations []string, search string) bool {
	for _, loc := range locations {
		if strings.Contains(strings.ToLower(loc), strings.ToLower(search)) {
			return true
		}
	}
	return false
}
