package secrets

import (
	secondaryports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/secrets/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/server/storage/postgres/specifications/basic"
	"github.com/jbakhtin/goph-keeper/internal/server/storage/postgres/specifications/session"
	"github.com/jbakhtin/goph-keeper/pkg/queryspec"
)

var _ secondaryports.SecretQuerySpecification = &Specification{}

type Specification struct {
	specifications []queryspec.QuerySpecification
}

func NewKeyValueQuerySpecification() (Specification, error) {
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

func (s Specification) UserID(id int) queryspec.QuerySpecification {
	return &session.UserIDSpecification{
		UserID: id,
	}
}
