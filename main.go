package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type CatFacts struct {
	Title   string `json:"Title"`
	Summary string `json:"summary"`
}

type SpaceNews struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type News struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

func handleNews(w http.ResponseWriter, r *http.Request) {
	space_news, err := fetchSpaceNews("https://api.spaceflightnewsapi.net/v4/articles/")

	if err != nil {
		http.Error(w, "Error has occured", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	cat_facts, err := fetchCatFacts("https://cat-fact.herokuapp.com/facts/")
	if err != nil {
		http.Error(w, "Error has occured", http.StatusInternalServerError)
		return
	}
	var news []News
	for i, sn, cf := 1, 1, 1; i <= 10; i++ {
		if i%3 != 0 && sn < len(space_news) {
			news = append(news, News{
				Title:   space_news[sn].Title,
				Summary: space_news[sn].Summary,
			})

			cf++
		} else {
			news = append(news, News{
				Title:   "Cat Fact",
				Summary: cat_facts[cf].Summary,
			})
			sn++
		}
	}
	jsonData, err := json.Marshal(news)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Write(jsonData)
}

func fetchCatFacts(url string) ([]CatFacts, error) {
	body, err := getRequest(url)
	if err != nil {
		return nil, err
	}

	var catfacts []CatFacts
	err = json.Unmarshal(body, &catfacts)
	if err != nil {
		return nil, err
	}
	return catfacts, nil
}

func fetchSpaceNews(url string) ([]SpaceNews, error) {
	body, err := getRequest(url)

	if err != nil {
		return nil, err
	}

	var spacenews []SpaceNews
	err = json.Unmarshal(body, &spacenews)
	if err != nil {
		fmt.Println("pop")
		return nil, err
	}
	return spacenews, nil
}

func getRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func newsTag(w http.ResponseWriter, req *http.Request) {
	v := req.FormValue("tag")
	fmt.Fprintln(w, "Do my search: "+v)
}

func main() {
	http.HandleFunc("/", foo)
	http.HandleFunc("/news", handleNews)
	http.ListenAndServe(":8000", nil)
}
