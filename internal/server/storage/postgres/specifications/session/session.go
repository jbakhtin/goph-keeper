package session

import (
	"encoding/json"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/models"
	secondary_ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/server/storage/postgres/specifications/basic"
)

var _ secondary_ports.SessionQuerySpecification = &Specification{}

type Specification struct {
	specifications []secondary_ports.QuerySpecification
}

func NewSessionQuerySpecification() (Specification, error) {
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

func (s Specification) UserID(userId int) secondary_ports.QuerySpecification {
	return &UserIDSpecification{
		UserID: userId,
	}
}

func (s Specification) IsNotClosed() secondary_ports.QuerySpecification {
	return &IsNotClosedSpecification{}
}

func (s Specification) FingerPrint(fingerPrint models.FingerPrint) secondary_ports.QuerySpecification {
	buf, _ := json.Marshal(fingerPrint)
	return &FingerPrintSpecification{
		FingerPrint: string(buf),
	}
}

func (s Specification) RefreshToken(refreshToken string) secondary_ports.QuerySpecification {
	return &RefreshTokenSpecification{
		RefreshToken: refreshToken,
	}
}
