# postgre-struct-maker

----------

### 简介
这是一个帮助go项目开发的小工具，能够通过连接数据库，自动得到数据库中已经建立的所有表格的字段名和类型并生成结构体的代码；优点是简单易用，扩展性搞，支持从配置文件读取数据库配置，并将生成结果保存到不同的文件中。目前支持postgres数据库,并且能够生成golang和typescript 的结构体。

### 使用方法
修改 main.exe 同级目录下的config.conf文件，点击运行main.exe或go run main.go，即可在该目录生成go.txt 和typescript.txt 分别保存生成的结构体代码。

### 配置格式

    # 数据库连接参数,json格式，字段名勿改
    database =  {
    "host":"localhost",
    "port":5432,
    "username":"testuser",
    "dbname":"testdb",
    "password":"testpassword"
    }
    # 模式名，将会查询这个模式下的所有表的字段名和类型，从而生成结构体   
    schema = "public"
    

----------
2019/7/9 16:31:22 