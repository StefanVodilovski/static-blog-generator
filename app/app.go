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

	// add blog title
	blogTitle := []byte("<h1>" + format.title + "</h1>")
	htmlDocument = append(blogTitle, []byte(htmlDocument)...)

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
