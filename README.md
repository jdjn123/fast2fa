# fast2fa
# 将大量linux配置同一个2fa

## 前提
** 需要一台管理机和目标主机 **
** 需要联网安装必备组件(编译安装google-authenticator.zip)**
后期随缘更新（也许）

```bash
 yum install -y unzip git gcc make autoconf libtool pam-devel automake
 apt-get update && apt-get install -y unzip git build-essential autoconf libtool automake libpam0g-dev 
```



## 使用方法

```bash
git clone  https://github.com/jdjn123/fast2fa.git
chmod 777 fast2fa

```
