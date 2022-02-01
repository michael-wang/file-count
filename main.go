package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	byExt = "group by extension"
	byDir = "group by directory"
)

type config struct {
	dir         string
	groupBy     string
	noRecursive bool
	ascending   bool
}

var cfg config

func init() {
	flag.StringVar(&cfg.dir, "d", ".", "Directory to start with")
	flag.StringVar(&cfg.groupBy, "g", "ext", "Group by file extension ('ext' or 'e') or by directory ('dir' or 'd')")
	flag.BoolVar(&cfg.noRecursive, "r", false, "Suppress recursive counting into sub-directories")
	flag.BoolVar(&cfg.ascending, "a", false, "Sorting results by ascending order (default by descending order)")
	flag.Parse()

	switch strings.ToLower(cfg.groupBy) {
	case "dir":
		fallthrough
	case "d":
		cfg.groupBy = byDir
	case "ext":
		fallthrough
	case "e":
		cfg.groupBy = byExt
	default:
		log.Fatalf("Unknown value for flag 'group-by': %s\n", cfg.groupBy)
	}
}

func main() {
	results := make(map[string]int)
	count(cfg.dir, results)

	pp := sortMapByValue(results)

	printResults(pp)
}

func count(dir string, results map[string]int) {
	ff, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range ff {
		if f.IsDir() {
			if !cfg.noRecursive {
				path := filepath.Join(dir, f.Name())
				count(path, results)
			}
			continue
		}

		if cfg.groupBy == byExt {
			ext := filepath.Ext(f.Name())
			c, ok := results[ext]
			if !ok {
				results[ext] = 1
			} else {
				results[ext] = c + 1
			}
		} else if cfg.groupBy == byDir {
			c, ok := results[dir]
			if !ok {
				results[dir] = 1
			} else {
				results[dir] = c + 1
			}
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

	if cfg.ascending {
		sort.Sort(pp)
	} else {
		sort.Sort(sort.Reverse(pp))
	}
	return pp
}

func printResults(pp Pairs) {
	ptr := message.NewPrinter(language.English)

	fmt.Printf("# of files\textension\n")
	for _, p := range pp {
		ptr.Printf("%8d\t%s\n", p.Value, p.Key)
	}
}
