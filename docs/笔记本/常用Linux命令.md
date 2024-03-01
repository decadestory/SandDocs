# 常用Linux命令

## 防火墙相关
```shell
开启防火墙：
systemctl start firewalld

关闭防火墙：
systemctl stop firewalld

查看防火墙状态：
systemctl status firewalld 
firewall-cmd --state
```

## 开放/关闭端口
```shell
添加端口
firewall-cmd --zone=public --add-port=9092/tcp --permanent && firewall-cmd --reload

删除端口
firewall-cmd --zone=public --remove-port=5905/tcp --permanent 

端口配置生效
firewall-cmd --reload

查看开放端口
firewall-cmd --zone=public --list-ports
```
## 开放/关闭端口
```shell
查看监听的端口
netstat -lnpt

查看进程网络连接
lsof -p 23501

查看网络：nmcli d
启用网络：nmtui 
```
## 文件系统
```shell
查看文件、文件夹大小
du -sh *

磁盘信息
df -hl 
```