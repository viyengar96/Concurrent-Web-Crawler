// Copyright © 2016 The Go Programming Language
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Findlinks3 crawls the web, starting with the URLs on the command line.
package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"time"
	"sync"

	"lab1/links"
)

type SafeMap struct {
	sync.Mutex
	URLs map[string]bool
}

var seen = SafeMap{URLs:map[string]bool{}}
var wg = sync.WaitGroup{}
var numIter = 1;

func breadthFirstMulti(url string, depth int){
	if(depth < 0){
		wg.Done()
		return;
	}

    seen.Lock()
	if seen.URLs[url] { 
		seen.Unlock()
		wg.Done()
        return
	}
	fmt.Println(url);
    seen.URLs[url] = true
    seen.Unlock()
	
	urls,err := links.Extract(url)

    if err != nil {
		fmt.Println(err)
		wg.Done()
        return
    }

    for _, u := range urls {
		wg.Add(1);
        go func(url string) {
            breadthFirstMulti(url, depth - 1)
        }(u)
	}
	wg.Done();
    return
}

//!+breadthFirst
// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item.
func breadthFirstSeq(f func(item string) []string, url string, depth int) {
	seen := make(map[string]bool)
	worklist := []string {url}
	for len(worklist) > 0 {
		if(depth < 0){
			fmt.Print("Done\n")
			return;
		}
		items := worklist
		depth -= 1;
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		fmt.Println(err)
	}
	return list
}

//!-crawl
func getDepth(input string) int{
	if !strings.Contains(input, "-depth="){
		fmt.Println("Usage: go run findlinks <url> -depth=N")
		os.Exit(1);
	}
	depthIndex := strings.Index(input, "=") + 1;
	depth, err := strconv.Atoi(string(input[depthIndex:]));
	if(err != nil){
		fmt.Println("Bad value for depth")
		fmt.Println("Usage: go run findlinks <url> -depth=N")
		os.Exit(1);
	}
	return depth;
}

//!+main
func main() {
	// Crawl the web breadth-first,
	// starting from the command-line arguments.
	if(len(os.Args) != 3){
		fmt.Println("Usage: go run findlinks <url> -depth=N")
		os.Exit(1);
	}

	avgRuntimeSeq := benchmarkSequential();
	fmt.Print("-------------------------------\n--------------------------------\n")
	avgRuntimeMulti := benchmarkMulti();
	fmt.Print("Time Elapsed with sequential crawler      : ", avgRuntimeSeq, "\n");
	fmt.Print("Time Elapsed with \"distributed\" crawler : ", avgRuntimeMulti, "\n");

}

func benchmarkSequential() float64{
	avgRuntime := 0.0;
	for i := 0; i < numIter; i++{
		start := time.Now();
		breadthFirstSeq(crawl, os.Args[1], getDepth(os.Args[2]))
		avgRuntime += float64(time.Since(start).Seconds());
	}
	avgRuntime /= float64(numIter);
	return avgRuntime;
}

func benchmarkMulti() float64{
	avgRuntime := 0.0;
	for i := 0; i< numIter; i++{
		start := time.Now();
		wg.Add(1);
		breadthFirstMulti(os.Args[1], getDepth(os.Args[2]))
		wg.Wait()
		avgRuntime += float64(time.Since(start).Seconds());
	}
	avgRuntime /= float64(numIter);
	return avgRuntime;
}


