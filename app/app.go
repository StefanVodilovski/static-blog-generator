package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/russross/blackfriday/v2"
	"github.com/urfave/cli/v2"
)

func generate(format outputFormat) error {
	// Read all markdown files from input folder
	// files, err := os.ReadDir(format.inputFolder)
	files, err := getMarkdownFilesWithDates(format.inputFolder)
	if err != nil {
		return err
	}
	//sort them
	sort.Slice(files, func(i, j int) bool {
		return files[i].Date < files[j].Date
	})

	// check if the output folder is empty
	isEmpty, err := checkOutput(format.outputFolder)
	if err != nil {
		return err
	}
	if !isEmpty {
		log.Panic("The output folder specified is not empty. Please provide an empty folder")
	}

	var htmlDocument []byte
	var postCount int
	var page int

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
			postCount++

			if format.posts != 0 && (postCount%format.posts == 0 || postCount == len(files)) {
				page += 1

				// Write HTML to output file
				err = output(htmlDocument, format, page)
				if err != nil {
					return err
				}

				//reset the HTML
				htmlDocument = []byte{}

			}

		}
	}

	//if we want every post in one HTML File
	if format.posts == 0 {
		err = output(htmlDocument, format, page)
		if err != nil {
			return err
		}
	}

	// add pagination
	if err := pagination(format); err != nil {
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
					&cli.IntFlag{
						Name:        "posts-per-page",
						Aliases:     []string{"ppp"},
						Usage:       "how many posts per page should there be",
						Destination: &outputFormat.posts,
						Required:    false,
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
