package model

type FPPPRecord struct {
	BusinessID         string
	TitleForm          string
	Divisi             string
	TglFPPP            string
	NoFPPP             string
	DeadlinePengiriman string
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
		r.EndTime,
	}
}

func HeaderRow() []interface{} {
	return []interface{}{
		"business_id",
		"Title Form",
		"Divisi",
		"Tgl FPPP",
		"No.FPPP",
		"Deadline Pengiriman",
		"End Time",
	}
}
