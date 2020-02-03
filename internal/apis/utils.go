package apis

type jsonFormat = map[string]interface{}

// ErrorsToJSON returns JSON format for multiple errors
func ErrorsToJSON(errors []error) map[string]interface{} {

}

func ObjectsToJSON(collection db.PaginatedCollection) map[string]interface{} {
	return jsonFormat{
		"data": collection.Objects,
		"pagination": {
			"current": collection.CurrentPage,
			"limit": collection.Limit,
			"max_page": collection.MaxPage,
			"total_count": collection.TotalCount,
		}
	}
}
