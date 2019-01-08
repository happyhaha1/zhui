package zhui

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Book struct {
	ID string `json:"_id"`
	Title string `json:"title"`
	Author string `json:"author"`
}

type Atoc struct {
	ID string `json:"_id"`
	Name string `json:"name"`
	LastChapter string `json:"lastChapter"`
}

type Chapter struct {
	Title string `json:"title"`
	Link string `json:"link"`
}

type Data struct {
	Books []Book `json:"books,omitempty"`
	Chapters []Chapter `json:"chapters,omitempty"`
}

func SearchBooks(searchKey string) ([]Book, error) {
	queryUrl := "http://api.zhuishushenqi.com/book/fuzzy-search?query="+url.QueryEscape(searchKey)
	resp, err := http.Get(queryUrl)

	bytes, err := ioutil.ReadAll(resp.Body)
	var data Data
	err = json.Unmarshal(bytes, &data)

	return data.Books,err
}

func SearchAtocs(book Book) ([]Atoc, error){
	queryUrl := "http://api.zhuishushenqi.com/atoc?view=summary&book="+book.ID
	resp, err := http.Get(queryUrl)

	bytes, err := ioutil.ReadAll(resp.Body)
	var atocs []Atoc
	err = json.Unmarshal(bytes, &atocs)

	return atocs,err
}

func Download(book Book,atoc Atoc, path string) error {
	chapters, e := SearchChapters(atoc)

	num := len(chapters)

	println(num)

	return e
}

func SearchChapters(atoc Atoc) ([]Chapter, error){
	queryUrl := "http://api.zhuishushenqi.com/atoc/"+atoc.ID+"?view=chapters"
	resp, err := http.Get(queryUrl)

	bytes, err := ioutil.ReadAll(resp.Body)
	var data Data
	err = json.Unmarshal(bytes, &data)

	return data.Chapters,err
}
