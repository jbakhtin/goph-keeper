package user

import (
	"github.com/jbakhtin/goph-keeper/internal/server/storage/postgres/specifications/basic"
	"github.com/jbakhtin/goph-keeper/pkg/queryspec"
)

type Specification struct {
	specifications []queryspec.QuerySpecification
}

func NewUserQuerySpecification() (Specification, error) {
	return Specification{
		specifications: make([]queryspec.QuerySpecification, 0),
	}, nil
}

func (s Specification) Limit(specification queryspec.QuerySpecification, i int) queryspec.QuerySpecification {
	return &basic.LimitSpecification{
		Specification: specification,
		Limit:         i,
	}
}

func (s Specification) Where(specification queryspec.QuerySpecification) queryspec.QuerySpecification {
	return &basic.WhereSpecification{
		Specification: specification,
	}
}

func (s Specification) Or(specifications ...queryspec.QuerySpecification) queryspec.QuerySpecification {
	return &basic.OrSpecification{
		Specifications: specifications,
	}
}

func (s Specification) And(specifications ...queryspec.QuerySpecification) queryspec.QuerySpecification {
	return &basic.AndSpecification{
		Specifications: specifications,
	}
}

func (s Specification) Email(email string) queryspec.QuerySpecification {
	return &EmailSpecification{
		Email: email,
	}
}
