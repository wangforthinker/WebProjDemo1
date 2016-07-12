package main

import (
	"reflect"
	"fmt"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
)

type SqlProxy struct {
	//sql config
	connection  string
	name string
	passwd string
	db *sql.DB
	
	//ReadData , UpdateData , DeleteData , AddData 
}

func (proxy *SqlProxy) Connect(connect_addr string, connect_username string, passwd string, database string) error {
	proxy.connection = connect_addr
	proxy.name = connect_username
	proxy.passwd = passwd
	db,err := sql.Open("mysql",connect_username+":"+passwd+"@tcp("+connect_addr+")/"+database+"?charset=utf8")
	
	if err != nil {
		return err
	}
	
	proxy.db = db
	return nil
}


//data should be simple struct
func (proxy *SqlProxy) InsertData (table_name string, data interface{}) (string,error) {
	if proxy.db == nil {
		return "get db error",nil;
	}
	

	relem := reflect.ValueOf(data).Elem()
	col_names := make([]string, relem.NumField())
	data_values := make([]interface{}, relem.NumField())
	
	rtype := relem.Type()
	
	for i:=0 ; i < relem.NumField(); i++ {
		f := relem.Field(i)
		col_names[i] = rtype.Field(i).Name
		data_values[i] = f.Interface()
	}
	
	var sql_col_str string
	var value_col_str string
	
	{
		sql_col_str = col_names[0]
		value_col_str = "?"
		for i := 1; i < len(col_names); i ++ {
			sql_col_str += ","
			sql_col_str += col_names[i]
			
			value_col_str += ","
			value_col_str += "?"
		}
		
		fmt.Printf("sql string:%s\n",sql_col_str)
	}
	
	stmt,err := proxy.db.Prepare("insert into " + table_name + "(" + sql_col_str +") values(" + value_col_str +")")
	if err != nil {
		fmt.Println("prepare error:",err.Error())
		return "prepare error",err
	}
	
	defer stmt.Close()
	
	//result, err := stmt.Exec("wang","123456")
	result, err := stmt.Exec(data_values...)
	
	
	if err != nil {
		fmt.Println("insert error")
		return	"insert data error",err
	} else {
		fmt.Println(result.LastInsertId())
	}
	
	return "",nil
}

// tablename , colname, val, colname,val ...
func (proxy *SqlProxy) CheckDup(table_name string, col_info ...interface{}) (bool,error) {
	
	var id int
	arrlen := len(col_info)
	//assert arrlen%2 == 0
	
	sqlstr := string("select id from " + table_name + " where ")
	checkval := make([]interface{}, arrlen/2)
	
	{
		sqlstr += col_info[0].(string)
		sqlstr += " = ? "
		checkval[0] = col_info[1]
		
		for i := 1 ; i < arrlen/2 ; i ++ {
			sqlstr += " and "
			sqlstr += col_info[i * 2].(string)
			sqlstr += " = ? "
			checkval[i] = col_info[i * 2 + 1]
		}	
	}
	
	fmt.Println(sqlstr)
	fmt.Println(checkval)
	
	row := proxy.db.QueryRow(sqlstr,checkval...)
	err := row.Scan(&id)
	
	if err != nil {
		return false,nil
	}else {
		return true,nil
	}
}

func (proxy *SqlProxy) DeleteData(colname string, table_name string, del_val interface{}) (del_rows int, err error) {
	stmt,e := proxy.db.Prepare("delete from " + table_name + " where " + colname + " = ?")
	if e != nil {
		return 0,e
	}
	
	defer stmt.Close()
	result,e := stmt.Exec(del_val)
	if e != nil{
		return 0,e
	}
	
	rows,_ := result.RowsAffected()
	return int(rows),nil
}

func (proxy* SqlProxy) UpdateData(query_colname string, update_colname string, table_name string,  query_val interface{} ,update_val interface{}) (update_rows int,err error) {
	stmt,e := proxy.db.Prepare("update " + table_name + " set " + update_colname + " = ? " + " where " + query_colname + " = ?")
	if e != nil {
		return 0,e
	}
	
	defer stmt.Close()
	result,e := stmt.Exec(update_val, query_val)
	if e != nil{
		return 0,e
	}
	
	rows,_ := result.RowsAffected()
	return int(rows),nil
}

func (proxy *SqlProxy) ExecSql(sql string) {
	
}

func (proxy *SqlProxy) Close() {
	if proxy.db != nil {
		proxy.db.Close()
	}
}