package session

import (
	"encoding/json"

	"github.com/jbakhtin/goph-keeper/pkg/queryspec"

	"github.com/jbakhtin/goph-keeper/internal/appmodules/auth/domain/models"
	secondary_ports "github.com/jbakhtin/goph-keeper/internal/appmodules/auth/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/storage/postgres/specifications/basic"
)

var _ secondary_ports.SessionQuerySpecification = &Specification{}

type Specification struct {
	specifications []queryspec.QuerySpecification
}

func NewSessionQuerySpecification() (Specification, error) {
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

func (s Specification) UserID(userID int) queryspec.QuerySpecification {
	return &UserIDSpecification{
		UserID: userID,
	}
}

func (s Specification) IsNotClosed() queryspec.QuerySpecification {
	return &IsNotClosedSpecification{}
}

func (s Specification) FingerPrint(fingerPrint models.FingerPrint) queryspec.QuerySpecification {
	buf, _ := json.Marshal(fingerPrint)
	return &FingerPrintSpecification{
		FingerPrint: string(buf),
	}
}

func (s Specification) RefreshToken(refreshToken string) queryspec.QuerySpecification {
	return &RefreshTokenSpecification{
		RefreshToken: refreshToken,
	}
}
