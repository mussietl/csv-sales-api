package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mussietl/csv-sales-api/internal/models"
	"github.com/mussietl/csv-sales-api/internal/services"
	"github.com/sirupsen/logrus"
)

// UploadHandler handles file upload requests
type UploadHandler struct {
	fileService *services.FileService
	csvService  *services.CSVService
	logger      *logrus.Logger
}

// NewUploadHandler creates a new UploadHandler instance
func NewUploadHandler(fileService *services.FileService, csvService *services.CSVService, logger *logrus.Logger) *UploadHandler {
	return &UploadHandler{
		fileService: fileService,
		csvService:  csvService,
		logger:      logger,
	}
}

// UploadCSV handles CSV file upload and processing
func (h *UploadHandler) UploadCSV(c *gin.Context) {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		h.logger.Errorf("Failed to get uploaded file: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "No file uploaded or invalid file format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate the file
	if err := h.fileService.ValidateFile(file); err != nil {
		h.logger.Errorf("File validation failed: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Save the uploaded file
	filePath, err := h.fileService.SaveUploadedFile(file)
	if err != nil {
		h.logger.Errorf("Failed to save uploaded file: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to save uploaded file",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Process the CSV file
	departmentSummaries, err := h.csvService.ProcessSalesCSV(filePath)
	if err != nil {
		h.logger.Errorf("Failed to process CSV file: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to process CSV file: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Save the result file
	resultFilePath, err := h.fileService.SaveResultFile(departmentSummaries)
	if err != nil {
		h.logger.Errorf("Failed to save result file: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to save result file",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Generate download URL
	downloadURL := h.fileService.GetDownloadURL(resultFilePath)

	// Calculate total sales across all departments
	totalSales := 0
	for _, summary := range departmentSummaries {
		totalSales += summary.TotalSales
	}

	// Create response
	response := models.UploadResponse{
		Success:          true,
		Message:          "CSV file processed successfully",
		DownloadURL:      downloadURL,
		TotalDepartments: len(departmentSummaries),
		TotalSales:       totalSales,
		ProcessedAt:      time.Now().Format(time.RFC3339),
	}

	h.logger.Infof("CSV processing completed successfully. Result file: %s", resultFilePath)
	c.JSON(http.StatusOK, response)
}
