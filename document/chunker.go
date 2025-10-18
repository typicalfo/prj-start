package document

import (
	"fmt"
	"github.com/typicalfo/prj-start/logger"
	"regexp"
	"strings"
)

type Chunk struct {
	Index    int
	Content  string
	Metadata map[string]string
}

type Chunker struct {
	maxChunkSize int
}

func NewChunker(maxChunkSize int) *Chunker {
	if maxChunkSize <= 0 {
		maxChunkSize = 1000 // default
	}
	return &Chunker{
		maxChunkSize: maxChunkSize,
	}
}

func (c *Chunker) ChunkDocument(fileInfo FileInfo) ([]Chunk, error) {
	logger.LogInfo(fmt.Sprintf("Chunking document: %s", fileInfo.RelativePath))

	content := fileInfo.Content
	ext := strings.ToLower(fileInfo.Extension)

	var chunks []Chunk
	var err error

	switch ext {
	case ".go":
		chunks, err = c.chunkGoCode(content)
	case ".md":
		chunks, err = c.chunkMarkdown(content)
	case ".sql":
		chunks, err = c.chunkSQL(content)
	case ".json", ".yaml", ".yml", ".toml":
		chunks, err = c.chunkConfig(content)
	case ".html":
		chunks, err = c.chunkHTML(content)
	default:
		chunks, err = c.chunkText(content)
	}

	if err != nil {
		return nil, fmt.Errorf("error chunking %s: %w", fileInfo.RelativePath, err)
	}

	// Add file metadata to each chunk
	for i := range chunks {
		if chunks[i].Metadata == nil {
			chunks[i].Metadata = make(map[string]string)
		}
		chunks[i].Metadata["filename"] = fileInfo.RelativePath
		chunks[i].Metadata["topic"] = fileInfo.Topic
		chunks[i].Metadata["extension"] = ext
		chunks[i].Metadata["total_chunks"] = fmt.Sprintf("%d", len(chunks))
	}

	logger.LogSuccess(fmt.Sprintf("Created %d chunks for %s", len(chunks), fileInfo.RelativePath))
	return chunks, nil
}

func (c *Chunker) chunkGoCode(content string) ([]Chunk, error) {
	var chunks []Chunk

	// Split by major Go constructs
	funcRegex := regexp.MustCompile(`(?m)^(func\s+\w+.*?{)`)
	typeRegex := regexp.MustCompile(`(?m)^(type\s+\w+\s+(struct|interface)\s*{)`)
	varRegex := regexp.MustCompile(`(?m)^(var\s+.*?\()`)
	constRegex := regexp.MustCompile(`(?m)^(const\s+.*?\()`)

	// Find all split points
	splitPoints := []int{0}

	matches := funcRegex.FindAllStringIndex(content, -1)
	for _, match := range matches {
		splitPoints = append(splitPoints, match[0])
	}

	matches = typeRegex.FindAllStringIndex(content, -1)
	for _, match := range matches {
		splitPoints = append(splitPoints, match[0])
	}

	matches = varRegex.FindAllStringIndex(content, -1)
	for _, match := range matches {
		splitPoints = append(splitPoints, match[0])
	}

	matches = constRegex.FindAllStringIndex(content, -1)
	for _, match := range matches {
		splitPoints = append(splitPoints, match[0])
	}

	// Sort and deduplicate split points
	splitPoints = c.uniqueSortedInts(splitPoints)

	// Create chunks
	for i := 0; i < len(splitPoints); i++ {
		start := splitPoints[i]
		end := len(content)
		if i+1 < len(splitPoints) {
			end = splitPoints[i+1]
		}

		chunk := strings.TrimSpace(content[start:end])
		if len(chunk) > 0 {
			chunks = append(chunks, Chunk{
				Index:   i,
				Content: chunk,
				Metadata: map[string]string{
					"chunk_type": "go_construct",
				},
			})
		}
	}

	if len(chunks) == 0 {
		// Fallback to text chunking
		return c.chunkText(content)
	}

	return chunks, nil
}

