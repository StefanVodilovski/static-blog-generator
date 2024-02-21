package main

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type outputFormat struct {
	inputFolder  string
	outputFolder string
	title        string
	posts        int
}

type fileInfoWithDate struct {
	fs.DirEntry
	Date string
}

func addTextToFile(container, filePath string) error {
	existingContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Open the file for writing (truncating it in the process)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the existing content back to the file
	_, err = file.Write(existingContent)
	if err != nil {
		return err
	}

	// Append the new text to the file
	_, err = file.WriteString(container)
	if err != nil {
		return err
	}

	return nil

}

func pagination(format outputFormat) error {
	files, err := os.ReadDir(format.outputFolder)
	if err != nil {
		return err
	}

	var container string
	for i, file := range files {
		if i == 0 {
			container = "<div class=\"pagination\"> <a href=\"" + strconv.Itoa(i+1) + ".html\"> < </a>"
		} else {
			container = "<div class=\"pagination\"> <a href=\"" + strconv.Itoa(i) + ".html\"> < </a>"
		}
		for j := range files {
			var link string
			if i == j {
				link = "<a href=\"" + strconv.Itoa(j+1) + ".html\" class=\"active\">" + strconv.Itoa(j+1) + "</a>"
			} else {
				link = "<a href=\"" + strconv.Itoa(j+1) + ".html\">" + strconv.Itoa(j+1) + "</a>"
			}
			container += link

		}

		if i == len(files)-1 {
			container += "<a href=\"" + strconv.Itoa(i+1) + ".html\"> > </a>"
		} else {
			container += "<a href=\"" + strconv.Itoa(i+2) + ".html\"> > </a>"
		}
		container += "</div>"
		filePath := filepath.Join(format.outputFolder, file.Name())
		err = addTextToFile(container, filePath)
		if err != nil {
			log.Fatal()
		}
		container = ""
	}

	return nil
}

func cssString() string {
	return `<style>
	.pagination {
		text-align: center;
	}

	.pagination a {
		color: black;
		text-decoration: none;
		padding: 8px 15px;
		display: inline-block;
	}

	.pagination a.active {
		background-color: hsl(120, 100%, 70%);
		font-weight: bold;
		border-radius: 5px;
	}

	.pagination a:hover:not(.active) {
		background-color: hsl(0, 0%, 77%);
		border-radius: 5px;
	}

	html * {
		font-family: Arial, sans-serif;
	}

	hr {
		border: solid 1px #ccc;
		margin-bottom: 50px;
		margin-top: 50px;
	}

	body {
		width: 750px;
		margin-left: auto;
		margin-right: auto;
	}

	h1 {
		text-align: right;
		color: #6d4aff;
		margin-bottom: 50px;
	}

	h2 {
		text-align: center;
		color: #372580;
		margin-bottom: 50px;
	}

	p {
		text-align: justify;
	}
	</style>
`
}

func style(html []byte, title string) []byte {
	boilerPlate := "<!DOCTYPE html><html lang=\"en\"> <head><meta http-equiv=\"content-type\" content=\"text/html; charset=UTF-8\"> <meta charset=\"utf-8\">"

	blogTitle := []byte("<title>" + title + "</title>")

	styleTag := cssString()

	blogHeading := []byte("<h1>" + title + "</h1>")

	// add the tags to the document
	html = append([]byte(blogHeading), []byte(html)...)
	html = append([]byte(styleTag), []byte(html)...)
	html = append(blogTitle, []byte(html)...)
	html = append([]byte(boilerPlate), []byte(html)...)

	return html
}

func output(htmlDoc []byte, format outputFormat, page int) error {

	// add style to the html
	htmlDoc = style(htmlDoc, format.title)

	var outputFilePath string
	if format.posts != 0 {
		outputFilePath = filepath.Join(format.outputFolder, strconv.Itoa(page)+".html")
	} else {
		outputFilePath = filepath.Join(format.outputFolder, "1.html")
	}

	if err := os.WriteFile(outputFilePath, []byte(htmlDoc), 0644); err != nil {
		return err
	}

	return nil
}

func checkOutput(directory string) (bool, error) {
	// Read the contents of the directory
	files, err := os.ReadDir(directory)
	if err != nil {
		return false, err
	}

	// If there are no files or subdirectories, the directory is empty
	return len(files) == 0, nil
}

func extractDateFromMarkdown(markdownBytes []byte) (string, error) {
	//iterate through each line
	lines := strings.Split(string(markdownBytes), "\n")
	for _, line := range lines {
		if strings.Contains(line, "*Published on ") {

			// Extract the date after "Published on"
			parts := strings.Split(line, " ")
			return strings.Split(parts[2], ".")[0], nil
			// if len(parts) >= 4 {
			// 	return parts[len(parts)-1], nil
			// }
		}
	}
	errMsg := "Unable to extract date from Markdown: no line containing 'Published on' found"
	log.Println(errMsg)
	return "", errors.New(errMsg)
}

func getMarkdownFilesWithDates(inputFolder string) ([]fileInfoWithDate, error) {
	var filesWithDates []fileInfoWithDate

	files, err := os.ReadDir(inputFolder)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".md" {
			markdownBytes, err := os.ReadFile(filepath.Join(inputFolder, file.Name()))
			if err != nil {
				return nil, err
			}

			date, err := extractDateFromMarkdown(markdownBytes)
			if err != nil {
				log.Printf("Error extracting date from %s: %v", file.Name(), err)
				continue
			}

			filesWithDates = append(filesWithDates, fileInfoWithDate{file, date})
		}
	}

	return filesWithDates, nil
}
