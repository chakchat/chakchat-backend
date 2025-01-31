package org.example.service;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import org.example.service.db.UserDAO;
import org.example.service.db.UserDAOInterface;
import org.example.service.handlers.Userservice;

import java.io.IOException;

public class GrpcServer {
    public static void main(String[] args) throws IOException, InterruptedException {
        Server server = ServerBuilder
                .forPort(8080)
                .addService(new Userservice()).build();
        server.start();
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            server.shutdown();
            System.out.println("Successfully stopped the server");
        }));
        server.awaitTermination();
    }
}