func (c *Chunker) chunkMarkdown(content string) ([]Chunk, error) {
	var chunks []Chunk

	// Split by headers
	headerRegex := regexp.MustCompile(`(?m)^(#{1,6}\s+.+)`)

	splitPoints := []int{0}
	matches := headerRegex.FindAllStringIndex(content, -1)
	for _, match := range matches {
		splitPoints = append(splitPoints, match[0])
	}

	splitPoints = c.uniqueSortedInts(splitPoints)

	for i := 0; i < len(splitPoints); i++ {
		start := splitPoints[i]
		end := len(content)
		if i+1 < len(splitPoints) {
			end = splitPoints[i+1]
		}

		chunk := strings.TrimSpace(content[start:end])
		if len(chunk) > 0 {
			chunks = append(chunks, Chunk{
				Index:   i,
				Content: chunk,
				Metadata: map[string]string{
					"chunk_type": "markdown_section",
				},
			})
		}
	}

	if len(chunks) == 0 {
		return c.chunkText(content)
	}

	return chunks, nil
}

func (c *Chunker) chunkSQL(content string) ([]Chunk, error) {
	var chunks []Chunk

	// Split by semicolons
	statements := strings.Split(content, ";")

	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if len(stmt) > 0 {
			chunks = append(chunks, Chunk{
				Index:   i,
				Content: stmt + ";",
				Metadata: map[string]string{
					"chunk_type": "sql_statement",
				},
			})
		}
	}

	if len(chunks) == 0 {
		return c.chunkText(content)
	}

	return chunks, nil
}

func (c *Chunker) chunkConfig(content string) ([]Chunk, error) {
	// For config files, try to split by top-level objects/sections
	var chunks []Chunk

	// Simple approach: split by double newlines for major sections
	sections := regexp.MustCompile(`\n\s*\n`).Split(content, -1)

	for i, section := range sections {
		section = strings.TrimSpace(section)
		if len(section) > 0 {
			chunks = append(chunks, Chunk{
				Index:   i,
				Content: section,
				Metadata: map[string]string{
					"chunk_type": "config_section",
				},
			})
		}
	}

	if len(chunks) == 0 {
		return c.chunkText(content)
	}

	return chunks, nil
}

func (c *Chunker) chunkHTML(content string) ([]Chunk, error) {
	var chunks []Chunk

	// Split by major HTML tags
	tagRegex := regexp.MustCompile(`(?i)<(div|section|article|header|footer|nav|main)[^>]*>`)

	splitPoints := []int{0}
	matches := tagRegex.FindAllStringIndex(content, -1)
	for _, match := range matches {
		splitPoints = append(splitPoints, match[0])
	}

	splitPoints = c.uniqueSortedInts(splitPoints)

	for i := 0; i < len(splitPoints); i++ {
		start := splitPoints[i]
		end := len(content)
		if i+1 < len(splitPoints) {
			end = splitPoints[i+1]
		}

		chunk := strings.TrimSpace(content[start:end])
		if len(chunk) > 0 {
			chunks = append(chunks, Chunk{
				Index:   i,
				Content: chunk,
				Metadata: map[string]string{
					"chunk_type": "html_section",
				},
			})
		}
	}

	if len(chunks) == 0 {
		return c.chunkText(content)
	}

	return chunks, nil
}

func (c *Chunker) chunkText(content string) ([]Chunk, error) {
	var chunks []Chunk

	// Split by paragraphs
	paragraphs := regexp.MustCompile(`\n\s*\n`).Split(content, -1)

	var currentChunk strings.Builder
	currentLength := 0
	chunkIndex := 0

	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if len(paragraph) == 0 {
			continue
		}

		if currentLength+len(paragraph) > c.maxChunkSize && currentChunk.Len() > 0 {
			// Save current chunk and start new one
			chunks = append(chunks, Chunk{
				Index:   chunkIndex,
				Content: currentChunk.String(),
				Metadata: map[string]string{
					"chunk_type": "text_paragraph",
				},
			})
			currentChunk.Reset()
			currentLength = 0
			chunkIndex++
		}

		if currentChunk.Len() > 0 {
			currentChunk.WriteString("\n\n")
		}
		currentChunk.WriteString(paragraph)
		currentLength += len(paragraph)
	}

	// Add final chunk
	if currentChunk.Len() > 0 {
		chunks = append(chunks, Chunk{
			Index:   chunkIndex,
			Content: currentChunk.String(),
			Metadata: map[string]string{
				"chunk_type": "text_paragraph",
			},
		})
	}

	return chunks, nil
}

func (c *Chunker) uniqueSortedInts(ints []int) []int {
	seen := make(map[int]bool)
	var result []int

	for _, i := range ints {
		if !seen[i] {
			seen[i] = true
			result = append(result, i)
		}
	}

	// Simple bubble sort for small slices
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i] > result[j] {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}
