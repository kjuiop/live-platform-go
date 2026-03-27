# live-platform-go

---

`live-platform-go`는 **실시간 채팅 시스템**을 위한 서버 애플리케이션으로, Kafka와 WebSocket을 활용하여 확장성 있는 채팅 서비스 구축을 목표로 합니다.


## 🧩 시스템 구성

![image](https://github.com/user-attachments/assets/44689192-f005-407d-9cd3-d9f52440683e)

- 사용자는 HTTP API를 통해 채팅방을 생성하거나 삭제할 수 있습니다.
- 채팅방이 생성될 때에 Kafka 의 Topic 도 생성됩니다.
- 채팅방에 사용자가 입장하게 되면 Kafka 에 등록된 채팅방을 Subscribe 하게 됩니다.
- 서버가 사용자로부터 WebSocket 메시지를 받으면 이를 Kafka 브로커에 Publish 합니다.
- Kafka 로부터 메시지를 서버가 전달받게 되면 서버와 연결된 접속자들에게 메시지를 전달합니다.


## ⚙️ 기술 스택

- **언어**: Go 1.23
- **웹 프레임워크**: Gin
- **WebSocket**: Gorilla WebSocket
- **데이터 저장소**: Redis
- **빌드 도구**: Makefile

## 🚀 시작하기


### TEST

```shell
make test
```


### BUILD

```shell
make build
```


### RUN
```shell
./live-chat-server
```

<br />


## 📄 API Spec Document


### ws join

```shell
ws://localhost:8091/ws/chat/join/rooms/N1-01JSVD2N05RD0F4GPGDHR5C73J/user/jake
```


### ws chat

```shell
{
    "Method":"chat",
    "SendUserId": "jake",
    "Message": "hello"
}
```