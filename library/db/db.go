package db

type Where map[string]interface{}
type RowData map[string]interface{}
type RowDataArr []map[string]interface{}

// 对外输出的接口，底层实现可以任意，比如mysql，postsql之类的
// type IDB interface {
// 	Select(tableName string, conds where, fields []string, target interface{}) error
// 	// Update(tableName string, conds where, data rowData) (int64, error)
// 	// Delete(tableName string, conds where) (int64, error)
// 	// Insert(tableName string, data rowData) (int64, error)
// 	// Query(query string, args ...interface{}) (*sql.Rows, error)
// 	NamedQuery(namedQuery string, namedConds where, target interface{}) error
// }

func GetConn() *Mysql {
	// 这里可以随意替换了就
	return GetMysqlDB("db")
}
