package tidyes

type SearchHit[T any] struct {
	Index  string  `json:"_index"`
	Type   string  `json:"_type"`
	ID     string  `json:"_id"`
	Score  float64 `json:"_score"`
	Source T       `json:"_source"`
}

type SearchTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type SearchHits[T any] struct {
	Total    SearchTotal    `json:"total"`
	MaxScore float64        `json:"max_score"`
	Hits     []SearchHit[T] `json:"hits"`
}

type SearchResult[T any] struct {
	Took     int           `json:"took"`
	TimedOut bool          `json:"timed_out"`
	Hits     SearchHits[T] `json:"hits"`
}
