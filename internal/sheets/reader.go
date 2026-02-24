package sheets

import (
	"fmt"
	"math"
	"strings"
	"time"

	"DataFPP_Dingtalk/internal/model"
	"DataFPP_Dingtalk/pkg/logger"

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
	resp, err := r.srv.Spreadsheets.Values.Get(r.sheetID, "FORM MASTER!A:CZ").Do()
	if err != nil {
		return nil, fmt.Errorf("reading FORM MASTER: %w", err)
	}

	if len(resp.Values) < 2 {
		return nil, nil
	}

	header := resp.Values[0]
	logger.Info("FORM MASTER headers: %v", header)
	colIdx := mapColumns(header, []string{
		"business_id", "status", "Title Form", "Divisi",
		"Tgl FPPP", "No. FPPP", "Deadline Pengiriman", "Deadline Pengambilan", "Waktu Produksi",
	})
	logger.Info("FORM MASTER column index map: %v", colIdx)

	// Allowed prefixes for No. FPPP filter
	allowedPrefixes := []string{"AST", "ASTA", "RSD", "RAE", "RAS"}

	var records []model.FPPPRecord
	for _, row := range resp.Values[1:] {
		bid := cellValue(row, colIdx["business_id"])
		if bid == "" {
			continue
		}

		// Filter by status (Complete/Completed or Running)
		status := strings.ToUpper(strings.TrimSpace(cellValue(row, colIdx["status"])))
		if status != "COMPLETE" && status != "COMPLETED" && status != "RUNNING" {
			continue
		}

		// Filter by No. FPPP prefix that contains allowed Prefixes
		noFPPP := strings.TrimSpace(cellValue(row, colIdx["No. FPPP"]))
		if !hasAllowedPrefix(noFPPP, allowedPrefixes) {
			continue
		}

		tglFPPP := strings.TrimSpace(cellValue(row, colIdx["Tgl FPPP"]))
		deadline := strings.TrimSpace(cellValue(row, colIdx["Deadline Pengiriman"]))
		waktuProduksi := strings.TrimSpace(cellValue(row, colIdx["Waktu Produksi"]))

		// Treat "null" string from sheet as empty
		if strings.EqualFold(waktuProduksi, "null") {
			waktuProduksi = ""
		}
		if strings.EqualFold(deadline, "null") {
			deadline = ""
		}
		if strings.EqualFold(tglFPPP, "null") {
			tglFPPP = ""
		}

		// Fallback: jika Deadline Pengiriman kosong, coba Deadline Pengambilan
		if deadline == "" {
			deadlineAmbil := strings.TrimSpace(cellValue(row, colIdx["Deadline Pengambilan"]))
			if !strings.EqualFold(deadlineAmbil, "null") {
				deadline = deadlineAmbil
			}
		}

		// If Waktu Produksi is empty, calculate from Deadline Pengiriman - Tgl FPPP
		if waktuProduksi == "" {
			waktuProduksi = calcWaktuProduksi(tglFPPP, deadline)
		}

		records = append(records, model.FPPPRecord{
			BusinessID:         bid,
			TitleForm:          cellValue(row, colIdx["Title Form"]),
			Divisi:             cellValue(row, colIdx["Divisi"]),
			TglFPPP:            tglFPPP,
			NoFPPP:             noFPPP,
			DeadlinePengiriman: deadline,
			WaktuProduksi:      waktuProduksi,
		})
	}

	return records, nil
}

// hasAllowedPrefix returns true if noFPPP contains any of the allowed prefixes (case-insensitive).
// Handles formats like "042/FPPP/ASTA/01/2026" where the prefix is in the middle.
func hasAllowedPrefix(noFPPP string, prefixes []string) bool {
	upper := strings.ToUpper(noFPPP)
	for _, p := range prefixes {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

// calcWaktuProduksi calculates the number of days between tglFPPP and deadline.
// Returns empty string if either date cannot be parsed.
func calcWaktuProduksi(tglFPPP, deadline string) string {
	layouts := []string{
		"2006-01-02",
		"02/01/2006",
		"01/02/2006",
		"2006/01/02",
		"02-01-2006",
		"January 2, 2006",
		"2 January 2006",
	}
	parse := func(s string) (time.Time, bool) {
		for _, l := range layouts {
			if t, err := time.Parse(l, s); err == nil {
				return t, true
			}
		}
		return time.Time{}, false
	}

	t1, ok1 := parse(tglFPPP)
	t2, ok2 := parse(deadline)
	if !ok1 || !ok2 {
		return ""
	}

	days := int(math.Round(t2.Sub(t1).Hours() / 24))
	return fmt.Sprintf("%d hari", days)
}

// ReadComments reads the Comment sheet and returns two maps keyed by business_id:
//   - endTime map: business_id -> end_time (only rows where show_name_true == "Finance Klaes")
//   - financeKlaes map: business_id -> show_name_true value ("Finance Klaes") if present
func (r *Reader) ReadComments() (map[string]string, map[string]string, error) {
	resp, err := r.srv.Spreadsheets.Values.Get(r.sheetID, "Comment!A:Z").Do()
	if err != nil {
		return nil, nil, fmt.Errorf("reading Comment: %w", err)
	}

	if len(resp.Values) < 2 {
		return nil, nil, nil
	}

	header := resp.Values[0]
	logger.Info("Comment headers (full): %v", header)
	colIdx := mapColumns(header, []string{"business_id", "show_name_true", "end_time"})
	logger.Info("Comment column index map: %v", colIdx)

	endTimeMap := make(map[string]string)
	financeKlaesMap := make(map[string]string)

	for _, row := range resp.Values[1:] {
		showName := strings.TrimSpace(cellValue(row, colIdx["show_name_true"]))
		if !strings.EqualFold(showName, "Finance Klaes") {
			continue
		}
		bid := cellValue(row, colIdx["business_id"])
		if bid == "" {
			continue
		}
		endTimeMap[bid] = cellValue(row, colIdx["end_time"])
		financeKlaesMap[bid] = showName
	}

	return endTimeMap, financeKlaesMap, nil
}

func mapColumns(header []interface{}, names []string) map[string]int {
	idx := make(map[string]int)
	for _, name := range names {
		idx[name] = -1
	}
	for i, cell := range header {
		s := strings.TrimSpace(fmt.Sprintf("%v", cell))
		for _, name := range names {
			// First-match-wins: skip jika sudah ditemukan sebelumnya
			if strings.EqualFold(s, name) && idx[name] == -1 {
				idx[name] = i
				break
			}
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
