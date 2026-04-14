package services

import (
	"bytes"
	"fmt"
	"time"

	"gatorpay-backend/models"

	"gorm.io/gorm"
)

type StatementService struct {
	db *gorm.DB
}

func NewStatementService(db *gorm.DB) *StatementService {
	return &StatementService{db: db}
}

func (s *StatementService) GenerateStatement(userID string, format string, startDate string, endDate string) ([]byte, string, error) {
	start, err1 := time.Parse("2006-01-02", startDate)
	end, err2 := time.Parse("2006-01-02", endDate)
	
	var txns []models.Transaction
	query := s.db.Where("from_user_id = ? OR to_user_id = ?", userID, userID)
	
	if err1 == nil && err2 == nil {
		query = query.Where("created_at BETWEEN ? AND ?", start, end.Add(24*time.Hour))
	}
	query.Order("created_at DESC").Find(&txns)

	if format == "csv" {
		return generateCSV(txns), "text/csv; charset=utf-8", nil
	} else if format == "pdf" {
		return generatePDF(txns), "application/pdf", nil
	}

	return nil, "", fmt.Errorf("unsupported format")
}

func generateCSV(txns []models.Transaction) []byte {
	var buf bytes.Buffer
	buf.WriteString("ID,Amount,Type,Description,Status,Date\n")
	for _, t := range txns {
		val, _ := t.Amount.Float64()
		row := fmt.Sprintf("%s,%.2f,%s,%s,%s,%s\n",
			t.ID, val, t.Type, t.Description, t.Status, t.CreatedAt.Format("2006-01-02 15:04:05"))
		buf.WriteString(row)
	}
	return buf.Bytes()
}

func generatePDF(txns []models.Transaction) []byte {
	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")
	buf.WriteString("=====================================\n")
	buf.WriteString("        ACCOUNT STATEMENT            \n")
	buf.WriteString("=====================================\n\n")
	
	totalIn := 0.0
	totalOut := 0.0

	for _, t := range txns {
		val, _ := t.Amount.Float64()
		if t.Type == "DEPOSIT" || t.Type == "CREDIT" {
			totalIn += val
		} else {
			totalOut += val
		}
	}
	
	buf.WriteString(fmt.Sprintf("Total In: $%.2f\nTotal Out: $%.2f\n\n", totalIn, totalOut))
	buf.WriteString("Transactions:\n--------------------------\n")
	for _, t := range txns {
		val, _ := t.Amount.Float64()
		buf.WriteString(fmt.Sprintf("%s | %s | %s | $%.2f | %s\n", 
			t.CreatedAt.Format("2006-01-02"), t.Type, t.Description, val, t.Status))
	}

	return buf.Bytes()
}
