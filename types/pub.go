package types

type PubsFile struct {
	Schema string `json:"$schema"`
	Pubs   []*Pub `json:"pubs"`
}

type Pub struct {
	CamraID    int     `json:"camraID"`
	GoodBeerID *int    `json:"goodBeerID,omitempty"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	RealAles   int     `json:"realAles"`
	NumBeers   int     `json:"numBeers"`
	HasRealAle bool    `json:"hasRealAle"`
	Notes      *string `json:"notes,omitempty"`
	Chain      *string `json:"chain,omitempty"`
	TempClosed bool    `json:"tempClosed,omitzero"`
}
