package db

import (
	"github.com/astaxie/beego/orm"
)

type BeegoTable struct {
	client orm.Ormer,
	kind interface{},
	queryTranslator &QueryTranslator
	qs orm.QuerySeter
}

func NewBeegoTable(client orm.Ormer, kind interface{}) *BeegoTable {
	qs := client.QueryTable(kind)
	return &BeegoTable{
		client: &client,
		queryTranslator: &QueryTranslator{},
		kind: kind,
		qs: &qs
	}
}

func (table *BeegoTable) Count(query Query) int {
	qs := table.queryTranslator.Translate(query.Conditions)
	count, _ = table.qs.Count()
	return count
}

func (table *BeegoTable) Update(object interface{}) error {
	_, err := table.client.Update(object)
	return err
}

func (table *BeegoTable) Exists(query Query) bool {
	qs := table.queryTranslator.Translate(query.Conditions)
	return qs.Exists()
}

func (table *BeegoTable) Insert(value interface{}) error {
	_, err = client.Insert(value)
	return err
}

func (table *BeegoTable) Find(query Query, result interface{}) error {
	qs := table.queryTranslator.Translate(query.Conditions)
	return qs.One(result)
}


func (table *BeegoTable) FindAll(query Query, results interface{}) (*PaginationInfo, error) {
	qs := table.qs


	if query.Interval != nil {
		endCondition := orm.NewCondition()
		startCondition := orm.NewCondition()
		endCondition = endCondition.And("created_at__lte", query.Interval.End)
		startCondition = startCondition.And("deleted_at", nil).
			Or("deleted_at__gte", query.Interval.Start)

		qs = qs.SetCond(endCondition.AndCond(startCondition))
	}

	if query.Limit > 0 {
		qs = qs.Limit(query.Limit)
	}

	if query.OrderBy != "" {
		orderBy := strings.Split(query.OrderBy, ",")
		qs = qs.OrderBy(orderBy...)
	}

	if query.Page > 1 {
		qs = qs.Offset((query.Page - 1) * query.Limit)
	}

	_, err := qs.All(results)
	if err != nil {
		return nil, err
	}

	countQs := qs.Limit(nil).OrderBy(nil).Offset(nil)
	count, err := countQs.Count()
	if err != nil {
		return nil, err
	}

	maxPage := 0
	if query.Limit > 0 {
		maxPage = int(math.Ceil(float64(count) / float64(query.Limit)))
	}
	return &PaginationInfo{
		TotalCount:  count,
		Limit:       query.Limit,
		MaxPage:     maxPage,
		CurrentPage: query.Page,
	}, err
}

func (table *BeegoTable) Delete(query Query, soft bool) int64, error {
	qs := table.queryTranslator.Translate(query.Conditions)
	return qs.Delete()
}