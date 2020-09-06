package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

var templates *template.Template

func main() {

	templates = template.Must(template.ParseGlob("templates/*.html"))
	http.HandleFunc("/", posthandler)
	fmt.Println("listining port:7777")
	http.ListenAndServe(":7777", nil)

}

func posthandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if r.URL.Path != "/" {
			http.Error(w, "ERROR-404\nPage not found(", http.StatusNotFound)
			return
		}
		templates.ExecuteTemplate(w, "index.html", nil)
	}

	if r.Method == "POST" {
		text := r.FormValue("text")
		font := r.FormValue("fonts")

		for _, v := range text {
			if !(v >= 32 && v <= 126) {

				http.Error(w, "ERROR-400\nBad request!", http.StatusBadRequest)
				return
			}
		}

		file, err := os.Open(FormatType(font))

		if err != nil {
			http.Error(w, "Internal Server Error!!!\nERROR-500", http.StatusInternalServerError)
			return
		}

		defer file.Close()
		banners := [][]string{}
		banner := []string{}

		scanner := bufio.NewScanner(file)
		i := 0

		for scanner.Scan() {
			i++
			banner = append(banner, scanner.Text())

			if i == 9 {
				banners = append(banners, banner)
				banner = []string{}
				i = 0
			}

		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error")
			return
		}
		c := ""
		str := []string{}
		array := strings.Split(text, "\\n")
		for i := 0; i < len(array); i++ {

			for j := 1; j <= 8; j++ {
				for _, value := range array[i] {

					str = banners[int(value)-32]

					c = c + str[j]
				}

				if len(array[i]) != 0 {
					c = c + "\n"
				}
			}
		}

		templates.ExecuteTemplate(w, "index.html", c)
	}

}

func FormatType(fs string) string {
	if fs == "shadow" {
		return "shadow.txt"
	}
	if fs == "thinkertoy" {
		return "thinkertoy.txt"
	}
	if fs == "standard" {
		return "standard.txt"
	}
	return "Error 500"
}
