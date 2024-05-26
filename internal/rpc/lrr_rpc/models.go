package lrr_rpc

type Archive struct {
	Arcid        string `json:"arcid"`
	Extension    string `json:"extension"`
	Isnew        string `json:"isnew"`
	Lastreadtime int    `json:"lastreadtime"`
	Pagecount    int    `json:"pagecount"`
	Progress     int    `json:"progress"`
	Tags         string `json:"tags"`
	Title        string `json:"title"`
}
