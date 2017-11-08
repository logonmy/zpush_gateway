# zpush长连接推送—接入网关
### 架构图

![架构图](http://on-img.com/chart_image/5a033e63e4b0f84f8978b4b7.png)

### 接入网关Gateway的职责

- 负责保持和客户端的TCP长连接
- 可同时启动多个Gateway Server，客户端通过负载均衡策略连接到任意一台Gateway上
- 接收并解析客户端的协议包（TLV格式：2字节的命令码 + 4字节的包体长度 + 可变长度的包体（包体采用protobuf序列化））
- 使用gRPC将解析之后的业务报文转发给Message Server处理
- 收到客户端的登录报文后，通过HTTP请求业务鉴权服务器对用户进行鉴权，鉴权成功才建立连接

