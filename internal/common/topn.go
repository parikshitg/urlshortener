package common

// TopN is a common struct to fetch and store top N domains.
type TopN struct {
	Rank      int    `json:"rank"`
	Domain    string `json:"domain"`
	Shortened int    `json:"shortened"`
}
