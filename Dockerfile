FROM golang:alpine
ENV TZ Asia/Shanghai
RUN echo "http://mirrors.aliyun.com/alpine/v3.4/main/" > /etc/apk/repositories \
        && apk --no-cache add tzdata zeromq \
        && ln -snf /usr/share/zoneinfo/${TZ} /etc/localtime \
        && echo "${TZ}" > /etc/timezone
COPY . /$GOPATH/src/hello_go/
WORKDIR /$GOPATH/src/hello_go/
#设置环境变量，开启go module和设置下载代理
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
#会在当前目录生成一个go.mod文件用于包管理
RUN go mod init
#增加缺失的包，移除没用的包
RUN go mod tidy
RUN go build src/main.go
CMD ["go","run","src/main.go"]