package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
)

var url_prefix string = "http://stackoverflow.com/"

func isSummary(t html.Token) (ok bool){
	for _, attr := range t.Attr {
		if (attr.Key == "class") && (attr.Val == "summary"){
			return true
		}
	}
	return false
}

func getQ(tknzer html.Tokenizer, ch chan string) {
		tknzer.Next()
		tknzer.Next()
		tknzer.Next()
		tknzer.Next()
		ch <- string(tknzer.Text())
}

func printQ(filename string, ch chan string) {
	file,_ := os.Create("./"+filename)
	for{
		value := <-ch
		if value == "END!" {
			break
		} else {
			file.Write([]byte(value+"\n"))
		}
	}
	fmt.Println("file: "+filename+" DONE")
	file.Close()
}

// crawl the page
func Crawl (url string, ch chan string) {
	resp, _ := http.Get(url_prefix + url)
    tokenizer := html.NewTokenizer(resp.Body)
	defer resp.Body.Close()

	for {
		token := tokenizer.Next()
	    switch {
			case token == html.ErrorToken:
				// End of page	
				ch<- "END!"
				return
			case token == html.StartTagToken:
				start_tt := tokenizer.Token()
				if start_tt.Data == "div" {
					//fmt.Println("get a div! %v", num)
					if isSummary(start_tt) {
						getQ(*tokenizer,ch)
					}
				} else{
					continue
				}
		}
	}
}

func main() {
	printCh1 := make(chan string)
	printCh2 := make(chan string)
	printCh3 := make(chan string)
	defer close(printCh1)
	defer close(printCh2)
	defer close(printCh3)
	go Crawl("questions/tagged/go",printCh1)
	go Crawl("questions/tagged/java",printCh2)
	go Crawl("questions/tagged/python",printCh3)
	printQ("GO",printCh1)
	printQ("JAVA",printCh2)
	printQ("PYTHON",printCh3)
}
