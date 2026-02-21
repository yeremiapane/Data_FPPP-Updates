package sheets

import (
	"fmt"

	"DataFPP_Dingtalk/internal/model"

	"google.golang.org/api/sheets/v4"
)

type Reader struct {
	srv     *sheets.Service
	sheetID string
}

func NewReader(srv *sheets.Service, sheetID string) *Reader {
	return &Reader{srv: srv, sheetID: sheetID}
}

func (r *Reader) ReadFormMaster() ([]model.FPPPRecord, error) {
	resp, err := r.srv.Spreadsheets.Values.Get(r.sheetID, "FORM MASTER!A:F").Do()
	if err != nil {
		return nil, fmt.Errorf("reading FORM MASTER: %w", err)
	}

	if len(resp.Values) < 2 {
		return nil, nil
	}

	header := resp.Values[0]
	colIdx := mapColumns(header, []string{"business_id", "Title Form", "Divisi", "Tgl FPPP", "No.FPPP", "Deadline Pengiriman"})

	var records []model.FPPPRecord
	for _, row := range resp.Values[1:] {
		bid := cellValue(row, colIdx["business_id"])
		if bid == "" {
			continue
		}
		records = append(records, model.FPPPRecord{
			BusinessID:         bid,
			TitleForm:          cellValue(row, colIdx["Title Form"]),
			Divisi:             cellValue(row, colIdx["Divisi"]),
			TglFPPP:            cellValue(row, colIdx["Tgl FPPP"]),
			NoFPPP:             cellValue(row, colIdx["No.FPPP"]),
			DeadlinePengiriman: cellValue(row, colIdx["Deadline Pengiriman"]),
		})
	}

	return records, nil
}

// ReadComments reads the Comment sheet and returns a map of business_id -> end_time
// filtered to only rows where show_name_true == "Finance Klaes".
func (r *Reader) ReadComments() (map[string]string, error) {
	resp, err := r.srv.Spreadsheets.Values.Get(r.sheetID, "Comment!A:C").Do()
	if err != nil {
		return nil, fmt.Errorf("reading Comment: %w", err)
	}

	if len(resp.Values) < 2 {
		return nil, nil
	}

	header := resp.Values[0]
	colIdx := mapColumns(header, []string{"business_id", "show_name_true", "end_time"})

	result := make(map[string]string)
	for _, row := range resp.Values[1:] {
		showName := cellValue(row, colIdx["show_name_true"])
		if showName != "Finance Klaes" {
			continue
		}
		bid := cellValue(row, colIdx["business_id"])
		if bid == "" {
			continue
		}
		result[bid] = cellValue(row, colIdx["end_time"])
	}

	return result, nil
}

func mapColumns(header []interface{}, names []string) map[string]int {
	idx := make(map[string]int)
	for _, name := range names {
		idx[name] = -1
	}
	for i, cell := range header {
		s := fmt.Sprintf("%v", cell)
		if _, ok := idx[s]; ok {
			idx[s] = i
		}
	}
	return idx
}

func cellValue(row []interface{}, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return fmt.Sprintf("%v", row[idx])
}
