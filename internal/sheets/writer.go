package sheets

import (
	"fmt"

	"DataFPP_Dingtalk/internal/model"
	"DataFPP_Dingtalk/pkg/logger"

	"google.golang.org/api/sheets/v4"
)

const destTab = "FPPP Data"

type Writer struct {
	srv     *sheets.Service
	sheetID string
}

func NewWriter(srv *sheets.Service, sheetID string) *Writer {
	return &Writer{srv: srv, sheetID: sheetID}
}

// ReadExisting reads the destination sheet and returns a map of business_id -> row number (1-indexed).
func (w *Writer) ReadExisting() (map[string]int, error) {
	rng := fmt.Sprintf("'%s'!A:A", destTab)
	resp, err := w.srv.Spreadsheets.Values.Get(w.sheetID, rng).Do()
	if err != nil {
		return nil, fmt.Errorf("reading existing destination data: %w", err)
	}

	existing := make(map[string]int)
	for i, row := range resp.Values {
		if i == 0 {
			continue // skip header
		}
		if len(row) > 0 {
			bid := fmt.Sprintf("%v", row[0])
			if bid != "" {
				existing[bid] = i + 1 // 1-indexed row number
			}
		}
	}

	return existing, nil
}

// Upsert writes records to the destination sheet using upsert logic.
func (w *Writer) Upsert(records []model.FPPPRecord) error {
	existing, err := w.ReadExisting()
	if err != nil {
		// If the sheet/tab doesn't exist yet or is empty, treat as empty
		logger.Info("Could not read existing data (sheet may be empty): %v", err)
		existing = make(map[string]int)
	}

	// Write header if no existing data
	if len(existing) == 0 {
		headerRange := fmt.Sprintf("'%s'!A1:G1", destTab)
		_, err := w.srv.Spreadsheets.Values.Update(w.sheetID, headerRange, &sheets.ValueRange{
			Values: [][]interface{}{model.HeaderRow()},
		}).ValueInputOption("RAW").Do()
		if err != nil {
			return fmt.Errorf("writing header row: %w", err)
		}
		logger.Info("Wrote header row to destination sheet")
	}

	var updateData []*sheets.ValueRange
	var appendRows [][]interface{}

	for _, rec := range records {
		rowNum, exists := existing[rec.BusinessID]
		if exists {
			rng := fmt.Sprintf("'%s'!A%d:G%d", destTab, rowNum, rowNum)
			updateData = append(updateData, &sheets.ValueRange{
				Range:  rng,
				Values: [][]interface{}{rec.ToRow()},
			})
		} else {
			appendRows = append(appendRows, rec.ToRow())
		}
	}

	// Batch update existing rows
	if len(updateData) > 0 {
		_, err := w.srv.Spreadsheets.Values.BatchUpdate(w.sheetID, &sheets.BatchUpdateValuesRequest{
			ValueInputOption: "RAW",
			Data:             updateData,
		}).Do()
		if err != nil {
			return fmt.Errorf("batch updating rows: %w", err)
		}
		logger.Info("Updated %d existing rows", len(updateData))
	}

	// Batch append new rows
	if len(appendRows) > 0 {
		appendRange := fmt.Sprintf("'%s'!A:G", destTab)
		_, err := w.srv.Spreadsheets.Values.Append(w.sheetID, appendRange, &sheets.ValueRange{
			Values: appendRows,
		}).ValueInputOption("RAW").Do()
		if err != nil {
			return fmt.Errorf("appending new rows: %w", err)
		}
		logger.Info("Appended %d new rows", len(appendRows))
	}

	return nil
}
