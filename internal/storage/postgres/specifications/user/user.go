package user

import (
	ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/server/storage/postgres/specifications/basic"
)

type Specification struct {
	specifications []ports.QuerySpecification
}

func NewUserQuerySpecification() (Specification, error) {
	return Specification{
		specifications: make([]ports.QuerySpecification, 0),
	}, nil
}

func (s Specification) Limit(specification ports.QuerySpecification, i int) ports.QuerySpecification {
	return &basic.LimitSpecification{
		Specification: specification,
		Limit:         i,
	}
}

func (s Specification) Where(specification ports.QuerySpecification) ports.QuerySpecification {
	return &basic.WhereSpecification{
		Specification: specification,
	}
}

func (s Specification) Or(specifications ...ports.QuerySpecification) ports.QuerySpecification {
	return &basic.OrSpecification{
		Specifications: specifications,
	}
}

func (s Specification) And(specifications ...ports.QuerySpecification) ports.QuerySpecification {
	return &basic.AndSpecification{
		Specifications: specifications,
	}
}

func (s Specification) Email(email string) ports.QuerySpecification {
	return &EmailSpecification{
		Email: email,
	}
}
