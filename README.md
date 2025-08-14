# fast2fa

**批量为多台 Linux 主机配置同一个 2FA（手机令牌）**  
需要用到的代码https://github.com/google/google-authenticator

---


## 前提条件

- 一台 **管理机**（运行本工具）  
- 多台 **目标主机**（需要开启 SSH 登录，并可联网安装依赖）  
- 需要联网以安装必备组件（首次运行时编译安装 `google-authenticator`）  
- 后续可能会随缘更新

---

## 必备依赖

管理机和目标机均需安装以下组件：

**CentOS / RHEL 系列**
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

日志功能还没做（以后会做）

***生产环境慎用，使用前，先做好telnet或者其他ssh配置。***