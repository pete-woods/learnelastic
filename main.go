package main

import (
	"context"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"

	"github.com/pete-woods/learnelastic/tidyes"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

type doc struct {
	Name       string `json:"name,omitempty"`
	Department string `json:"department,omitempty"`
	Age        int    `json:"age,omitempty"`
	Country    string `json:"country,omitempty"`
}

func run() (err error) {
	esClient, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	es := tidyes.NewHelper(esClient)

	const today = "index"
	const yesterday = "yesterday"
	indices := []string{yesterday, today}

	for _, idx := range indices {
		err = es.Do(ctx, nil, esapi.IndicesCreateRequest{
			Index: idx,
			Body: esutil.NewJSONReader(tidyes.IndicesCreateRequest{
				Mappings: tidyes.Mappings{
					Properties: map[string]map[string]any{
						"age":        {"type": "integer"},
						"country":    {"type": "keyword"},
						"department": {"type": "keyword"},
						"name":       {"type": "text"},
					},
				},
			}),
		})
		if err != nil {
			return err
		}
	}

	if false {
		return nil
	}

	err = es.Do(ctx, nil, esapi.IndexRequest{
		Index:      today,
		DocumentID: "1",
		Body: esutil.NewJSONReader(doc{
			Name:       "bob",
			Age:        50,
			Department: "IT",
			Country:    "UK",
		}),
	})
	if err != nil {
		return err
	}

	err = es.Do(ctx, nil, esapi.IndexRequest{
		Index:      today,
		DocumentID: "2",
		Body: esutil.NewJSONReader(doc{
			Name:       "alice",
			Age:        56,
			Department: "IT",
			Country:    "UK",
		}),
	})
	if err != nil {
		return err
	}

	err = es.Do(ctx, nil, esapi.IndexRequest{
		Index:      today,
		DocumentID: "3",
		Body: esutil.NewJSONReader(doc{
			Name:       "charlie",
			Age:        26,
			Department: "IT",
			Country:    "France",
		}),
	})
	if err != nil {
		return err
	}

	err = es.Do(ctx, nil, esapi.IndexRequest{
		Index:      yesterday,
		DocumentID: "1",
		Body: esutil.NewJSONReader(doc{
			Name:       "abbie",
			Age:        45,
			Department: "IT",
			Country:    "Germany",
		}),
	})
	if err != nil {
		return err
	}
	err = es.Do(ctx, nil, esapi.IndexRequest{
		Index:      yesterday,
		DocumentID: "2",
		Body: esutil.NewJSONReader(doc{
			Name:       "bobby",
			Age:        26,
			Department: "IT",
			Country:    "Germany",
		}),
	})
	if err != nil {
		return err
	}

	var d tidyes.GetResult[doc]
	err = es.Do(ctx, &d, esapi.GetRequest{
		Index:      today,
		DocumentID: "1",
	})
	if err != nil {
		return err
	}

	log.Println("get")
	log.Printf("%#+v", d.Source)

	var s tidyes.SearchResult[doc]
	err = es.Do(ctx, &s, esapi.SearchRequest{
		Index: indices,
		Body: esutil.NewJSONReader(tidyes.Query{
			Query: tidyes.NewMatch(doc{
				Department: "IT",
			}),
			Sort: []any{
				map[string]string{"country": "asc"},
				map[string]string{"age": "asc"},
			},
		}),
	})
	if err != nil {
		return err
	}

	log.Println("search")
	for _, h := range s.Hits.Hits {
		log.Printf("%s:%s - %#+v", h.Index, h.ID, h.Source)
	}
	log.Println("search end")

	return nil
}
