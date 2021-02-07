package ecommerce;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import java.util.logging.Logger;

/**
* productInfo 서비스에 대한 gRPC 클라이언트 샘플
*/
public class ProductInfoClient {
    public static void main(String[] args) throws InterruptedException {
    ManagedChannel channel = ManagedChannelBuilder
        .forAddress("localhost", 50051)
        .usePlaintext()
        .build();

    ProductInfoGrpc.ProductInfoBlockingStub stub =
        ProductInfoGrpc.newBlockingStub(channel);

    ProductInfoOuterClass.ProductID productID = stub.addProduct(
        ProductInfoOuterClass.Product.newBuilder()
        .setName("Apple iPhone 11")
        .setDescription("Meet Apple iPhone 11. " +
            "All-new dual-camera system with " +
            "Ultra Wide and Night mode.")
        .setPrice(1000.0f)
        .build());
    System.out.println(productID.getValue());

    ProductInfoOuterClass.Product product = stub.getProduct(productID);
    System.out.println(product.toString());
    channel.shutdown();
    }
}