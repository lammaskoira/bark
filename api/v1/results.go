package v1

type Report struct {
	// Defines the version of the result format.
	Version string `json:"version" yaml:"version"`
	// Overall defines the overall status.
	Overall ResultStatus `json:"overall" yaml:"overall"`
	// Details contains the details of the result.
	Results []*ContextualResult `json:"results" yaml:"results"`
}

func (r *Report) AddResult(cr *ContextualResult) {
	r.Results = append(r.Results, cr)
}

func (r *Report) GatherOverall() {
	overall := ResultStatusPass
	for _, cr := range r.Results {
		for _, result := range cr.Results {
			if result.Status == ResultStatusError {
				r.Overall = ResultStatusError
				return
			}
			if result.Status == ResultStatusFail {
				overall = ResultStatusFail
			}
		}
	}

	r.Overall = overall
}

// Specifies a result for a single context.
type ContextualResult struct {
	// Defines the version of the result format.
	Version string   `json:"version" yaml:"version"`
	Target  string   `json:"target" yaml:"target"`
	Results []Result `json:"results" yaml:"results"`
}

func (r *ContextualResult) AddResult(result Result) {
	// TODO: Generate ID.
	r.Results = append(r.Results, result)
}

type ResultStatus string

const (
	ResultStatusPass  ResultStatus = "pass"
	ResultStatusFail  ResultStatus = "fail"
	ResultStatusError ResultStatus = "error"
)

type Result struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Status       ResultStatus `json:"status"`
	ErrorMessage string       `json:"errorMessage,omitempty"`
}
