package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/russross/blackfriday/v2"
	"github.com/urfave/cli/v2"
)

type outputFormat struct {
	inputFolder  string
	outputFolder string
	title        string
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

func generate(format outputFormat) error {
	// Read all markdown files from input folder
	files, err := os.ReadDir(format.inputFolder)
	if err != nil {
		return err
	}

	var htmlDocument []byte

	//Iterate through all files
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".md" {

			// Read the markdown file
			markdownBytes, err := os.ReadFile(filepath.Join(format.inputFolder, file.Name()))
			if err != nil {
				return err
			}

			// Convert markdown to HTML
			html := blackfriday.Run(markdownBytes)

			// Add to  the html document
			htmlDocument = append(htmlDocument, html...)

		}
	}

	// add style
	htmlDocument = style(htmlDocument, format.title)

	// Write HTML to output file
	outputFilePath := filepath.Join(format.outputFolder, "index.html")
	if err := os.WriteFile(outputFilePath, []byte(htmlDocument), 0644); err != nil {
		return err
	}

	return nil
}

func main() {

	var outputFormat outputFormat

	app := &cli.App{
		Name:  "gen-blog",
		Usage: "Static blog generator - Turn your Markdown files in to HTML files",
		Commands: []*cli.Command{
			{
				Name: "generate",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "input",
						Aliases:     []string{"i"},
						Usage:       "the path of the input folder containing the markdown files.",
						Destination: &outputFormat.inputFolder,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "output",
						Aliases:     []string{"o"},
						Usage:       "the path of the output folder where the static HMTL content will be created.",
						Destination: &outputFormat.outputFolder,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "title",
						Aliases:     []string{"t"},
						Usage:       "the title of the blog.",
						Destination: &outputFormat.title,
						Required:    true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					return generate(outputFormat)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
