package database

type SqlPredicate string

var (
	SqlEqualPredicate              = SqlPredicate("= ?")
	SqlNotEqualPredicate           = SqlPredicate("<> ?")
	SqlGreaterThanPredicate        = SqlPredicate("> ?")
	SqlGreaterThanOrEqualPredicate = SqlPredicate(">= ?")
	SqlSmallerThanPredicate        = SqlPredicate("< ?")
	SqlSmallerThanOrEqualPredicate = SqlPredicate("<= ?")
	SqlIn                          = SqlPredicate("IN (?)")
	SqlNotIn                       = SqlPredicate("Not IN (?)")
	SqlBetweenAnd                  = SqlPredicate("BETWEEN ? AND ?")
	SqlLikeAll                     = SqlPredicate("LIKE  ?")
	SqlLikeLeft                    = SqlPredicate("LIKE  ?")
	SqlLikeRight                   = SqlPredicate("LIKE  ?")
)

/*
数据库查询匹配
*/
type SqlCondition struct {
	QueryName string
	Predicate SqlPredicate
	Value     interface{}
}

//查询条件
type SqlOptions struct {
	SelectField string
	OrderBy     map[string]string
}
