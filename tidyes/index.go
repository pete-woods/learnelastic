package tidyes

type Properties map[string]map[string]any

type IndicesCreateRequest struct {
	Mappings Mappings `json:"mappings"`
}

type Mappings struct {
	Properties Properties `json:"properties"`
}

type PutMappingRequest struct {
	Properties Properties `json:"properties"`
}
