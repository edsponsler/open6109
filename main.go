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

type Section struct {
	Title      string      `json:"title"`
	Paragraphs []Paragraph `json:"paragraphs"`
}

type Book struct {
	ID          string    `json:"id"`
	Title       string    `json:"title,omitempty"`
	Author      string    `json:"author,omitempty"`
	FrontMatter []Section `json:"front_matter,omitempty"`
	Chapters    []Chapter `json:"chapters"`
	BackMatter  []Section `json:"back_matter,omitempty"`
}

var textSanitizer = strings.NewReplacer(
	"“", "\"",
	"”", "\"",
	"‘", "'",
	"’", "'",
	"—", "--",
	"æ", "ae",
	"œ", "oe",
	"Œ", "Oe",
	"é", "e",
	"è", "e",
	"\u200e", "", // Remove invisible direction marks
	"\u200f", "",
)

type BookConfig struct {
	FrontMatterStart []string
	BackMatterStart  []string
}

var bookConfigs = map[string]BookConfig{
	"2701": {
		FrontMatterStart: []string{"ETYMOLOGY.", "EXTRACTS"},
		BackMatterStart:  []string{"Epilogue"},
	},
}

type ParseContext int

const (
	CtxNone ParseContext = iota
	CtxFrontMatter
	CtxChapter
	CtxBackMatter
)

func cleanTitle(t string) string {
	t = strings.ToLower(t)
	var sb strings.Builder
	for _, r := range t {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func ParseGutenberg(filePath string, bookID string) (*Book, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	book := &Book{ID: bookID}
	scanner := bufio.NewScanner(file)

	// Regex for identifying chapters. Using \b to ensure word boundaries
	chapterRegex := regexp.MustCompile(`(?i)^CHAPTER\s+([IVXLCDM\d]+)\b`)

	config := bookConfigs[bookID]

	var activeChapter *Chapter
	var currentParagraph strings.Builder

	inHeader := true
	chapterCounter := 0
	paragraphCounter := 0

	var currentContext ParseContext = CtxNone
	var parsingTitle bool
	var titleBuilder strings.Builder

	flushParagraph := func() {
		if currentParagraph.Len() > 0 {
			p := Paragraph{
				Index: paragraphCounter,
				Body:  currentParagraph.String(),
			}
			paragraphCounter++

			switch currentContext {
			case CtxFrontMatter:
				if len(book.FrontMatter) > 0 {
					idx := len(book.FrontMatter) - 1
					book.FrontMatter[idx].Paragraphs = append(book.FrontMatter[idx].Paragraphs, p)
				}
			case CtxChapter:
				if activeChapter != nil {
					activeChapter.Paragraphs = append(activeChapter.Paragraphs, p)
				}
			case CtxBackMatter:
				if len(book.BackMatter) > 0 {
					idx := len(book.BackMatter) - 1
					book.BackMatter[idx].Paragraphs = append(book.BackMatter[idx].Paragraphs, p)
				}
			}
			currentParagraph.Reset()
		}
	}

	flushActiveChapter := func() {
		if activeChapter != nil {
			cleanNewTitle := cleanTitle(activeChapter.Title)
			dupIndex := -1
			for i, ch := range book.Chapters {
				if cleanTitle(ch.Title) == cleanNewTitle {
					dupIndex = i
					break
				}
			}
			if dupIndex != -1 {
				// Duplicate detected (TOC entry). Discard everything up to and including it.
				book.Chapters = book.Chapters[dupIndex+1:]
			}
			book.Chapters = append(book.Chapters, *activeChapter)
			activeChapter = nil
		}
	}

	finishTitleParsing := func() {
		if !parsingTitle {
			return
		}
		parsingTitle = false
		fullTitle := strings.TrimSpace(titleBuilder.String())

		switch currentContext {
		case CtxFrontMatter:
			cleanNewTitle := cleanTitle(fullTitle)
			dupIndex := -1
			for i, s := range book.FrontMatter {
				if cleanTitle(s.Title) == cleanNewTitle {
					dupIndex = i
					break
				}
			}
			if dupIndex != -1 {
				book.FrontMatter = book.FrontMatter[dupIndex+1:]
			}
			book.FrontMatter = append(book.FrontMatter, Section{Title: fullTitle})

		case CtxChapter:
			chapterCounter++
			activeChapter = &Chapter{
				Title: fullTitle,
				Order: chapterCounter,
			}

		case CtxBackMatter:
			cleanNewTitle := cleanTitle(fullTitle)
			dupIndex := -1
			for i, s := range book.BackMatter {
				if cleanTitle(s.Title) == cleanNewTitle {
					dupIndex = i
					break
				}
			}
			if dupIndex != -1 {
				book.BackMatter = book.BackMatter[dupIndex+1:]
			}
			book.BackMatter = append(book.BackMatter, Section{Title: fullTitle})
		}
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = textSanitizer.Replace(line)

		// Phase 1: Handle Gutenberg boundaries
		if inHeader {
			if strings.Contains(line, "*** START OF") {
				inHeader = false
			}
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

		// Check for Back Matter triggers
		isBackMatter := false
		for _, trigger := range config.BackMatterStart {
			if strings.HasPrefix(line, trigger) {
				isBackMatter = true
				break
			}
		}
		if isBackMatter {
			flushParagraph()
			flushActiveChapter()
			currentContext = CtxBackMatter
			parsingTitle = true
			titleBuilder.Reset()
			titleBuilder.WriteString(line)
			continue
		}

		// Check for Front Matter triggers
		isFrontMatter := false
		for _, trigger := range config.FrontMatterStart {
			if strings.HasPrefix(line, trigger) {
				isFrontMatter = true
				break
			}
		}
		if isFrontMatter {
			flushParagraph()
			flushActiveChapter()
			currentContext = CtxFrontMatter
			parsingTitle = true
			titleBuilder.Reset()
			titleBuilder.WriteString(line)
			continue
		}

		// Look for a new chapter declaration
		if chapterRegex.MatchString(line) {
			flushParagraph()
			flushActiveChapter()
			currentContext = CtxChapter
			parsingTitle = true
			titleBuilder.Reset()
			titleBuilder.WriteString(line)
			continue
		}

		// Phase 2 & 3: Parse Titles, Chapters, and Paragraphs
		if line == "" {
			if parsingTitle {
				finishTitleParsing()
			} else {
				flushParagraph()
			}
			continue
		}

		if parsingTitle {
			titleBuilder.WriteString(" ")
			titleBuilder.WriteString(line)
			continue
		}

		// Reassemble wrapped text line into the paragraph buffer
		if currentParagraph.Len() > 0 {
			currentParagraph.WriteString(" ")
		}
		currentParagraph.WriteString(line)
	}

	// Flush remaining data
	if parsingTitle {
		finishTitleParsing()
	} else {
		flushParagraph()
	}
	flushActiveChapter()

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
