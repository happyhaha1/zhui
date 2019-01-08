package zhui

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/schollz/progressbar"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

	bar := progressbar.New(num)
	page := 100
	var start = 0
	var end = 0

	tempDirPath, e := ioutil.TempDir("", book.Title)

	ch := make(chan int)
	var files []string
	for i := 0; end < num; i++ {
		start = 0 + i*page
		end = start + page
		if end > num {
			end = num
		}
		path := tempDirPath + "/" + book.Title + strconv.Itoa(i) + ".txt"
		files = append(files, path)
		go saveToFile(chapters[start:end], path, ch)
	}
	for range ch {
		e := bar.Add(1)
		if e != nil {
			return e
		}
		if bar.State().CurrentPercent == 1 {
			close(ch)
		}
	}

	file, _ := os.Create(path + "/" + book.Title + ".txt")
	writer := bufio.NewWriter(file)
	for _,filePath := range files {
		file, e := os.Open(filePath)
		if e != nil {
			return e
		}
		reader := bufio.NewReader(file)
		_, _ = io.Copy(writer, reader)
	}

	e = os.RemoveAll(tempDirPath)
	return e
}

func saveToFile(chapters []Chapter, path string, done chan int)  {
	file, e := os.Create(path)
	bufferedWriter := bufio.NewWriter(file)
	for _,chapter := range chapters {
		content := content(chapter.Link)
		content = fmt.Sprintf("%s\n%s\n",chapter.Title,content)
		_, e := bufferedWriter.WriteString(content)
		if e != nil {
			fmt.Printf("error= %s %s",e,chapter.Title)
		}
		done <- 1
	}
	e = bufferedWriter.Flush()
	if e != nil {
		fmt.Printf("error= %s %s",e,chapters)
	}
}

func content(link string) string {
	queryUrl := "http://chapterup.zhuishushenqi.com/chapter/"+url.QueryEscape(link)
	resp, err := http.Get(queryUrl)
	bytes, err := ioutil.ReadAll(resp.Body)

	content := gjson.GetBytes(bytes, "chapter.cpContent").Str
	if len(content) == 0 {
		content = gjson.GetBytes(bytes, "chapter.body").Str
	}
	if err != nil {
		fmt.Printf("error= %s %s",err,link)
	}
	return content
}

func SearchChapters(atoc Atoc) ([]Chapter, error){
	queryUrl := "http://api.zhuishushenqi.com/atoc/"+atoc.ID+"?view=chapters"
	resp, err := http.Get(queryUrl)

	bytes, err := ioutil.ReadAll(resp.Body)
	var data Data
	err = json.Unmarshal(bytes, &data)

	return data.Chapters,err
}
