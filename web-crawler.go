package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
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

func printQ(ch chan string) {
	for{
		value := <- ch
		if value == "END!" {
			break
		} else {
			fmt.Println(value)
		}
	}
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
		printCh := make(chan string)
		defer close(printCh)
		go Crawl("questions/tagged/go",printCh)
		printQ(printCh)
}
