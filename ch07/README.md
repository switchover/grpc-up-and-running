# Chapter 7 : 서비스 수준 gRPC 실행

## 예제 코드 리스트
- 코드 7-1 (서버측 테스트) : [main.go](01-Testing/productinfo/server/main_test.go)

## 정오
### 코드 부분
- 217 페이지 ghz 실행 : IP(`0.0.0.0:50051`) 부분  
    ```
    0.0.0.0:50051
    ```
    :arrow_right:
    ```
    localhost:50051
    ```
- 220 페이지 docker image build 실행 : 뒤 부분 경로 추가 필요  
    ```
    docker image build -t grpc-productinfo-server -f server/Dockerfile
    ```
    :arrow_right:
    ```
    docker image build -t grpc-productinfo-server -f server/Dockerfile .
    ```
---
# 세부 세션별 예제

* 테스팅 : [Testing](./01-Testing)
* 도커 배포 : [Docoker](./02-Docker)

---
# 최종 코드

gRPC의 CI/CD, 모니터링 등에 대한 예제 코드는 원서의 소스 저장소 [7장 부분](https://github.com/grpc-up-and-running/samples/tree/master/ch07)을 참고합니다.
