package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type FileResponse struct {
	Path        string
	Filename    string
	ContentType string
	Inline      bool
}

func SendFile(ctx *gin.Context, filePath string) {
	ctx.FileFromFS(filePath, http.Dir("./public"))
}

func SendFileWithName(ctx *gin.Context, filePath, filename string) {
	ctx.FileAttachment(filePath, filename)
}

func SendFileInline(ctx *gin.Context, filePath, filename string) {
	ctx.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
	ctx.File(filePath)
}

func SendFileResponse(ctx *gin.Context, response FileResponse) error {
	if _, err := os.Stat(response.Path); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return err
	}

	if response.Filename == "" {
		response.Filename = filepath.Base(response.Path)
	}

	if response.ContentType != "" {
		ctx.Header("Content-Type", response.ContentType)
	}

	if response.Inline {
		ctx.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", response.Filename))
	} else {
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", response.Filename))
	}

	ctx.File(response.Path)
	return nil
}

func SendBytesAsFile(ctx *gin.Context, data []byte, filename, contentType string) {
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	ctx.Data(http.StatusOK, contentType, data)
}

func SendReaderAsFile(ctx *gin.Context, reader io.Reader, filename, contentType string) error {
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	ctx.Header("Content-Type", contentType)

	_, err := io.Copy(ctx.Writer, reader)
	return err
}

func StreamFile(ctx *gin.Context, filePath string) error {
	cleanPath := filepath.Clean(filePath)
	if strings.Contains(cleanPath, "..") {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file path"})
		return fmt.Errorf("invalid file path")
	}

	file, err := os.Open(cleanPath)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log the error but do not fail the request
			fmt.Fprintf(os.Stderr, "failed to close file %s: %v\n", cleanPath, closeErr)
		}
	}()

	stat, err := file.Stat()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
		return err
	}

	ctx.Header("Content-Length", fmt.Sprintf("%d", stat.Size()))
	ctx.Header("Content-Type", getContentType(filePath))

	http.ServeContent(ctx.Writer, ctx.Request, stat.Name(), stat.ModTime(), file)
	return nil
}

func getContentType(filePath string) string {
	ext := filepath.Ext(filePath)
	contentTypes := map[string]string{
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".gif":  "image/gif",
		".svg":  "image/svg+xml",
		".zip":  "application/zip",
		".txt":  "text/plain",
		".csv":  "text/csv",
		".json": "application/json",
		".xml":  "application/xml",
	}

	if contentType, ok := contentTypes[ext]; ok {
		return contentType
	}
	return "application/octet-stream"
}

func DownloadUserResume(ctx *gin.Context) {
	userID := ctx.Param("id")
	resumePath := fmt.Sprintf("./uploads/resumes/%s.pdf", userID)

	if err := SendFileResponse(ctx, FileResponse{
		Path:        resumePath,
		Filename:    fmt.Sprintf("resume_%s.pdf", userID),
		ContentType: "application/pdf",
		Inline:      false,
	}); err != nil {
		return
	}
}

func ViewUserResume(ctx *gin.Context) {
	userID := ctx.Param("id")
	resumePath := fmt.Sprintf("./uploads/resumes/%s.pdf", userID)

	if err := SendFileResponse(ctx, FileResponse{
		Path:        resumePath,
		Filename:    fmt.Sprintf("resume_%s.pdf", userID),
		ContentType: "application/pdf",
		Inline:      true,
	}); err != nil {
		return
	}
}

func ExportJobApplicationsCSV(ctx *gin.Context) {
	csvData := []byte("Name,Email,Status\nJohn Doe,john@example.com,Applied\n")

	SendBytesAsFile(ctx, csvData, "applications.csv", "text/csv")
}
