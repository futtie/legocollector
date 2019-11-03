package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func getPartImageAndDescription(partnumber string, color int) (string, error) {
	descriptionURL := "https://www.bricklink.com/v2/catalog/catalogitem.page?P=" + partnumber + "#T=S&C=" + strconv.Itoa(color)
	imageURL := "http://img.bricklink.com/ItemImage/PN/" + strconv.Itoa(color) + "/" + partnumber + ".png"
	destination := localPartImageStorage + partnumber + "-" + strconv.Itoa(color) + ".png"

	// Get the data
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(destination)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	resp, err = http.Get(descriptionURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	htmlAsString := string(bytes)

	var description string
	re := regexp.MustCompile(`<title>BrickLink - Part .* : (.*) - BrickLink Reference Catalog</title>`)
	subMatches := re.FindStringSubmatch(htmlAsString)
	if len(subMatches) == 2 {
		description = subMatches[1]
	}
	return description, err
}

func saveImageByURL(ID int, imageURL string) error {
	ext := filepath.Ext(imageURL)
	destination := localSetImageStorage + strconv.Itoa(ID) + ext
	// Get the data
	resp, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
