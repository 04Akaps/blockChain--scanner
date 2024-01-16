package types

type errString string

const (
	BlockNumberInvalid errString = "startBlock over latestBlock"
	CanNotFindBlock    errString = "can't find block"
)

func ErrToString(s errString) string {
	return string(s)
}
