# Scoop使用文档(windows)

<!-- Scoop使用文档,安装以及开发工具安装 [2024-02-01] -->

## 1.安装Scoop (可能需要梯子) 
```ps
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression
```

## 2.添加bucket extras
```shell
scoop bucket add extras
```

## 3.应用
```shell
# 搜索app https://scoop.sh/#/apps

# vscode
scoop install extras/vscode

# golang
scoop install main/go

# bun
scoop install main/bun

# chsrc换源工具
scoop install main/chsrc

# 7-zip
scoop install main/7zip

```