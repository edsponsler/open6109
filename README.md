# open6109

open6109 is a Go project designed to parse and structure text from Project Gutenberg books. It takes raw text files and converts them into structured JSON containing metadata (Book, Chapter, Paragraph) for further processing or indexing.

## Features
- Parses Project Gutenberg raw text files.
- Automatically handles Gutenberg header and footer boundaries.
- Extracts Title and Author metadata.
- Identifies and structures chapters and paragraphs.
- Outputs clean, structured JSON.

## Usage

1. Place your Project Gutenberg text files in the `corpus/` directory.
2. Update the `ParseGutenberg` function call in `main.go` with the path to your file and its ID.
3. Run the program:
   ```bash
   go run main.go
   ```
4. Check the generated `structured_book.json` file for the output.

## License
MIT
