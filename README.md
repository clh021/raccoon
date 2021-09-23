# steps

#### 介绍
项目启动步骤配置工具

最初是为了打包 lithium 项目的 appImage，提供静态文件web服务后启动浏览器中间件而制作。为了适应更多变的项目部署场景而衍生出灵活配置的需求……


#### 使用说明

```bash
make build # 编译
make buildcross # 交叉编译 mips64le, arm64, amd64 版本
make run # 运行编译好的程序
make serve # 编译并运行
```

#### 配置
示例
```yaml
  -
    tag: background
    exec: nitr
  -
    tag: web # 建议 tag 用英文，会作为环境变量的一部分传递给后面的执行步骤
    webroot: dist
    addr: :0 # `:0` 表示自动寻找可用端口，无论是否自动，后续步骤都可通过 `APP_${upper(tag)}_URL` 得到服务地址
  -
    tag: client
    exec: start.sh
```

注意事项:

1. webroot 优先级最高，设置此项则认定为该步骤仅需提供 web 服务，忽略其他参数
2. webroot 应该是一个准备工作而非最终步骤，所以不应设置在最后一项
3. webroot 设置后，后面的步骤都将可以通过环境变量 `APP_${upper(tag)}_URL` 获取到 webroot 的监听地址 `127.0.0.1:9999` 目标是传递随机设置的端口
4. 当最后一个步骤退出时，前面启动的所有程序和服务都将退出。