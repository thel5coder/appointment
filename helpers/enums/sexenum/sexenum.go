package sexenum

type sexMapEnums map[string]interface{}
var (
	sexEnums = []sexMapEnums{
		map[string]interface{}{
			"key":"male",
			"text":"Laki - laki",
		},
		map[string]interface{}{
			"key":"female",
			"text":"Perempuan",
		},
	}
)

func GetKey(index int) string{
	value := sexEnums[index]

	return value["key"].(string)
}

func GetText(index int) string{
	value := sexEnums[index]

	return value["text"].(string)
}

func GetEnums() []sexMapEnums{
	return sexEnums
}

func GetEnum(index int) sexMapEnums{
	return sexEnums[index]
}
