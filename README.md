# EastMoneySpider

requirements
```shell script
go get -u github.com/mitchellh/mapstructure
go get -u golang.org/x/text
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get github.com/PuerkitoBio/goquery
```

config database
1. install [docker](http://docker.com/)
2. pull mysql: `docker pull mysql`
3. run mysql: `docker run --name eastmoney -p 3306:3306 -e MYSQL_ROOT_PASSWORD=abc -d mysql:latest`
4. change mysql connect dsn from main.go file
```go
dsn := "root:abc@tcp(0.0.0.0:3306)/eastmoney?charset=utf8mb4"
```

build
```shell script
go build .
```

run
```shell script
./EasyMoneySpider
```

