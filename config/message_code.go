package config

var messageList = map[string]string{
	"PARAM_ERROR":   "MSG_V0000", // param error
	"MAX_LENGTH":    "MSG_V0001", // param max length
	"FIX_LENGTH":    "MSG_V0002", // param fix length
	"FORMAT_NUMBER": "MSG_V0004", // param is number
	"FORMAT_DATE":   "MSG_V0003", // param format date is YYYY-MM-DD. Ex: 2023-01-01
	"REQUIRE":       "MSG_V0001", // Param require

	"KEY_NOT_FOUND":   "MSG_S0000",      // key error not found
	"SYSTEM_ERROR":    "MSG_S0001",      // system error
	"TOKEN_INCORRECT": "MSG_S0002",      // token invalid
	"GET_DATA_FAIL":   "MSG_RE0001",     //get data fail
	"CREATE_SUCCESS":	"MSG_CI0001", //Create new data success
	"NOT_ID_EXISTS" : "MSG_RE0002",//No item with that Id exists 
	"GET_DATA_SUCCESS": "MSG_RI0001", //Get data success

	"USERNAME_PASSWORD_INCORRECT": "MSG_N0000",
	"MISSING_FIELDS": "MSG_V1000",
	"UPDATE_SUCCESS": "MSG_UI0001",
	"DELETE_SUCCESS": "MSG_DI0001",

}

func GetMessageCode(key string) string {
	// var message = "MSG04000"
	var message = key
	if msg, ok := messageList[key]; ok {
		message = msg
	}

	return message
}
