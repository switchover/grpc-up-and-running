# Chapter 7 : 서비스 수준 gRPC 실행

## 예제 코드 리스트
- 코드 7-1 (서버측 테스트) : [main.go](01-Testing/productinfo/server/main_test.go)
- 코드 7-2 (클라이언트 mock 테스트) : [prodinfo_mock_test.go](01-Testing/productinfo/client/mock_prodinfo/prodinfo_mock_test.go)
- 코드 7-4, 7-5 (서버용 쿠버네티스 배포 및 서비스 기술자) : [grpc-prodinfo-server.yaml](03-Kubernetes/productinfo/server/grpc-prodinfo-server.yaml)
- 코드 7-6 (클라이언트용 쿠버네티스 잡 기술자) : [grpc-prodinfo-client-job.yaml](03-Kubernetes/productinfo/client/grpc-prodinfo-client-job.yaml)
- 코드 7-7 (서비스용 쿠버네티스 인그레스 기술자) : [grpc-prodinfo-ingress.yaml](03-Kubernetes/productinfo/ingress/grpc-prodinfo-ingress.yaml)

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
* 쿠버네티스 배포 : [Kubernetes](./03-Kubernetes)
* 오픈센서스 활용 : TBD (원서 저장소 참조 : https://github.com/grpc-up-and-running/samples/tree/master/ch07/grpc-opencensus, https://github.com/grpc-up-and-running/samples/tree/master/ch07/grpc-opencensus-tracing)
* 프로메테우스 활용 : TBD (원서 저장소 참조 : https://github.com/grpc-up-and-running/samples/tree/master/ch07/grpc-prometheus/go)

---
# 최종 코드

gRPC의 CI/CD, 모니터링 등에 대한 예제 코드는 원서의 소스 저장소 [7장 부분](https://github.com/grpc-up-and-running/samples/tree/master/ch07)을 참고합니다.
