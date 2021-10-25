package message

var MessageFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "invalid params",
}

func GetMessage(code int) string {
	message, ok := MessageFlags[code]
	if ok {
		return message
	}

	return MessageFlags[ERROR]
}
