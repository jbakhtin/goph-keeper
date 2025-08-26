package queryspec

type QuerySpecification interface {
	Query() string
	Value() []any
}
