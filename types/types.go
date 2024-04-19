package types

type OBUData struct {
	OBUID int     `json:"obuId"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

type Distance struct {
	Value float64 `json:"value"`
	OBUID int     `json:"obuID"`
	Unix  int64   `json:"unix"`
}

type Invoice struct {
	OBUID         int     `json:"obuID"`
	TotalDistance float64 `json:"totalDistance"`
	InvoiceAmount float64 `json:"invoiceAmount"`
}

func NewDistance(value float64, obuID int, unix int64) Distance {
	return Distance{
		Value: value,
		OBUID: obuID,
		Unix:  unix,
	}
}

func NewInvoice(obuID int, totalDist float64, invAmount float64) *Invoice {
	return &Invoice{
		OBUID:         obuID,
		TotalDistance: totalDist,
		InvoiceAmount: invAmount,
	}
}
