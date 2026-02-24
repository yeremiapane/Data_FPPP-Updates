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

// ClearAndWrite clears the destination sheet then writes header + all records from scratch.
func (w *Writer) ClearAndWrite(records []model.FPPPRecord) error {
	clearRange := fmt.Sprintf("'%s'", destTab)

	// 1. Clear all existing content
	_, err := w.srv.Spreadsheets.Values.Clear(w.sheetID, clearRange, &sheets.ClearValuesRequest{}).Do()
	if err != nil {
		return fmt.Errorf("clearing destination sheet: %w", err)
	}
	logger.Info("Cleared destination sheet '%s'", destTab)

	// 2. Build all rows: header first, then data
	allRows := [][]interface{}{model.HeaderRow()}
	for _, rec := range records {
		allRows = append(allRows, rec.ToRow())
	}

	// 3. Write everything in one call
	writeRange := fmt.Sprintf("'%s'!A1", destTab)
	_, err = w.srv.Spreadsheets.Values.Update(w.sheetID, writeRange, &sheets.ValueRange{
		Values: allRows,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return fmt.Errorf("writing data to destination sheet: %w", err)
	}

	logger.Info("Wrote header + %d data rows to '%s'", len(records), destTab)
	return nil
}
