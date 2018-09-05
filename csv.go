package main

import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
)

// CSV Format
// IMG | ID | TOP | RIGHT | BOTTOM | LEFT | USER

var images map[string]bool
var file *os.File
var csvw *csv.Writer

func openCSV(filename string) (err error) {
	// First, load all images into a map
	// (map will be better than an array for performance reasons)

	files, err := ioutil.ReadDir("images")
	if err != nil {
		return err
	}

	images = make(map[string]bool, len(files))

	for _, v := range files {
		images[v.Name()] = true
	}

	file, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return
	}

	r := csv.NewReader(file)
	_, err = r.Read() // First line should define the format
	if err == io.EOF {
		csvw = csv.NewWriter(file)
		err = csvw.Write([]string{
			"IMAGE NAME", "METEOR ID",
			"TOP", "RIGHT", "BOTTOM", "LEFT",
			"USER",
		})

		if err != nil {
			file.Close()
			return
		}

		csvw.Flush()
		if err = csvw.Error(); err != nil {
			file.Close()
			return
		}

		return nil
	}

	records, err := r.ReadAll()
	if err == io.EOF {
		return nil
	}

	if err != nil {
		file.Close()
		return
	}

	for _, rec := range records {
		if len(rec) > 0 {
			delete(images, rec[0])
		}
	}

	csvw = csv.NewWriter(file)

	return nil
}

func getImage() string {
	keys := make([]int, len(images))

	i := 0
	for k := range images {
		keys[i] = k
		i++
	}

	i = rand.Intn(len(keys))
	for ; !images[keys[i]]; i = rand.Intn(len(keys)) {
	}

	images[keys[i]] = false

	return keys[i]
}
