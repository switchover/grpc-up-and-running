# Chapter 6 : 보안 적용 gRPC

## 예제 코드 리스트
- 코드 6-1 (서버 단방향 보안 연결) : [main.go](01-SecureChannel/productinfo/server/main.go)
- 코드 6-2 (클라이언트 단방향 보안 연결) : [main.go](01-SecureChannel/productinfo/client/main.go)
- 코드 6-3 (서버 mTLS 보안 연결) : [main.go](02-mTLS/productinfo/server/main.go)
- 코드 6-4 (클라이언트 mTLS 보안 연결) : [main.go](02-mTLS/productinfo/client/main.go)

## 정오
### 코드 부분
- 185 페이지 코드 6-1. `grpc.ServerOption` 구조체 리터럴 부분 : 뒤에 `,` 추가 필요 (여러 줄을 사용한 경우)
    ```
    opts := []grpc.ServerOption{
        grpc.Creds(credentials.NewServerTLSFromCert(&cert))
    }
    ```
    :arrow_right:
    ```
    opts := []grpc.ServerOption{
        grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
    }
    ```
    또는
    ```
    opts := []grpc.ServerOption{grpc.Creds(credentials.NewServerTLSFromCert(&cert))}
    ```
- 187 페이지 코드 6-2. `var` 선언 중 `hostname` 지정 부분 : 뒤에 `"` 누락
    ```
    address = "localhost:50051
    ```
    :arrow_right:
    ```
    address = "localhost:50051"
    ```

---
# 보안 인증서 만들기
6장의 예제들을 실행하기 위해서는 SSL/TLS용 인증서 등이 필요하며, 이를 위해 [OpenSSL](https://www.openssl.org/)을 사용합니다. (프로그램 설치 필요)

따라서 [Certificates](./00-Certificates) 부분을 통해 먼저 예제 실행을 위한 인증서를 먼저 생성합니다.

---
# 세부 세션별 예제

* 인증서 만들기 : [Certificates](./00-Certificates) (6장 공통 사용 인증서)
* 보안 설정 : [Secure Channel](./01-SecureChannel)
* mTLS(mutual TLS) 설정 : [mTLS](./02-mTLS)


---
# 최종 코드

gRPC의 보안 기능에 대한 예제 코드는 원서의 소스 저장소 [6장 부분](https://github.com/grpc-up-and-running/samples/tree/master/ch06)을 참고합니다.
