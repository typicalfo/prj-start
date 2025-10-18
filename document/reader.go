package document

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"prj-start/logger"
	"strings"
)

type FileInfo struct {
	Path         string
	RelativePath string
	Topic        string // folder name
	Extension    string
	Content      string
	Size         int64
}

type Reader struct {
	rootDir string
}

func NewReader(rootDir string) *Reader {
	return &Reader{
		rootDir: rootDir,
	}
}

func (r *Reader) ReadAllDocuments() ([]FileInfo, error) {
	logger.LogInfo(fmt.Sprintf("Reading documents from: %s", r.rootDir))

	var documents []FileInfo
	err := filepath.WalkDir(r.rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.LogError(fmt.Sprintf("Error accessing path %s: %v", path, err))
			return nil // Continue walking
		}

		if d.IsDir() {
			return nil
		}

		// Skip certain files
		if r.shouldSkipFile(path) {
			return nil
		}

		fileInfo, err := r.readFile(path)
		if err != nil {
			logger.LogError(fmt.Sprintf("Error reading file %s: %v", path, err))
			return nil // Continue walking
		}

		documents = append(documents, fileInfo)
		logger.LogInfo(fmt.Sprintf("Read file: %s", fileInfo.RelativePath))
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	logger.LogSuccess(fmt.Sprintf("Successfully read %d documents", len(documents)))
	return documents, nil
}

func (r *Reader) shouldSkipFile(path string) bool {
	// Skip binary files, large files, and common non-text files
	ext := strings.ToLower(filepath.Ext(path))
	skipExts := map[string]bool{
		".png":   true,
		".jpg":   true,
		".jpeg":  true,
		".gif":   true,
		".ico":   true,
		".pdf":   true,
		".zip":   true,
		".tar":   true,
		".gz":    true,
		".exe":   true,
		".dll":   true,
		".so":    true,
		".dylib": true,
	}

	// Skip files larger than 1MB
	info, err := os.Stat(path)
	if err != nil || info.Size() > 1024*1024 {
		return true
	}

	return skipExts[ext]
}

func (r *Reader) readFile(path string) (FileInfo, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return FileInfo{}, err
	}

	relativePath, err := filepath.Rel(r.rootDir, path)
	if err != nil {
		return FileInfo{}, err
	}

	// Extract topic from the first directory level
	parts := strings.Split(relativePath, string(filepath.Separator))
	topic := "root"
	if len(parts) > 1 {
		topic = parts[0]
	}

	// Read file content
	content, err := r.readTextFile(path)
	if err != nil {
		return FileInfo{}, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return FileInfo{}, err
	}

	return FileInfo{
		Path:         absPath,
		RelativePath: relativePath,
		Topic:        topic,
		Extension:    filepath.Ext(path),
		Content:      content,
		Size:         info.Size(),
	}, nil
}

func (r *Reader) readTextFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var builder strings.Builder
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		builder.WriteString(scanner.Text())
		builder.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return builder.String(), nil
}
