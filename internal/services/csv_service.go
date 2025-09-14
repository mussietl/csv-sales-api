package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// CSVService handles CSV processing operations
type CSVService struct {
	logger *logrus.Logger
}

// NewCSVService creates a new CSVService instance
func NewCSVService(logger *logrus.Logger) *CSVService {
	return &CSVService{
		logger: logger,
	}
}

// ProcessSalesCSV processes a CSV file and returns aggregated sales data by department
func (cs *CSVService) ProcessSalesCSV(filePath string) ([]DepartmentSummary, error) {
	// Open the CSV file
	file, err := openFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)

	// Read header row first
	header, err := reader.Read()
	if err != nil {
		cs.logger.Errorf("Failed to read CSV header: %v", err)
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Parse header to find department and sales columns
	departmentIndex, salesIndex, err := cs.findColumnIndices(header)
	if err != nil {
		return nil, fmt.Errorf("failed to find required columns: %w", err)
	}

	// Process data rows using streaming
	departmentSales := make(map[string]int)
	rowNumber := 1 // Start from 1 since we already read the header

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			cs.logger.Errorf("Failed to read CSV record at row %d: %v", rowNumber+1, err)
			return nil, fmt.Errorf("failed to read CSV record at row %d: %w", rowNumber+1, err)
		}

		rowNumber++

		if len(record) <= departmentIndex || len(record) <= salesIndex {
			cs.logger.Warnf("Skipping row %d: insufficient columns", rowNumber)
			continue
		}

		department := strings.TrimSpace(record[departmentIndex])
		salesStr := strings.TrimSpace(record[salesIndex])

		if department == "" {
			cs.logger.Warnf("Skipping row %d: empty department", rowNumber)
			continue
		}

		sales, err := strconv.Atoi(salesStr)
		if err != nil {
			cs.logger.Warnf("Skipping row %d: invalid sales value '%s': %v", rowNumber, salesStr, err)
			continue
		}

		departmentSales[department] += sales
	}

	// Check if we processed any data
	if len(departmentSales) == 0 {
		return nil, fmt.Errorf("no valid data rows found in CSV file")
	}

	// Convert map to slice
	var summaries []DepartmentSummary
	for department, totalSales := range departmentSales {
		summaries = append(summaries, DepartmentSummary{
			Department: department,
			TotalSales: totalSales,
		})
	}

	cs.logger.Infof("Processed %d departments from CSV file", len(summaries))
	return summaries, nil
}

// findColumnIndices finds the indices of department and sales columns
func (cs *CSVService) findColumnIndices(header []string) (int, int, error) {
	var departmentIndex, salesIndex int = -1, -1

	for i, col := range header {
		colLower := strings.ToLower(strings.TrimSpace(col))

		// Look for department column
		if departmentIndex == -1 && (colLower == "department" ||
			colLower == "department name" ||
			strings.Contains(colLower, "department") ||
			colLower == "dept") {
			departmentIndex = i
		}

		// Look for sales column
		if salesIndex == -1 && (colLower == "sales" ||
			colLower == "total_sales" ||
			colLower == "total sales" ||
			colLower == "number of sales" ||
			colLower == "amount" ||
			colLower == "revenue") {
			salesIndex = i
		}
	}

	if departmentIndex == -1 {
		return -1, -1, fmt.Errorf("department column not found in CSV header")
	}
	if salesIndex == -1 {
		return -1, -1, fmt.Errorf("sales column not found in CSV header")
	}

	cs.logger.Infof("Found department column at index %d, sales column at index %d", departmentIndex, salesIndex)
	return departmentIndex, salesIndex, nil
}

// openFile opens a file for reading
func openFile(filePath string) (io.ReadCloser, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}
