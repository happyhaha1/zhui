package main

import (
	"flag"
	"fmt"
	"github.com/happyhaha1/zhui/homdir"
	"github.com/happyhaha1/zhui/zhui"
)

func main() {
	var searchKey string
	flag.StringVar(&searchKey,"s","happyhaha","搜索")

	var downloadDir string
	homedir, e := homedir.Dir()
	if e != nil {
		fmt.Printf("error= %s",e)
	}
	flag.StringVar(&downloadDir,"d",homedir+"/Downloads","下载路径")

	flag.Parse()

	books, e := zhui.SearchBooks(searchKey)
	if e != nil {
		fmt.Printf("error= %s",e)
	}

	for index,book := range books {
		fmt.Printf("%d 书籍名称=%s 作者=%s\n",index,book.Title,book.Author)
	}
	if len(books) == 0 {
		fmt.Println("查找不到结果")
		return
	}
	fmt.Println("请输入需要下载的书籍序列号")
	var index int
	_, e = fmt.Scanln(&index)
	if e != nil {
		fmt.Printf("error= %s",e)
	}
	book := books[index]
	atocs, e := zhui.SearchAtocs(book)
	for index,atoc := range atocs {
		fmt.Printf("%d 源=%s 最后一章=%s\n",index,atoc.Name,atoc.LastChapter)
	}
	fmt.Println("请输入需要下载的源序列号")
	_, e = fmt.Scanln(&index)
	if e != nil {
		fmt.Printf("error= %s",e)
	}
	atoc := atocs[index]

	e = zhui.Download(book, atoc, downloadDir)
	if e != nil {
		fmt.Printf("error= %s",e)
	}
}
