package org.example.user;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import org.example.user.User;

import java.io.IOException;

public class GrpcServer {
    public static void main(String[] args) throws IOException, InterruptedException {
        Server server = ServerBuilder
                .forPort(8080)
                .addService(new Userservice()).build();
        server.start();
        server.awaitTermination();
    }
}
