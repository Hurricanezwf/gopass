* [简介](#简介)
* [演示](#演示)
* [快速开始](#快速开始)
* [快捷键](#快捷键)


### 简介
gopass是一个跨平台的密码管理工具，该工具旨在解决linux终端下使用密码不方便的问题。  

* 直接终端操作，无须在各个应用之间切换，比较适用于重度使用终端的同学。  
* 支持在终端以图形界面的方式搜索匹配，复制粘贴，简化了上手难度、提高了工作效率。
* 支持自定义预留SecretKey加解密数据，支持修改预留SecretKey，更安全。


### 演示
[![asciicast](https://asciinema.org/a/0fPa5CJzUiue5Ilt1aJyd0I1x.png)](https://asciinema.org/a/0fPa5CJzUiue5Ilt1aJyd0I1x)


### 快速开始
* 安装xclip, linux下依赖xclip复制到剪贴板
 
```Go
CentOS: sudo yum -y install xclip
Ubuntu: sudo apt -y install xclip
Mac:    依赖于pbcopy,已默认安装
```

* 将配置文件和gopass二进制文件放在同一目录下(也可用-c指定配置路径)，直接运行gopass -h查看用法    


```Go
NAME:
   gopass - A tool for managing your password in terminal

USAGE:
   gopass [command] [-c ConfigFile]

COMMANDS:
     add      add password into gopass
     del      delete password from gopass
     update   update password into gopass
     ui       display ui to search and copy password
     chsk     change auth SecretKey which you provide for authentication when init the app
     help     show help
     version  show version
```


### 快捷键
gopass ui支持如下快捷键：

|按键|功能|
|----|---|
|小写q|退出UI|  
|Enter|复制对应的密码内容|
|Ctrl+Q|清除搜索框的内容|
