package service

import (
	"fmt"

	"DataFPP_Dingtalk/internal/sheets"
	"DataFPP_Dingtalk/pkg/logger"
)

type SyncService struct {
	reader *sheets.Reader
	writer *sheets.Writer
}

func NewSyncService(reader *sheets.Reader, writer *sheets.Writer) *SyncService {
	return &SyncService{reader: reader, writer: writer}
}

func (s *SyncService) SyncData() error {
	logger.Info("Starting data sync...")

	// 1. Read from FORM MASTER (already filtered by status & No. FPPP prefix)
	records, err := s.reader.ReadFormMaster()
	if err != nil {
		return fmt.Errorf("reading form master: %w", err)
	}
	logger.Info("Read %d records from FORM MASTER", len(records))

	// 2. Read from Comment (filtered by Finance Klaes)
	endTimeMap, financeKlaesMap, err := s.reader.ReadComments()
	if err != nil {
		return fmt.Errorf("reading comments: %w", err)
	}
	logger.Info("Read %d Finance Klaes comments", len(endTimeMap))

	// 3. Merge: left join FORM MASTER with Comment on business_id
	for i := range records {
		if endTime, ok := endTimeMap[records[i].BusinessID]; ok {
			records[i].EndTime = endTime
		}
		if fk, ok := financeKlaesMap[records[i].BusinessID]; ok {
			records[i].FinanceKlaes = fk
		}
	}

	// 4. Clear sheet and write all records fresh
	if err := s.writer.ClearAndWrite(records); err != nil {
		return fmt.Errorf("writing data: %w", err)
	}

	logger.Info("Data sync completed successfully. Total records: %d", len(records))
	return nil
}
