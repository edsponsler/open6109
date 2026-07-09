package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Paragraph struct {
	Index int    `json:"index"`
	Body  string `json:"body"`
}

type Chapter struct {
	Title      string      `json:"title"`
	Order      int         `json:"order"`
	Paragraphs []Paragraph `json:"paragraphs"`
}

type Book struct {
	ID       string    `json:"id"`
	Title    string    `json:"title,omitempty"`
	Author   string    `json:"author,omitempty"`
	Chapters []Chapter `json:"chapters"`
}

func ParseGutenberg(filePath string, bookID string) (*Book, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	book := &Book{ID: bookID}
	scanner := bufio.NewScanner(file)

	// Regex for identifying chapters
	chapterRegex := regexp.MustCompile(`(?i)^CHAPTER\s+([IVXLCDM\d]+)`)

	var currentChapter *Chapter
	var currentParagraph strings.Builder

	inHeader := true
	chapterCounter := 0
	paragraphCounter := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Phase 1: Handle Gutenberg boundaries
		if inHeader {
			if strings.Contains(line, "*** START OF") {
				inHeader = false
			}
			// Optional: Extract metadata from the header text (Title:, Author:) here
			if strings.HasPrefix(line, "Title:") {
				book.Title = strings.TrimSpace(strings.TrimPrefix(line, "Title:"))
			}
			if strings.HasPrefix(line, "Author:") {
				book.Author = strings.TrimSpace(strings.TrimPrefix(line, "Author:"))
			}
			continue
		}
		if strings.Contains(line, "*** END OF") {
			break
		}

		// Phase 2 & 3: Parse Chapters and Paragraphs
		if line == "" {
			// Empty line: check if we need to flush a completed paragraph
			if currentParagraph.Len() > 0 {
				if currentChapter != nil {
					currentChapter.Paragraphs = append(currentChapter.Paragraphs, Paragraph{
						Index: paragraphCounter,
						Body:  currentParagraph.String(),
					})
					paragraphCounter++
				}
				currentParagraph.Reset()
			}
			continue
		}

		// Look for a new chapter declaration
		if chapterRegex.MatchString(line) {
			// Save previous chapter if it exists
			if currentChapter != nil {
				book.Chapters = append(book.Chapters, *currentChapter)
			}
			chapterCounter++
			currentChapter = &Chapter{
				Title: line,
				Order: chapterCounter,
			}
			continue
		}

		// Reassemble wrapped text line into the paragraph buffer
		if currentParagraph.Len() > 0 {
			currentParagraph.WriteString(" ")
		}
		currentParagraph.WriteString(line)
	}

	// Flush out remaining data
	if currentParagraph.Len() > 0 && currentChapter != nil {
		currentChapter.Paragraphs = append(currentChapter.Paragraphs, Paragraph{
			Index: paragraphCounter,
			Body:  currentParagraph.String(),
		})
	}
	if currentChapter != nil {
		book.Chapters = append(book.Chapters, *currentChapter)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return book, nil
}

func main() {
	// Execute parsing
	book, err := ParseGutenberg("./corpus/pg2701.txt", "2701")
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	// Save parsed structural data to JSON file
	outFile, err := os.Create("structured_book.json")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	encoder := json.NewEncoder(outFile)
	encoder.SetIndent("", "  ") // Indented for readability
	if err := encoder.Encode(book); err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
}
