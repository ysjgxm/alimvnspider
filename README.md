# alimvnspider

获取maven阿里源全量仓库

使用方式： 编译： D:\downloads\gopath\src\alimvnspider>go build .

运行： 第一个参数是要下载的目录URL（结尾要有“/”）， 第二个参数是要存储的路径（绝对路径，结尾要有“/”） D:\downloads\gopath\src\alimvnspider>alimvnspider.exe "https://maven.aliyun.com/browse/tree?_input_c harset=utf-8&repoId=central&path=activemq/" "D:/alitest/"

注意： 1 在win10上测试阿里源的activemq目录成功下载所有目录结构和文件。在linux如果有问题，请帮助完善该程序，多谢~ 2 使用过程如果出现下载文件的请求超时、远端拒绝，请尝试控制下载并发数，并完善该程序，多谢~（目前通过随机休眠，使下载文件时间小于访问目录时间，简单控制下载并发数）
