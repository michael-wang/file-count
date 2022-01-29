package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	dir := parseArgs()
	counts := make(map[string]int)
	listDir(dir, counts)

	pp := sortMapByValue(counts)

	printResults(pp)
}

func parseArgs() string {
	if len(os.Args) < 2 {
		return "."
	}
	return os.Args[1]
}

func listDir(dir string, counts map[string]int) {
	ff, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range ff {
		if f.IsDir() {
			path := filepath.Join(dir, f.Name())
			listDir(path, counts)
			continue
		}

		ext := filepath.Ext(f.Name())
		c, ok := counts[ext]
		if !ok {
			counts[ext] = 1
		} else {
			counts[ext] = c + 1
		}
	}
}

type Pair struct {
	Key   string
	Value int
}

type Pairs []Pair

func (pp Pairs) Swap(i, j int) {
	pp[i], pp[j] = pp[j], pp[i]
}
func (pp Pairs) Len() int {
	return len(pp)
}
func (pp Pairs) Less(i, j int) bool {
	return pp[i].Value < pp[j].Value
}

func sortMapByValue(m map[string]int) Pairs {
	pp := make(Pairs, len(m))
	i := 0
	for k, v := range m {
		pp[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pp))
	return pp
}

func printResults(pp Pairs) {
	ptr := message.NewPrinter(language.English)

	fmt.Printf("# of files\textension\n")
	for _, p := range pp {
		ptr.Printf("%8d\t%s\n", p.Value, p.Key)
	}
}
