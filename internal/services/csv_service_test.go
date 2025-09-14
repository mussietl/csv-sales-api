package services

import (
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCSVServiceProcessSalesCSV(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Suppress logs during testing

	csvService := NewCSVService(logger)

	tests := []struct {
		name           string
		csvContent     string
		expectedResult []DepartmentSummary
		expectError    bool
	}{
		{
			name: "valid CSV with department and sales columns",
			csvContent: `department,sales
							Electronics,1000
							Clothing,500
							Electronics,1500
							Books,300
							Clothing,200`,
			expectedResult: []DepartmentSummary{
				{Department: "Electronics", TotalSales: 2500},
				{Department: "Clothing", TotalSales: 700},
				{Department: "Books", TotalSales: 300},
			},
			expectError: false,
		},
		{
			name: "CSV with different column names",
			csvContent: `dept,amount
Electronics,1000
Clothing,500`,
			expectedResult: []DepartmentSummary{
				{Department: "Electronics", TotalSales: 1000},
				{Department: "Clothing", TotalSales: 500},
			},
			expectError: false,
		},
		{
			name: "CSV with empty department",
			csvContent: `department,sales
Electronics,1000
,500
Books,300`,
			expectedResult: []DepartmentSummary{
				{Department: "Electronics", TotalSales: 1000},
				{Department: "Books", TotalSales: 300},
			},
			expectError: false,
		},
		{
			name: "CSV with invalid sales values",
			csvContent: `department,sales
Electronics,1000
Clothing,invalid
Books,300`,
			expectedResult: []DepartmentSummary{
				{Department: "Electronics", TotalSales: 1000},
				{Department: "Books", TotalSales: 300},
			},
			expectError: false,
		},
		{
			name:           "CSV with only header",
			csvContent:     `department,sales`,
			expectedResult: []DepartmentSummary{},
			expectError:    true,
		},
		{
			name:           "CSV with empty data rows",
			csvContent:     `department,sales\n,1000\nElectronics,\n,`,
			expectedResult: []DepartmentSummary{},
			expectError:    true,
		},
		{
			name: "CSV without required columns",
			csvContent: `name,age
John,25
Jane,30`,
			expectedResult: nil,
			expectError:    true,
		},
		{
			name: "Large CSV with many rows",
			csvContent: func() string {
				// Generate a large CSV with 1000 rows
				content := "department,sales\n"
				for i := 0; i < 1000; i++ {
					dept := "Department" + string(rune('A'+(i%5))) // 5 different departments
					sales := (i + 1) * 10                          // Simple incrementing sales values
					content += fmt.Sprintf("%s,%d\n", dept, sales)
				}
				return content
			}(),
			expectedResult: []DepartmentSummary{
				{Department: "DepartmentA", TotalSales: 997000},  // 200 rows: 10+60+110+160+...+9950 = 200*4985 = 997000
				{Department: "DepartmentB", TotalSales: 999000},  // 200 rows: 20+70+120+170+...+9960 = 200*4995 = 999000
				{Department: "DepartmentC", TotalSales: 1001000}, // 200 rows: 30+80+130+180+...+9970 = 200*5005 = 1001000
				{Department: "DepartmentD", TotalSales: 1003000}, // 200 rows: 40+90+140+190+...+9980 = 200*5015 = 1003000
				{Department: "DepartmentE", TotalSales: 1005000}, // 200 rows: 50+100+150+200+...+9990 = 200*5025 = 1005000
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary CSV file
			tempFile, err := os.CreateTemp("", "test_*.csv")
			require.NoError(t, err)
			defer os.Remove(tempFile.Name())

			_, err = tempFile.WriteString(tt.csvContent)
			require.NoError(t, err)
			tempFile.Close()

			// Process the CSV
			result, err := csvService.ProcessSalesCSV(tempFile.Name())

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedResult), len(result))

				// Convert results to map for easier comparison
				resultMap := make(map[string]int)
				for _, r := range result {
					resultMap[r.Department] = r.TotalSales
				}

				expectedMap := make(map[string]int)
				for _, r := range tt.expectedResult {
					expectedMap[r.Department] = r.TotalSales
				}

				assert.Equal(t, expectedMap, resultMap)
			}
		})
	}
}

func TestCSVServiceFindColumnIndices(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	csvService := NewCSVService(logger)

	tests := []struct {
		name               string
		header             []string
		expectedDeptIndex  int
		expectedSalesIndex int
		expectError        bool
	}{
		{
			name:               "standard headers",
			header:             []string{"department", "sales"},
			expectedDeptIndex:  0,
			expectedSalesIndex: 1,
			expectError:        false,
		},
		{
			name:               "different case headers",
			header:             []string{"Department", "Sales"},
			expectedDeptIndex:  0,
			expectedSalesIndex: 1,
			expectError:        false,
		},
		{
			name:               "alternative headers",
			header:             []string{"dept", "amount"},
			expectedDeptIndex:  0,
			expectedSalesIndex: 1,
			expectError:        false,
		},
		{
			name:               "headers with extra columns",
			header:             []string{"id", "department", "sales", "date"},
			expectedDeptIndex:  1,
			expectedSalesIndex: 2,
			expectError:        false,
		},
		{
			name:        "missing department column",
			header:      []string{"name", "sales"},
			expectError: true,
		},
		{
			name:        "missing sales column",
			header:      []string{"department", "name"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deptIndex, salesIndex, err := csvService.findColumnIndices(tt.header)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedDeptIndex, deptIndex)
				assert.Equal(t, tt.expectedSalesIndex, salesIndex)
			}
		})
	}
}
