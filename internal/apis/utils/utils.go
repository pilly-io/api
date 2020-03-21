package utils

import (
	"strconv"
	"time"
)

type jsonFormat = map[string]interface{}

// ConvertTimestampToTime : convert a ts to a time
func ConvertTimestampToTime(ts string) (*time.Time, error) {
	toInt, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return nil, err
	}
	toTime := time.Unix(int64(toInt), 0)
	return &toTime, err
}

// ErrorsToJSON returns JSON format for multiple errors
func ErrorsToJSON(errors ...error) map[string]interface{} {
	var errorsString []string
	for _, e := range errors {
		errorsString = append(errorsString, e.Error())
	}
	return jsonFormat{
		"errors": errorsString,
	}
}

//ObjectToJSON returns JSON format for the data
func ObjectToJSON(object interface{}) map[string]interface{} {
	return jsonFormat{
		"data": object,
	}
}

/*func PaginatedObjectToJSON(collection db.PaginationInfo) map[string]interface{} {
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

// GetFullKindName : Given a shortname retrieve the full name of the Kubernetes kind
func GetFullKindName(shortname string) string {
	switch shortname {
	case "po":
		return "pod"
	case "sts":
		return "statefulset"
	case "dep", "deploy":
		return "deployment"
	case "rs":
		return "replicaset"
	case "cj":
		return "cronjob"
	case "ds":
		return "daemonset"
	default:
		return shortname
	}
}
