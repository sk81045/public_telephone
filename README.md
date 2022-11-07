## 基于TCP的公共电话SERVER端

## 多平台交叉编译打包

# Linux
# X86
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
 
# ARM
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build main.go

# Windows
# X86
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
 
# ARM
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build main.go

# MacOS
# X86
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go
 
# ARM
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build main.go


### 4.3 目录结构

```
     。。。。。      

```

