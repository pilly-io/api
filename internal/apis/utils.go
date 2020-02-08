package apis

type jsonFormat = map[string]interface{}

// ErrorsToJSON returns JSON format for multiple errors
func ErrorsToJSON(object interface{}) map[string]interface{} {
	return jsonFormat{
		"errors": "errors",
	}
}

func ObjectToJSON(object interface{}) map[string]interface{} {
	return jsonFormat{
		"data": object,
	}
}

/*func ObjectsToJSON(collection db.PaginationInfo) map[string]interface{} {
	return jsonFormat{
		"data": collection.Objects,
		"pagination": {
			"current":     collection.CurrentPage,
			"limit":       collection.Limit,
			"max_page":    collection.MaxPage,
			"total_count": collection.TotalCount,
		},
	}
}*/
