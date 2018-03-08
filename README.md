### 简介
gopass是一个跨平台的密码管理工具，该工具旨在解决linux终端下使用密码不方便的问题。  

* 直接终端操作，无须在各个应用之间切换，比较适用于重度使用终端的同学。  
* 支持在终端以图形界面的方式搜索匹配，复制粘贴，简化了上手难度、提高了工作效率  


### 快速开始
* 安装xclip, linux下依赖xclip复制到剪贴板
 
```Go
CentOS: sudo yum -y install xclip
Ubuntu: sudo apt -y install xclip
Mac:    
```

* 将配置文件和gopass二进制文件放在同一目录下，直接运行gopass -h查看用法    


```Go
NAME:
   gopass - A tool for managing your password in terminal

USAGE:
   gopass [command]

COMMANDS:
     add      add password into gopass
     del      delete password from gopass
     update   update password into gopass
     ui       display ui to search and copy password
     help     show help
     version  show version
```

## Demo
[![asciicast](https://asciinema.org/a/0fPa5CJzUiue5Ilt1aJyd0I1x.png)](https://asciinema.org/a/0fPa5CJzUiue5Ilt1aJyd0I1x)

### UI快捷键
|按键|功能|
|----|---|
|小写q|退出UI|  
|Enter|复制对应的密码内容|
|Ctrl+Q|清楚搜索框的内容|


