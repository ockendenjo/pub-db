package types

type PubsFile struct {
	Pubs []*Pub `json:"pubs"`
}

type Pub struct {
	CamraID    int     `json:"camraID"`
	GoodBeerID *int    `json:"goodBeerID,omitempty"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	RealAles   int     `json:"realAles"`
	Notes      *string `json:"notes,omitempty"`
	Chain      *string `json:"chain,omitempty"`
	TempClosed bool    `json:"tempClosed,omitzero"`
}
