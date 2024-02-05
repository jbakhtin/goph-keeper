package session

import (
	"encoding/json"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/models"
	ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/server/storage/postgres/specifications/basic"
)

var _ ports.SessionQuerySpecification = &Specification{}

type Specification struct {
	specifications []ports.QuerySpecification
}

func NewSessionQuerySpecification() (Specification, error) {
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

func (s Specification) UserID(userID int) ports.QuerySpecification {
	return &UserIDSpecification{
		UserID: userID,
	}
}

func (s Specification) IsNotClosed() ports.QuerySpecification {
	return &IsNotClosedSpecification{}
}

func (s Specification) FingerPrint(fingerPrint models.FingerPrint) ports.QuerySpecification {
	buf, _ := json.Marshal(fingerPrint)
	return &FingerPrintSpecification{
		FingerPrint: string(buf),
	}
}

func (s Specification) RefreshToken(refreshToken string) ports.QuerySpecification {
	return &RefreshTokenSpecification{
		RefreshToken: refreshToken,
	}
}
