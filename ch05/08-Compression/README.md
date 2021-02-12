# Section 8 : 압축 (Chapter 5 : gRPC: 고급 기능)

----
# 압축 (Compression) (Go)

## 1. 서비스 구현 
서버 부분의 압축은 다음과 같이 `google.golang.org/grpc/encoding/gzip` 패키지만 등록하면 됩니다.

이에 대한 코드는 서버의 [main.go](order-service/server/main.go)를 참조합니다.

```go
import (
	// ...
	_ "google.golang.org/grpc/encoding/gzip" // gzip compressor 등록
)
// ...
```

## 2. 클라이언트 구현
쿨라이언트 구현도 간단하게 `grpc.UseCompressor()` 함수를 사용하면 됩니다.

이에 대한 코드는 클라이언트 [main.go](order-service/client/main.go)를 참조합니다.