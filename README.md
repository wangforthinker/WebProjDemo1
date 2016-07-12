# WebProjDemo1
简介：
1、简单的web系统，使用mysql存储数据，简化sql 增删改查使用方式（通过使用SqlProxy代理），其中使用了反射机制；
2、本系统只是简单的一个登陆系统，可以申请用户，修改用户密码，删除用户和登陆用户，后台连接mysql存储；
3、其中，在添加用户时，采用了反射机制，使得向数据库添加记录时可以泛化至任意数据表，不用每次都自己重复写sql语句，详情参见SqlProxy.InsertData方法。

使用方法：
下载后，在src/WebClient和src/WebServer 目录分别执行 go build -v,生成可执行文件
服务端运行：./WebServer port database_addr database_username database_passwd database_name
客户端运行： ./WebClient methods_name ...  
客户端不同方法不同参数，详情见源代码

