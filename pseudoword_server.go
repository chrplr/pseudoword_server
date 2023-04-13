package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	//"encoding/json"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

// fileCache is used by readLines to cache files in memory
var fileCache = make(map[string][]string)

// readLines returns the lines from a textfile
func readLines(path string) ([]string, error) {

	if lines, ok := fileCache[path]; ok {
		return lines, nil
	}

	// as the file is not in the cache, load it from disk
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	fileCache[path] = lines
	return lines, scanner.Err()
}

// RandSample returns 'k' random integers from 1..n
func RandomSample(n, k int) []int {
	perm := make([]int, n, n)
	for i := 0; i < n; i++ {
		perm[i] = i
	}
	if k > n {
		k = n
	}
	for i := 0; i < k; i++ {
		j := i + rand.Intn(n-i)
		perm[i], perm[j] = perm[j], perm[i]
	}
	return perm[:k]
}

// This could be read from a json configuration file
var srcFiles = map[string]string{}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}

	//var err error
	var n int64

	langl, ok := params["lang"]
	if !ok {
		fmt.Fprintf(w, "lang parameter missing")
	}
	lang := langl[0]

	ns, ok := params["n"]
	if !ok {
		n = 100
	} else {
		if n, err = strconv.ParseInt(ns[0], 10, 32); err != nil {
			fmt.Fprintf(w, "n parameter not an int")
		}
	}

	if lines, err := readLines(srcFiles[lang]); err != nil {
		panic(err)
	} else {
		var n0 = int(n)
		samp := RandomSample(len(lines), n0)
		for i, index := range samp {
			if i == 0 {
				fmt.Fprintf(w, "%s", lines[index])
			} else {
				fmt.Fprintf(w, "\t%s", lines[index])
			}
		}
	}

}

func main() {
	srcFiles = map[string]string{"en": "lists/English_1000pseudos.txt",
		"fr": "lists/French_1000pseudos.txt"}

	http.HandleFunc("/", handleQuery)

	http.HandleFunc("/help", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Usage: user ?lang=[en|fr]&n=100\n")
	})

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
