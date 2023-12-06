package types

type LogoutType int32

const (
	LogoutTypeThis LogoutType = 0
	LogoutTypeAll  LogoutType = 1
)

var (
	LogoutTypeName = map[LogoutType]string{
		LogoutTypeThis: "this",
		LogoutTypeAll:  "all",
	}

	LogoutTypeValue = map[string]LogoutType{
		LogoutTypeName[0]: LogoutTypeThis,
		LogoutTypeName[1]: LogoutTypeAll,
	}
)
