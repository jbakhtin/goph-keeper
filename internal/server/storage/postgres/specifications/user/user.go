package user

import (
	secondary_ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/server/storage/postgres/specifications/basic"
)

type Specification struct {
	specifications []secondary_ports.QuerySpecification
}

func NewUserQuerySpecification() (Specification, error) {
	return Specification{
		specifications: make([]secondary_ports.QuerySpecification, 0),
	}, nil
}

func (s Specification) Limit(specification secondary_ports.QuerySpecification, i int) secondary_ports.QuerySpecification {
	return &basic.LimitSpecification{
		Specification: specification,
		Limit:         i,
	}
}

func (s Specification) Where(specification secondary_ports.QuerySpecification) secondary_ports.QuerySpecification {
	return &basic.WhereSpecification{
		Specification: specification,
	}
}

func (s Specification) Or(specifications ...secondary_ports.QuerySpecification) secondary_ports.QuerySpecification {
	return &basic.OrSpecification{
		Specifications: specifications,
	}
}

func (s Specification) And(specifications ...secondary_ports.QuerySpecification) secondary_ports.QuerySpecification {
	return &basic.AndSpecification{
		Specifications: specifications,
	}
}

func (s Specification) Email(email string) secondary_ports.QuerySpecification {
	return &EmailSpecification{
		Email: email,
	}
}
