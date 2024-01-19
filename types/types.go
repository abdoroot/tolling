package types

type Distance struct {
	OBUID int     `json:"obuid"`
	Value float64 `json:"value"`
	Unix  int64   `json:"unix"`
}

type OBUdata struct {
	OBUID int     `json:"obuid"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}
