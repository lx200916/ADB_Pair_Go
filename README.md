# ADB_Pair_Go
Try to Make WiFi-ADB Easier.(Android 11+ Only.
> [!IMPORTANT]
> 本程序支持Sixel协议,在支持的终端中可以以像素形式渲染二维码,而在不支持的终端中只能回落到字符画格式.Sixel渲染的二维码清晰度更高,相对更易扫描.请尽量使用如iTerm2[macOS] konsole[KDE] 或者VSCode[开启 "terminal.integrated.experimentalImageSupport": true] 等终端模拟器运行.
> 
模仿Android Studio中插件的行为编写的WiFi ADB配对工具.在开发者设置中开启`无线调试`后,使用`通过二维码配对`选项扫描二维码即可完成配对.

配对原理是查看 AS 2020.3.1 Canary Preview 版本的源码得来,在博客中有详细介绍:https://saltedfish.fun/index.php/archives/29/

需要安装ADB并加入到环境变量.
