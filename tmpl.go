package main

import (
	"html/template"
	"io/ioutil"
	"path/filepath"
)

type alert struct {
	Msg   string
	Level string
}

type loginData struct {
	// Submit string
	Email  string
	Alerts []alert
	Tab    int
}

var tmpls *template.Template

func loadSrc() error {
	var t []string
	fs, err := ioutil.ReadDir("src")
	if err != nil {
		// fmt.Println("Error opening src directory")
		// fmt.Println(err.Error())
		// os.Exit(1)
		return err
	}

	for _, f := range fs {
		if filepath.Ext(f.Name()) == ".html" {
			t = append(t, filepath.Join("src", f.Name()))
		}
	}

	tmpls, err = template.ParseFiles(t...)
	if err != nil {
		// fmt.Println("Error parsing templates")
		// fmt.Println(err.Error())
		// os.Exit(1)

		return err
	}

	return nil
}
