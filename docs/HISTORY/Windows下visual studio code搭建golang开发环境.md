# Windows下visual studio code搭建golang开发环境
## 序幕

其实环境搭建没什么难的，但是遇到一些问题，主要是有些网站资源访问不了(如：golang.org)，
导致一些包无法安装，最终会导致环境搭建失败，跟据这个教程几步，我们将可以快速的构建golang的开发环境。

### 开发环境：
![image](http://images.cnblogs.com/cnblogs_com/JerryNo1/818380/o_QQ%e6%88%aa%e5%9b%be20160420134611.jpg)


## 一、安装
这里我用需要安装一些工具：

1.[Visual Studio Code 1.0.0](https://www.visualstudio.com/products/code-vs)

2.[Golang下载](http://www.golangtc.com/download) 
```
这里我使用的是Go1.6.
```

3.[git下载](https://git-scm.com/download/) 
```
这一步跟建环境没什么关系，
但是之后要引用一些包需要它.
```

4.[安装vscode-go 插件](https://marketplace.visualstudio.com/items?itemName=lukehoban.Go)

![image](http://images.cnblogs.com/cnblogs_com/JerryNo1/818380/o_QQ%e6%88%aa%e5%9b%be20160420141250.jpg)

```
安装方法:打开vscode 按F1 然后选择安装扩展 选择go 下载后 重启
```
![image](http://images.cnblogs.com/cnblogs_com/JerryNo1/818380/o_QQ%e6%88%aa%e5%9b%be20160420141705.jpg)

5.[下载vscode扩展依赖工具](http://pan.baidu.com/s/1qY7pRa8)

```
vscode-go 插件需要一些工具，这些功具默认需要这样安装：

go get -u -v github.com/nsf/gocode
go get -u -v github.com/rogpeppe/godef
go get -u -v github.com/golang/lint/golint
go get -u -v github.com/lukehoban/go-outline
go get -u -v sourcegraph.com/sqs/goreturns
go get -u -v golang.org/x/tools/cmd/gorename
go get -u -v github.com/tpng/gopkgs
go get -u -v github.com/newhook/go-symbols
go get -u -v golang.org/x/tools/cmd/guru
(以及dlv调试运行工具)

这些内容最终目的是要在gopath中bin文件夹中生成一些exe文件。
但是，这里安装时遇到了一些问题，有一些会安装失败，报错信息主要意思是说无法链接啊，无法找到啊之类的。
这里我省略这些步骤。直接下载我已经生成的exe包，放在gopath的bin文件夹中就行了。如下图：
```
![image](http://images.cnblogs.com/cnblogs_com/JerryNo1/818380/o_QQ%e6%88%aa%e5%9b%be20160420142649.jpg)

## 二、配置

1.GOPATH：这个是你的工作目录，请在系统新加一个GOPATH指定你的工作目录,并且在工作目录中新建
bin,pkg,src三个文件夹。

2.Git：将git安装目录下的bin目录加入Path环境变量。

3.用vscode打开src文件夹添加两个文件如下：
![image](http://images.cnblogs.com/cnblogs_com/JerryNo1/818380/o_QQ%e6%88%aa%e5%9b%be20160420143835.jpg)
```
{
    "files.autoSave": "onFocusChange",
    "go.buildOnSave": true,
    "go.lintOnSave": true,
    "go.vetOnSave": true,
    "go.buildTags": "",
    "go.buildFlags": [],
    "go.lintFlags": [],
    "go.vetFlags": [],
    "go.coverOnSave": false,
    "go.useCodeSnippetsOnFunctionSuggest": false,
    "go.formatOnSave": true,
    "go.formatTool": "goreturns",
    "go.goroot": "C:\\Go",
    "go.gopath": "H:\\Projects\\GitHub\\Go",
    "go.gocodeAutoBuild": true
}
```
![image](http://images.cnblogs.com/cnblogs_com/JerryNo1/818380/o_QQ%e6%88%aa%e5%9b%be20160420144012.jpg)
```
{
	"version": "0.2.0",
	"configurations": [
		{
			"name": "Launch",
			"type": "go",
			"request": "launch",
			"mode": "debug",
			"remotePath": "",
			"port": 2345,
			"host": "127.0.0.1",
			"program": "${workspaceRoot}\\helloworld",
			"env": {},
			"args": []
		}
	]
}
```
## 三、演示

##运行
![image](http://images.cnblogs.com/cnblogs_com/JerryNo1/818380/o_QQ%e6%88%aa%e5%9b%be20160420151454.jpg)
##调试
![image](http://images.cnblogs.com/cnblogs_com/JerryNo1/818380/o_QQ%e6%88%aa%e5%9b%be20160420151303.jpg)

# 四、总结
这些库因为不能获取，我在github上找了好久，最后编译成exe,国内找不到几篇文章提到这个问题，
所以写个博客，让大家方便的开发。希望大家快乐的学习go语言。