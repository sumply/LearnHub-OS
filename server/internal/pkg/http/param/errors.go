package param

type QueryParam interface {
	Name() string
	Variants() []string
}

type QueryError struct {
	params []QueryParam
}

func NewQueryError(queryNumber int) *QueryError {
	return &QueryError{
		params: make([]QueryParam, 0, queryNumber),
	}
}

func (q *QueryError) Add(param QueryParam) {
	q.params = append(q.params, param)
}

func (q *QueryError) Empty() bool {
	return len(q.params) == 0
}

func (q *QueryError) Error() string {
	return "invalid query params"
}

func (q *QueryError) ToMap() map[string]any {
	m := make(map[string]any)

	for i := range q.params {
		m[q.params[i].Name()] = map[string]any{
			"variants": q.params[i].Variants(),
		}
	}

	return m
}
