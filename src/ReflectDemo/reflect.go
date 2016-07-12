package main

import (
	"fmt"
	"reflect"
)

type Teacher struct {
	Name string
	Age int
}

type Student struct {
	Name string
	Age int
	Points [] int
	Tea Teacher
}

func show (args...interface{}) {
	for _,v := range args {
		fmt.Println(v)
	}
}


func test(data interface{}) {
	
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
	
	show(data_values...)
}

func main() {
	
	st := Student {
		Name : "wang",
		Age : 11,
		Points : []int {1,2,3,4},
		Tea : Teacher {Name:"Li", Age:30},
	}
	
	
	fmt.Println("type:",reflect.TypeOf(st))
	s := reflect.ValueOf(&st).Elem()
	t := s.Type()
	
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%s %s = %v\n", t.Field(i).Name, f.Type(), f.Interface())
		if f.Kind() == reflect.String {
			fmt.Printf("struct type: %s\n", t.Field(i).Name)
		}
	}
	
	teacher := Teacher {
		Name : "LiLi",
		Age : 40,
	}
	
	test(&teacher)
}