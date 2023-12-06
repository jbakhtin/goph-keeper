package types

type LogoutType int32

const (
	LogoutType_THIS LogoutType = 0
	LogoutType_ALL  LogoutType = 1
)

var (
	LogoutType_name = map[LogoutType]string{
		LogoutType_THIS: "this",
		LogoutType_ALL:  "all",
	}

	LogoutType_value = map[string]LogoutType{
		LogoutType_name[0]: LogoutType_THIS,
		LogoutType_name[1]: LogoutType_ALL,
	}
)
