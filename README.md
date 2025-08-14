# fast2fa
# 将大量linux配置同一个2fa

## 前提
**需要一台管理机和目标主机**
***需要联网安装必备组件(编译安装google-authenticator.zip)**
***配置好时间同步！！！***
后期随缘更新（也许）

```bash
 yum install -y unzip git gcc make autoconf libtool pam-devel automake
 apt-get update && apt-get install -y unzip git build-essential autoconf libtool automake libpam0g-dev 
```

自行编辑hosts.csv写上目的主机

## 使用方法

```bash
git clone  https://github.com/jdjn123/fast2fa.git
chmod 777 fast2fa
./fast2fa --hosts hosts.csv 

```

日志功能还没做
***生产环境慎用，使用前，先做好telnet或者其他配置。***