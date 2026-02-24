package model

type FPPPRecord struct {
	BusinessID         string
	TitleForm          string
	Divisi             string
	TglFPPP            string
	NoFPPP             string
	DeadlinePengiriman string
	WaktuProduksi      string
	FinanceKlaes       string
	EndTime            string
}

func (r *FPPPRecord) ToRow() []interface{} {
	return []interface{}{
		r.BusinessID,
		r.TitleForm,
		r.Divisi,
		r.TglFPPP,
		r.NoFPPP,
		r.DeadlinePengiriman,
		r.WaktuProduksi,
		r.FinanceKlaes,
		r.EndTime,
	}
}

func HeaderRow() []interface{} {
	return []interface{}{
		"business_id",
		"Title Form",
		"Divisi",
		"Tgl FPPP",
		"No. FPPP",
		"Deadline Pengiriman",
		"Waktu Produksi",
		"Finance Klaes",
		"End Time",
	}
}
