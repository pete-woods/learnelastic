package tidyes

type Match struct {
	Match any `json:"match"`
}

func NewMatch(match any) *Match {
	return &Match{Match: match}
}

type Term struct {
	Term any `json:"term"`
}

func NewTerm(term any) *Term {
	return &Term{Term: term}
}

type Query struct {
	Query any   `json:"query"`
	Sort  []any `json:"sort,omitempty"`
}
