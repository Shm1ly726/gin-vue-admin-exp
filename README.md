# gin-vue-admin-exp
## 简介
学习golang编程练习写的一款简单的gin-vue-admin框架利用小工具

能力有限，代码写得很垃圾，各位师傅们轻喷，也欢迎各位师傅指导与讨论

**声明：此工具仅作学习交流用，不要做任何违法行为，如果违法&恶意操作，与本人无关！！！**

## 主要功能
相关漏洞原理：[http://47.109.35.4/posts/src-gin-vue-admin%E6%94%BB%E5%87%BB%E9%9D%A2/](http://47.109.35.4/posts/src-gin-vue-admin%E6%94%BB%E5%87%BB%E9%9D%A2/)

1. 前端绕过
   1. 绕过思路1：只判断`captchaId`与`captcha`是否为空，不校验验证码是否正确
   2. 绕过思路2：只判断`captchaId`与`captcha`是否为空，不校验验证码是否正确
2. 绕过后撞库
   1. Username：内置字典，默认`"admin", "test", "demo", "guest", "test1", "user", "ceshi", "test123", "system", "web", "sys", "admin1"`
   2. password：加载根目录`bigpasswdDict.txt`文件
3. 漏洞利用
   1. 提权/后门利用(CNVD-2024-00979)
   2. 读取系统配置利用
## 使用说明
```go
gin-vue-admin-exp.exe -u https://www.baidu.com (判断前端是否绕过)
gin-vue-admin-exp.exe -u https://www.baidu.com -x xxxx (撞库后获取x-token进行漏洞利用)
```
编译命令
```go
go mod init
go mod tidy
go build -o gin-vue-admin-exp.exe
```
## 运行截图
```
gin-vue-admin-exp.exe -u https://www.baidu.com
```
![image.png](https://cdn.nlark.com/yuque/0/2024/png/29227563/1705476662529-bffba752-f7e0-4db1-977e-86caace9500d.png#averageHue=%231c1c1c&clientId=u7ff73711-4205-4&from=paste&height=333&id=u62eb13b7&originHeight=416&originWidth=1341&originalType=binary&ratio=1.25&rotation=0&showTitle=false&size=67311&status=done&style=none&taskId=ue15461db-f553-4447-b24e-8f147b710b2&title=&width=1072.8)

```
gin-vue-admin-exp.exe -u https://www.baidu.com -x xxxx
```
![image.png](https://cdn.nlark.com/yuque/0/2024/png/29227563/1705476694652-20534861-f57f-4e0f-93e0-129abfc96a97.png#averageHue=%23161616&clientId=u7ff73711-4205-4&from=paste&height=572&id=u91987a6e&originHeight=715&originWidth=1360&originalType=binary&ratio=1.25&rotation=0&showTitle=false&size=72828&status=done&style=none&taskId=u110d563d-0906-4235-8a2b-a60deec5025&title=&width=1088)

## 计划
1. 添加第二种撞库模式，引入字典dict，固定密码撞库(实现不难，但是太懒了)
