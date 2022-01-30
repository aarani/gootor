package commands

const Version = 7
const NetInfo = 8
const Certs = 129
const AuthChallenge = 130

func IsVariableLength(command byte) bool {
	return command == 7 || command >= 128
}
