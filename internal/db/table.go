package db

type Table interface {
	Find(query Query, result interface{}) error
	FindAll(query Query, results interface{}) (*PaginationInfo, error)
	Exists(query Query) bool
	Insert(value interface{}) error
	BulkInsert(value interface{}) error
	Count(query Query) int
	Update(object interface{}) error
	Delete(query Query, soft bool) (int64, error)
}
