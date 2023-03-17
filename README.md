# 文件同步器

<p>用于扫描指定目录中的静态文件并将其同步到另一台计算机上的相对目录中。</p>


## install
`go get -u github.com/Madou-Shinni/file-rsync`

## 快速开始


注意：windows盘符目录需要修改 D: -> /cygdrive/D
目前支持 windows -> linux
        linux -> linux
```shell
GLOBAL OPTIONS:
   --ip value, -I value      远程主机的ip
   --src value, -S value     需要同步文件的目录
   --dist value, -D value    需要同步文件的目标目录，default:src目录
   --module value, -M value  rsync的模组(需要提前在配置文件中创建)，default:test (default: test)
   --help, -h                show help


go run main.go --ip 192.168.110.94 --src /cygdrive/D/go-project/frisbee-officer-backend-GVA/server/uploads/file/ --dist /root/rsfile --module test
go run main.go -I 192.168.110.94 -S /cygdrive/D/go-project/frisbee-officer-backend-GVA/server/uploads/file/ -D /root/rsfile -M test
```

## rsync下载地址
<h3>Introduction</h3>
本机采用的是windows10

安装可参考：https://blog.csdn.net/qq_42684504/article/details/105433988

服务端：https://www.backupassist.com/rsync/
客户端：https://www.itefix.net/

rsyncd.conf是启动的配置文件

运行rsync server `rsync.exe --config=../../rsyncd.conf --daemon --no-detach`

运行client端 `>rsync -avz --progress /cygdrive/D/go-project/frisbee-officer-backend-GVA/server/uploads/file/ 192.168.110.94::test`

注意：D:需要使用 /cygdrive/D 替换

## 使用参考
https://www.ruanyifeng.com/blog/2020/08/rsync.html

## Linux Server
<h3>Introduction</h3>
linux自带了rsync，你做要做的就是修改配置文件
`vim /etc/rsyncd.conf`

```conf
uid = 0
gid = 0
use chroot = false
strict modes = false
hosts allow = *
log file = rsyncd.log

# Module definitions
# Remember cygwin naming conventions : c:\work becomes /cygwin/c/work
#
[test] # 模组
path = / # 只需要修改这里，将你同步的文件保存在这里
read only = false
transfer logging = yes
```
<h3>run</h3>
```shell
rsync /etc/rsyncd.conf --daemon
```