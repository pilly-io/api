package db

import (
	"github.com/astaxie/beego/orm"
)

type QueryTranslator struct {
	qs orm.QuerySeter
}

//Translate a QueryCondition in Where clauses
func (translator QueryTranslator) Translate(conditions *QueryConditions) orm.QuerySeter {
	qs := translator.qs
	for key, value := range *conditions {
		if key == "or__" {
			orQuerySet := translator.Translate(value.(*QueryConditions))
			mergedCondition := qs.GetCond().AndCond(orQuerySet.GetCond())
			qs = qs.SetCond(mergedCondition)
		} else {
			qs = qs.Filter(key, value)
		}
	}
	return qs
}
