package org.example.user;

import io.grpc.stub.StreamObserver;
import org.example.user.UserServiceGrpc;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class Userservice extends UserServiceGrpc.UserServiceImplBase {


    private static final Logger log = LoggerFactory.getLogger(Userservice.class);

    @Override
    public void getUser(User.UserRequest request, StreamObserver<User.UserResponse> responseObserver) {
        super.getUser(request, responseObserver);
    }
}
