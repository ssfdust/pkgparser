package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"github.com/bbrks/wrap"
	"github.com/weldpua2008/go-dialog"
	"fmt"
	"github.com/atotto/clipboard"
)

func main() {
	selected := make([]string, 0)
	// init the dialog
	window := dialog.New(dialog.CONSOLE, 0)
	window.SetTitle("Selection")
	window.SetBackTitle("Select Packages")

	// parse data from stdin
	data := read_stdin()
	if len(data) == 0 {
		os.Exit(2)
	}
	choices := parse_packages(data)
	res, err := window.Checklist(0, choices...)
	if err != nil {
		log.Println(err)
	}
	for _, pkg := range res {
		if strings.Contains(pkg, "/"){
			name := strings.Split(pkg, "/")[1]
			selected = append(selected, name)
		}
	}
	packages_str := strings.Join(selected, " ")
	fmt.Println("selected packages: " + packages_str)
	clipboard.WriteAll(packages_str)
}

func read_stdin() string {
	data := ""
	r := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 4*1024)
	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		data += string(buf)
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		// process buf
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}

	return data
}

func parse_packages(info string) []string {
	reg := `\n([a-zA-Z0-9]+/.*)\n    `

	nodes := SplitAfter(reg, info)
	return Extract(nodes)
}
func SplitAfter(reg string, str string) []string {
	var (
		r []string
		p int
	)
	re := regexp.MustCompilePOSIX(reg)
	is := re.FindAllStringIndex(str, -1)
	if is == nil {
		return append(r, str)
	}
	for _, i := range is {
		r = append(r, str[p:i[0]])
		r = append(r, str[i[0]:i[1]])
		p = i[1]
	}
	return append(r, str[p:])
}

func Extract(data []string) []string {
	ret := make([]string, 0)
	chunk := make([]string, 0)
	textwrap := wrap.NewWrapper()
	for idx, item := range data {
		if idx == 0 {
			// the first is not splited
			// we split it here
			d := strings.SplitN(item, "\n", -1)
			for _, i :=  range d {
				chunk = append(chunk, i)
		 	}
		} else if idx % 2 == 1 {
			// put the package name to list
			chunk[0] = item
			continue
		} else {
			// put the package desc to list
			chunk[1] = item
		}

		pkg := strings.TrimFunc(chunk[0], stripe)
		desc := strings.TrimFunc(chunk[1], stripe)
		if len(desc) > 78{
			desc = textwrap.Wrap(desc, 78)
		}
		for ix, d := range strings.Split(desc, "\n") {
			d = strings.TrimSpace(d)
			if ix == 0 {
				ret = append(ret, pkg, d, "off")
			} else {
				if d != ""{
					ret = append(ret, "", d, "off")
				}
			}
		}
	}
	return ret
}

func stripe(ch rune) bool {
	if ch == '\n' || ch == ' ' {
		return true
	} else {
		return false
	}
}
