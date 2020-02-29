package db

import (
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/pilly-io/api/internal/models"
)

type BeegoTable struct {
	client          orm.Ormer
	kind            interface{}
	queryTranslator *QueryTranslator
	qs              orm.QuerySeter
}

func NewBeegoTable(client orm.Ormer, kind interface{}) *BeegoTable {
	qs := client.QueryTable(kind)
	return &BeegoTable{
		client:          client,
		queryTranslator: &QueryTranslator{qs},
		kind:            kind,
		qs:              qs,
	}
}

func (table *BeegoTable) Count(query Query) int {
	qs := table.queryTranslator.Translate(&query.Conditions)
	if query.ExcludeDeleted {
		qs = qs.Filter("deleted_at__isnull", true)
	}
	count, _ := qs.Count()
	return int(count)
}

func (table *BeegoTable) Exists(query Query) bool {
	qs := table.queryTranslator.Translate(&query.Conditions)
	if query.ExcludeDeleted {
		qs = qs.Filter("deleted_at__isnull", true)
	}
	return qs.Exist()
}

func (table *BeegoTable) Insert(value interface{}) error {
	value.(models.PersistedModel).BeforeSave()
	_, err := table.client.Insert(value)
	return err
}

func (table *BeegoTable) Update(value interface{}) error {
	value.(models.PersistedModel).BeforeSave()
	_, err := table.client.Update(value)
	return err
}

func (table *BeegoTable) Find(query Query, result interface{}) error {
	qs := table.queryTranslator.Translate(&query.Conditions)
	if query.ExcludeDeleted {
		qs = qs.Filter("deleted_at__isnull", true)
	}
	err := qs.One(result)
	if err != nil {
		return err
	}
	result.(models.PersistedModel).AfterLoad()
	return nil
}

func (table *BeegoTable) FindAll(query Query, results interface{}) (*PaginationInfo, error) {
	qs := table.queryTranslator.Translate(&query.Conditions)
	if query.Interval != nil {
		endCondition := orm.NewCondition()
		startCondition := orm.NewCondition()
		endCondition = endCondition.And("created_at__lte", query.Interval.End)
		startCondition = startCondition.And("deleted_at__isnull", true).
			Or("deleted_at__gte", query.Interval.Start)

		qs = qs.SetCond(qs.GetCond().AndCond(endCondition.AndCond(startCondition)))
	} else if query.ExcludeDeleted {
		qs = qs.Filter("deleted_at__isnull", true)
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

	count, err := qs.Count()
	if err != nil {
		return nil, err
	}

	maxPage := 0
	if query.Limit > 0 {
		maxPage = int(math.Ceil(float64(count) / float64(query.Limit)))
	}

	objects := reflect.Indirect(reflect.ValueOf(results))
	for i := 0; i < objects.Len(); i++ {
		model := objects.Index(i).Interface().(models.PersistedModel)
		model.AfterLoad()
	}

	return &PaginationInfo{
		TotalCount:  int(count),
		Limit:       query.Limit,
		MaxPage:     maxPage,
		CurrentPage: query.Page,
	}, err
}

func (table *BeegoTable) Delete(query Query, soft bool) (int64, error) {
	qs := table.queryTranslator.Translate(&query.Conditions)
	if soft {
		return qs.Update(orm.Params{"deleted_at": time.Now().UTC()})
	}
	return qs.Filter("id__isnull", false).Delete()
}
