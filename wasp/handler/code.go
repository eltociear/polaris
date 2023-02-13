package handler

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400

	ERROR_BLOCK_NOT_FOUND = 100001
	ERROR_TXN_NOT_FOUND   = 100002
)

var MsgFlags = map[int]string{
	SUCCESS:               "ok",
	ERROR:                 "fail",
	INVALID_PARAMS:        "Invalid params.",
	ERROR_BLOCK_NOT_FOUND: "Block not found.",
	ERROR_TXN_NOT_FOUND:   "Txn not found.",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
