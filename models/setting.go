package model

type (
	// Setting collection
	Setting struct {
		TopNav []TopNav `json:"top_nav" bson:"top_nav"`
	}

	// TopNav columns
	TopNav struct {
		UID     string `json:"uid"`
		View    string `json:"view"`
		Class   string `json:"class"`
		Title   string `json:"title"`
		IsView  bool   `json:"isView"`
		Color   string `json:"color"`
		HashKey string `json:"$$hashKey"`
	}
)
