# open6109

open6109 is a Go project designed to parse and structure text from Project Gutenberg books. It takes raw text files and converts them into structured JSON containing metadata (Book, Front Matter, Chapters, Back Matter, and Paragraphs) for further processing or indexing.

## Features

- Parses Project Gutenberg raw text files.
- Automatically handles Gutenberg header and footer boundaries.
- Extracts Title and Author metadata.
- Identifies and structures chapters, paragraphs, front matter (e.g., Etymology, Extracts), and back matter (e.g., Epilogue).
- Automatically filters out the Table of Contents (TOC) list using duplicate chapter title detection.
- Gracefully handles multi-line chapter titles.
- Normalizes and sanitizes ambiguous Unicode characters (like smart quotes, em dashes, and ligatures) into standard ASCII equivalents.
- **Quality Checker (QC) Utility**: Compares the structured JSON against the original raw text file to ensure complete text coverage and sequential integrity, and flags potential empty sections or title splits.
- Outputs clean, structured JSON.

## Usage

### 1. Run the Parser

1. Place your Project Gutenberg text files in the `corpus/` directory.
2. Update the `ParseGutenberg` function call in `main.go` with the path to your file and its ID.
3. Run the parser:

   ```bash
   go run main.go
   ```

4. Check the generated `structured_book.json` file for the output.

### 2. Run the Quality Checker

You can verify the structural accuracy of the generated JSON output against the original source file by running the Quality Checker utility:

```bash
go run ./qc/main.go --source ./corpus/pg2701.txt --json structured_book.json
```

The QC utility checks for:
- **Empty Sections**: Any chapters or front/back matter sections that contain no paragraphs.
- **Duplicate Sections**: Any duplicate section or chapter titles.
- **Suspected Title Splits**: Chapters with unusually short first paragraphs (useful for checking if a title split occurred).
- **Missing/Mismatched Text**: Performs a fuzzy, sequential search of the text to verify that every paragraph and title in the JSON matches the source text in order, ensuring no text was lost or mangled.

## License

MIT
