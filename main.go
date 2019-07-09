package main

import (
	"database/sql"
	"github.com/BlackCarDriver/log"
	"github.com/BlackCarDriver/config"
	_ "github.com/lib/pq"
	"fmt"
	"strings"
)

var(
	db *sql.DB
	goRes *log.Logger
	tsRes *log.Logger
	conf *config.ConfigMachine
	schemeName string
	queryResult map[string]column
) 

type column map[string]string	

type dbconfig struct {
	Host string  `json:"host"`
	Port int	`json:"port"`
	UserName  string `json:"username"`
	DbName	string	`json:"dbname"`
	Password string 	`json:"password"`
}

var template string = `
SELECT 
tb.tablename as tablename,
a.attname AS columnname,
t.typname AS type
FROM
pg_class as c,
pg_attribute as a, 
pg_type as t,
(select tablename from pg_tables where schemaname = $1 ) as tb
WHERE  a.attnum > 0 
and a.attrelid = c.oid
and a.atttypid = t.oid
and c.relname = tb.tablename 
order by tablename
` 

var tconf = dbconfig{}


//postgresql type -> go type
var pgMap = map[string]string{
	"int4":"int32", 
	"int8":"int64",
	"float4":"float32",
	"float8":"float64",
	"double":"float64",
	"varchar":"string",
	"boolean":"bool",
	"timestamp":"time.Time",
	"date":"time.Time",
}

//postgresql type -> typeScript type
var ptMap = map[string]string{
	"int4":"number", 
	"int8":"number",
	"float4":"number",
	"float8":"number",
	"double":"number",
	"varchar":"string",
	"boolean":"boolean",
	"timestamp":"string",
	"date":"string",
}

func init() {
	queryResult = make(map[string]column)

	//set the file path that result save in
	log.SetLogPath("./")
	goRes = log.NewLogger("Go.txt")
	tsRes = log.NewLogger("TypeScript.txt")
	
	//get database connect config from config file
	conf,err := config.NewConfig("./")
	if err != nil {
		fmt.Println("无法获取配置文件!!")
		panic(err)
	}
	conf.SetIsStrict(true)
	schemeName, err = conf.GetString("schema")
	conf.GetStruct("database", &tconf)

	//connect to the database
	confstr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		tconf.Host,
		tconf.Port,
		tconf.UserName,
		tconf.Password,
		tconf.DbName,
	)
	db, err = sql.Open("postgres", confstr)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Database Connect Scuess!")
}

func QueryTables(){
	rows, err := db.Query(template, schemeName)
	if err!=nil {
		panic(err)
	}
	for rows.Next() {
		tableName, colName, colType := "", "", ""
		rows.Scan(&tableName, &colName, &colType)
		_, ok := queryResult[tableName]
		if !ok {
			queryResult[tableName] = make( map[string]string )
		}
		queryResult[tableName][colName] = colType
	}
}

func Report(){
	for tableName, ctMap := range queryResult {
		goRes.Write("\ntype %s struct {\n", tableName)
		tsRes.Write("\nexport class %s {\n", tableName)
		for colName, coltype := range ctMap {
			//translate to go struct foramt
			retype := pgMap[coltype]
			if retype==""{
				retype = coltype
			}
			jsonName := strings.ToLower(colName)
			vname := strings.ToUpper(jsonName[0:1]) + jsonName[1:]
			goRes.Write("\t%-10s %-10s `json:\"%s\"`\n", vname, retype, jsonName)
			//translate to typescript format
			retype = ptMap[coltype]
			if retype==""{
				retype = coltype
			}
			tsRes.Write("\t%s:%s; \n", jsonName, retype)
		}
		goRes.Write("}\n")
		tsRes.Write("}\n") 
	}
}

func main(){
	QueryTables()
	Report()
}
