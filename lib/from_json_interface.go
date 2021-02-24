package lib

type JsonSerializable interface {
	FromJson(map[string]interface{}) interface{}

	ToJson()[]byte
}
