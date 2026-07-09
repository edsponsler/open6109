package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
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

func normalizeText(s string) string {
	s = strings.ToLower(s)
	s = textSanitizer.Replace(s)
	var sb strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func getBodyText(srcText string) string {
	lines := strings.Split(srcText, "\n")
	var bodyBuilder strings.Builder
	inHeader := true

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if inHeader {
			if strings.Contains(line, "*** START OF") {
				inHeader = false
			}
			continue
		}
		if strings.Contains(line, "*** END OF") {
			break
		}
		bodyBuilder.WriteString(line)
		bodyBuilder.WriteString("\n")
	}
	return bodyBuilder.String()
}

func checkEmptySections(book *Book) int {
	emptyCount := 0
	for _, s := range book.FrontMatter {
		if len(s.Paragraphs) == 0 {
			fmt.Printf("[EMPTY FRONT MATTER] Section: %q has 0 paragraphs.\n", s.Title)
			emptyCount++
		}
	}
	for _, ch := range book.Chapters {
		if len(ch.Paragraphs) == 0 {
			fmt.Printf("[EMPTY CHAPTER] Chapter: %q (Order: %d) has 0 paragraphs.\n", ch.Title, ch.Order)
			emptyCount++
		}
	}
	for _, s := range book.BackMatter {
		if len(s.Paragraphs) == 0 {
			fmt.Printf("[EMPTY BACK MATTER] Section: %q has 0 paragraphs.\n", s.Title)
			emptyCount++
		}
	}
	return emptyCount
}

func checkDuplicateSections(book *Book) int {
	dupCount := 0

	// Check FrontMatter
	seenFront := make(map[string]int)
	for _, s := range book.FrontMatter {
		seenFront[s.Title]++
	}
	for title, count := range seenFront {
		if count > 1 {
			fmt.Printf("[DUPLICATE FRONT MATTER] Section: %q appears %d times.\n", title, count)
			dupCount += (count - 1)
		}
	}

	// Check Chapters
	seenChapters := make(map[string]int)
	for _, ch := range book.Chapters {
		seenChapters[ch.Title]++
	}
	for title, count := range seenChapters {
		if count > 1 {
			fmt.Printf("[DUPLICATE CHAPTER] Chapter: %q appears %d times.\n", title, count)
			dupCount += (count - 1)
		}
	}

	// Check BackMatter
	seenBack := make(map[string]int)
	for _, s := range book.BackMatter {
		seenBack[s.Title]++
	}
	for title, count := range seenBack {
		if count > 1 {
			fmt.Printf("[DUPLICATE BACK MATTER] Section: %q appears %d times.\n", title, count)
			dupCount += (count - 1)
		}
	}

	return dupCount
}

func checkTitleSplits(book *Book) int {
	splitCount := 0
	for _, ch := range book.Chapters {
		if len(ch.Paragraphs) > 0 {
			first := ch.Paragraphs[0].Body
			// If the first paragraph is extremely short, it is likely a split chapter title
			if len(first) < 80 {
				fmt.Printf("[SUSPECTED TITLE SPLIT] Chapter: %q has a very short first paragraph: %q (Length: %d)\n", ch.Title, first, len(first))
				splitCount++
			}
		}
	}
	return splitCount
}

func verifySectionCoverage(title string, paragraphs []Paragraph, normalizedSource string, currentOffset *int) int {
	mismatchCount := 0

	// Check title
	normTitle := normalizeText(title)
	idx := strings.Index(normalizedSource[*currentOffset:], normTitle)
	if idx == -1 {
		globalIdx := strings.Index(normalizedSource, normTitle)
		if globalIdx != -1 {
			fmt.Printf("[OUT OF ORDER] Section title %q found at index %d, expected after offset %d\n", title, globalIdx, *currentOffset)
		} else {
			fmt.Printf("[MISSING TITLE] Section title %q not found in source text\n", title)
		}
		mismatchCount++
	} else {
		*currentOffset += idx + len(normTitle)
	}

	// Check paragraphs
	for _, p := range paragraphs {
		normBody := normalizeText(p.Body)
		if len(normBody) == 0 {
			continue
		}

		pIdx := strings.Index(normalizedSource[*currentOffset:], normBody)
		if pIdx == -1 {
			prefixLen := 100
			if len(normBody) < prefixLen {
				prefixLen = len(normBody)
			}
			pPrefix := normBody[:prefixLen]
			pIdx2 := strings.Index(normalizedSource[*currentOffset:], pPrefix)
			if pIdx2 == -1 {
				fmt.Printf("[MISSING TEXT] Section %q, Paragraph %d not found in source text\n", title, p.Index)
				mismatchCount++
			} else {
				fmt.Printf("[PARTIAL MATCH] Section %q, Paragraph %d matched partially (prefix matched)\n", title, p.Index)
				*currentOffset += pIdx2 + len(pPrefix)
			}
		} else {
			*currentOffset += pIdx + len(normBody)
		}
	}

	return mismatchCount
}

func checkFuzzyTextCoverage(book *Book, srcText string) int {
	normalizedSource := normalizeText(getBodyText(srcText))
	currentOffset := 0
	mismatchCount := 0

	for _, s := range book.FrontMatter {
		if len(s.Paragraphs) == 0 {
			continue
		}
		mismatchCount += verifySectionCoverage(s.Title, s.Paragraphs, normalizedSource, &currentOffset)
	}
	for _, ch := range book.Chapters {
		if len(ch.Paragraphs) == 0 {
			continue
		}
		mismatchCount += verifySectionCoverage(ch.Title, ch.Paragraphs, normalizedSource, &currentOffset)
	}
	for _, s := range book.BackMatter {
		if len(s.Paragraphs) == 0 {
			continue
		}
		mismatchCount += verifySectionCoverage(s.Title, s.Paragraphs, normalizedSource, &currentOffset)
	}
	return mismatchCount
}

func main() {
	sourcePath := flag.String("source", "./corpus/pg2701.txt", "Path to original Gutenberg source file")
	jsonPath := flag.String("json", "structured_book.json", "Path to structured JSON output file")
	flag.Parse()

	// 1. Read files
	jsonBytes, err := os.ReadFile(*jsonPath)
	if err != nil {
		fmt.Printf("Error reading JSON file %q: %v\n", *jsonPath, err)
		os.Exit(1)
	}

	srcBytes, err := os.ReadFile(*sourcePath)
	if err != nil {
		fmt.Printf("Error reading source file %q: %v\n", *sourcePath, err)
		os.Exit(1)
	}

	var book Book
	if err := json.Unmarshal(jsonBytes, &book); err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== RUNNING QUALITY CHECKER ===")
	fmt.Printf("Source File: %s\n", *sourcePath)
	fmt.Printf("JSON File:   %s\n", *jsonPath)
	fmt.Println("-------------------------------")

	emptyCount := checkEmptySections(&book)
	dupCount := checkDuplicateSections(&book)
	splitCount := checkTitleSplits(&book)
	mismatchCount := checkFuzzyTextCoverage(&book, string(srcBytes))

	fmt.Println("-------------------------------")
	fmt.Println("=== QC SUMMARY ===")
	fmt.Printf("Empty Sections:         %d\n", emptyCount)
	fmt.Printf("Duplicate Sections:     %d\n", dupCount)
	fmt.Printf("Suspected Title Splits: %d\n", splitCount)
	fmt.Printf("Missing/Mismatched Text: %d\n", mismatchCount)

	if emptyCount > 0 || dupCount > 0 || mismatchCount > 0 {
		fmt.Println("STATUS: FAIL")
		os.Exit(1)
	}

	fmt.Println("STATUS: PASS")
}
