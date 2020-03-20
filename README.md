# Shalog

Shalog是一个基于Golang开源的轻量级内容管理系统，告别PHP类的繁琐的部署方式，超底的资源占用率，并且支持图床功能当然也可以使用七牛作为存储方式，支持Docker、k8s部署。

![](http://source.qiniu.cnd.nsini.com/images/2020/03/b9/d7/37/20200318-de49e256577173333dd85ec0d7fb9dda.jpg?imageView2/2/w/1280/interlace/0/q/70)

## 项目设计

支持内容自定义，支持模版自定义，支持图床功能，支持Mweb，Metaweblog API。

### 内容展示

项目开源地址：[https://github.com/icowan/shalog](github.com/icowan/shalog)

![](http://source.qiniu.cnd.nsini.com/images/2020/03/f5/a8/fb/20200311-4d17f4b35d2fb28cf53ca480a88f57d5.jpg?imageView2/2/w/1280/interlace/0/q/70)

### 管理后台前端

开源地址: [https://github.com/icowan/blog-view](https://github.com/icowan/blog-view)

使用ReaceJS作为管理后台的前端展示，如下图:

![](http://source.qiniu.cnd.nsini.com/images/2020/03/d9/00/ad/20200311-ea32f012517ad52fa9aa953500bd9cf0.jpg?imageView2/2/w/1280/interlace/0/q/70)


## 演示Demo

演示地址: [https://shalog.nsini.com](https://shalog.nsini.com/)

演示管理后台地址: [https://shalog.nsini.com/admin/](https://shalog.nsini.com/admin/)

用户名: `shalog`

密码: `admin`

## 安装说明

平台后端基于[go-kit](https://github.com/go-kit/kit)、前端基于 [umijs](https://umijs.org/) 和 [ant-design](https://github.com/ant-design/ant-design)框架进行开发。

后端所使用到的依赖全部都在[go.mod](go.mod)里，前端的依赖在`package.json`，详情的请看`yarn.lock`，感谢开源社区的贡献。

后端代码: [https://github.com/icowan/shalog](https://github.com/icowan/shalog)

前端代码: [https://github.com/icowan/shalog-view](https://github.com/icowan/shalog-view)

## 快速开始

配置文件准备, **app.cfg**以下为参考:

```ini
[server]
app_name = shalog
app_key = R*9N*Q#ROFJI
debug = false # 是否启用调试模式
log_level = error # warning error info debug
logs_path = /var/log/shalog.log
session_timeout = 14400 # 管理后台登录token失效时间

[mysql]
host = mysql # 数据库地址
port = 3306 # 数据库端口
user = root
password = admin
database = shalog

[cors]
allow = false # 是否支持跨域
origin = http://localhost:8000
methods = GET,POST,OPTIONS,PUT,DELETE
headers = Origin,Content-Type,Authorization,mode,cors,x-requested-with,Access-Control-Allow-Origin,Access-Control-Allow-Credentials
```

### docker-compose 启动

在您的电脑上安装docker-compose命令，请参考: [https://docs.docker.com/compose/install/](https://docs.docker.com/compose/install/)

创建 `docker-compose.yaml` 文件:

```yaml
version: '3'
services:
  mysql:
    image: mysql:5.7.29
    environment:
      MYSQL_ROOT_PASSWORD: "admin"
      MYSQL_DATABASE: "shalog"
    command: [
      '--character-set-server=utf8mb4',
      '--collation-server=utf8mb4_unicode_ci',
    ]
    expose:
      - "3306"
    ports:
      - "3306:3306"
  shalog:
    image: dudulu/shalog:latest
    command: /go/bin/shalog start -p :8080 -c /etc/shalog/app.cfg
    environment:
      GOPATH: "/go"
      USERNAME: "shalog"
      PASSWORD: "admin"
      SQL_PATH: ./database/db.sql
    volumes:
      - ./app.cfg:/etc/shalog/app.cfg
    depends_on:
      - mysql
    restart: always
    ports:
      - "8080:8080"
```

将上面准备好的app.cfg放到当前目录，然后执行以下命令:

```
$ docker-compose start
```

浏览器输入: `http://localhost:8080` 访问

### 本地启动

- Golang 1.13+ [安装手册](https://golang.org/dl/)
- MySQL 5.7+ (大多数据都存在mysql)


修改 `app.cfg` 文件，将mysql地址配置为您自己的数据库地址。

克隆代码，及本地启动

```
$ git clone github.com/icowan/shalog.git
$ cd shalog/
$ make run
```

浏览器输入: `http://localhost:8080` 访问

## 文档

- [内容发布编辑](https://www.lattecake.com/post/20135)
- [图片处理](https://lattecake.com/post/20132)
- [站点设置和更换模版](https://www.lattecake.com/post/20136)
- [友链申请审核](https://lattecake.com/post/20131)
- [通过mweb发布](https://lattecake.com/post/20134)
- [开普勒云平台部署](https://lattecake.com/post/20133)
- 更多内容请访问[https://lattecake.com/search?tag=Shalog](https://lattecake.com/search?tag=Shalog)

## 支持我

![](http://source.qiniu.cnd.nsini.com//static/pay/wechat-pay.JPG?imageView2/2/w/360/interlace/0/q/70)

![](http://source.qiniu.cnd.nsini.com//static/pay/alipay.JPG?imageView2/2/w/360/interlace/0/q/70)
