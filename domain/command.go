// Copyright (c) 2024 Benedict Weis. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package domain

import "strconv"

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

// CommandTypeFromString is used to mach a received command as a string to the CommandType used to communicate the command.
func CommandTypeFromString(s string) CommandType {
	for currentCommandType := Unknown; currentCommandType <= Quit; currentCommandType++ {
		if currentCommandType.String() == s {
			return currentCommandType
		}
	}
	return Unknown
}

// String implements the string variants of CommandType.
func (c CommandType) String() string {
	commandTypeToStringMapping := []string{"unknown", "name", "msg", "acc", "login", "passwd", "info", "who", "quit"}
	if c < 0 || int(c) > len(commandTypeToStringMapping)-1 {
		return strconv.Itoa(int(c))
	}
	return commandTypeToStringMapping[c]
}

// Command represents a command a user wants to be executed.
type Command struct {
	SessionID   string
	CommandType CommandType
	Arguments   []string
}
