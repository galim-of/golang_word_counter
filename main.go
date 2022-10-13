package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"sync"
)

func main() {
	needle := flag.String("needle", "Go", "specify string that will be searched for in HTTP response")
	maxHandlers := flag.Int("max-handlers", 5, "maximum number of concurrenlty running handlers")
	debug := flag.Bool("debug", false, "show number of running goroutines with errors, if any")
	flag.Parse()

	urls := flag.Args()
	if len(flag.Args()) == 0 {
		log.Fatal("specify atleast one URL as: <program_name> <url> [<url>...]")
	}
	if *debug {
		log.Printf("count of urls = %d\n", len(urls))
	}

	total := 0
	waitCh := make(chan struct{}, *maxHandlers)
	wg := sync.WaitGroup{}
	for _, url := range urls {
		waitCh <- struct{}{}
		wg.Add(1)
		go func(wg *sync.WaitGroup, url, needle string) {
			<-waitCh
			if *debug {
				log.Printf("num of running goroutines %d\n", runtime.NumGoroutine())
			}
			count, err := worker(url, needle)
			if err != nil && *debug {
				log.Printf("%v\n", err)
			}
			total += count
			wg.Done()
		}(&wg, url, *needle)
	}
	wg.Wait()
	close(waitCh)
	log.Printf("Total: %d\n", total)
}

func worker(target string, needle string) (int, error) {
	url, err := url.ParseRequestURI(target)
	if err != nil {
		return 0, fmt.Errorf("wrong url found %q will ignore", target)
	}
	resp, err := http.Get(url.String())
	if err != nil {
		return 0, fmt.Errorf("can't make request to %q: %v", target, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("response status code for %q: %v will ignore", target, resp.StatusCode)
	}
	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("can't read response data for %q: %v will ignore", target, err)
	}
	count := strings.Count(string(rawData), needle)
	log.Printf("Count for %s: %d\n", target, count)
	return count, nil

}
