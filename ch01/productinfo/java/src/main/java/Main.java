//import com.google.protobuf.StringValue;
import ecommerce.ProductInfoGrpc;
import ecommerce.ProductInfoOuterClass;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;

public class Main {
    public static void main(String[] arg) {
        // 코드 1-4 부분

        // 원격 서버 주소를 사용해 채널(channel)을 생성
        ManagedChannel channel = ManagedChannelBuilder.forAddress("localhost", 8080)
                .usePlaintext(true)
                .build();
        // 채널을 사용해 블로킹(blocking) 방식 스텁 생성
        ProductInfoGrpc.ProductInfoBlockingStub stub;
        stub = ProductInfoGrpc.newBlockingStub(channel);

        // 블로킹 스텁을 통한 원격 메서드 호출
        ProductInfoOuterClass.ProductID productID = stub.addProduct(
                ProductInfoOuterClass.Product.newBuilder()
                        .setName("Apple iPhone 11")
                        .setDescription("Meet Apple iPhone 11." +
                                "All-new dual-camera system with " +
                                "Ultra Wide and Night mode.")
                        .build());

        // 책에서는 stub.addProduct() 메서드에 의해 리턴되는 productID의 타입이 StringValue로 되어 있으나,
        // 정의된 proto 파일(ProductInfo.proto)는 일반 객체인 ProductID의 타입을 갖는다.
    }
}
