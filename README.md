# example

大连民族大学官方微信公众号前后端代码，涉及技术为GIN、Gorm、微信公众号开发

## 文档
[wechat sdk文档](https://silenceper.com/wechat)

## Build 
```go
  注意在linux环境下打包，windows下打包会有godror包引用错误
  go build -o bin/example cmd/example/main.go
```
## Run
请修改config.yaml中的相关为自己的参数再运行
```go
./bin/example
```

## 关注公众号

![学点程序](https://silenceper.oss-cn-beijing.aliyuncs.com/qrcode/search_study_program.png)