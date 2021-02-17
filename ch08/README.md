# # Chapter 8 : gRPC 생태계

# 참고 사항
## gRPC Gateway 버전 변경
기존 v1에서 v2로 변경(20년 10월)됨에 따라 일부 패키지 등이 변경되었습니다.
자세한 변경 사항은 객 세션별 예제 부분을 참조해 주세요.
(단, 코드 상에서 `github.com/grpc-up-and-running/samples`를 사용하는 경우 v1.16 사용되어, 코드 활용에는 문제가 없음)
- 참고 : https://github.com/grpc-ecosystem/grpc-gateway/blob/master/docs/docs/development/grpc-gateway_v2_migration_guide.md 

## 정오
### 코드 부분
- 257 페이지 맨 아래 `curl` 호출 부분 : 마지막 라인은 명령이 아닌 출력 내용으로 마지막 `\` 부분 삭제
    ```shell
    $ curl -X POST http://localhost:8081/v1/product \
    -d '{"name": "Apple", "description": "iphone7", "price": 699}' \
    "38e13578-d91e-11e9"
    ```
    :arrow_right:
    ```
    $ curl -X POST http://localhost:8081/v1/product \
    -d '{"name": "Apple", "description": "iphone7", "price": 699}'
    ```
    - 마지막 `"38e13578-d91e-11e9"` 라인은 출력된 결과로 생성된 제품 ID입니다.

---
# 세부 세션별 예제
* gRPC 게이트웨이 : [gRPC Gateway](./01-Gateway)
* 서버 리플렉션 : TBD (원서 저장소 참조 : https://github.com/grpc-up-and-running/samples/tree/master/ch08/server-reflection)

---
# 최종 코드

gRPC 생태계에 대한 예제 코드는 원서의 소스 저장소 [8장 부분](https://github.com/grpc-up-and-running/samples/tree/master/ch08)을 참고합니다.
