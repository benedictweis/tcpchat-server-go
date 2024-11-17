package domain

type CommandType int

const (
	Unknown CommandType = iota
	ChangeName
	PrivateMessage
	CreateAccount
	Login
	ChangePassword
	Info
	Who
	Quit
)

// String implements the string variants of CommandType
func (c CommandType) String() string {
	return [...]string{"unknown", "name", "msg", "acc", "login", "passwd", "info", "who", "quit"}[c]
}

// MatchCommandTypeStringToCommandType is used to mach a received command as a string to the CommandType used to communicate the command
func MatchCommandTypeStringToCommandType(s string) CommandType {
	for currentCommandType := Unknown; currentCommandType <= Quit; currentCommandType++ {
		if currentCommandType.String() == s {
			return currentCommandType
		}
	}
	return Unknown
}

// Command represents a command a user wants to be executed
type Command struct {
	SessionId   string
	CommandType CommandType
	Arguments   []string
}
