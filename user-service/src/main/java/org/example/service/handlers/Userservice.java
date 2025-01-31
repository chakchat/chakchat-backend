package org.example.service.handlers;

import io.grpc.stub.StreamObserver;
import org.example.service.User;
import org.example.service.UserServiceGrpc;
import org.example.service.port.UserServiceRepository;
import org.example.service.services.CreateStatus;
import org.example.service.services.GetUserStatus;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import net.devh.boot.grpc.server.service.GrpcService;
import org.springframework.beans.factory.annotation.Autowired;

@GrpcService
public class Userservice extends UserServiceGrpc.UserServiceImplBase {

    public static final Logger logger = LoggerFactory.getLogger(Userservice.class);

    @Autowired
    private UserServiceRepository repository;

    @Override
    public void getUser(User.UserRequest request, StreamObserver<User.UserResponse> responseObserver) {

        String phone_number = request.getPhoneNumber();

        var user = repository.GetUserByPhone(phone_number);

        User.UserResponse.Builder responseBuilder = User.UserResponse.newBuilder();
        if (user.getReturnStatus() == GetUserStatus.FAILED) {
            RuntimeException runtimeException = new RuntimeException("User can't be found" + user.getResult().getPhone());
            logger.info("Access to user denied.{}", user, runtimeException);
            responseBuilder
                    .setStatus(User.UserResponseStatus.FAILED);

        } else if (user.getReturnStatus() == GetUserStatus.NOT_FOUND) {

            RuntimeException runtimeException = new RuntimeException("User doesn't exist " + user.getResult().getPhone());
            logger.info("User is not found.{}", user, runtimeException);
            responseBuilder
                    .setStatus(User.UserResponseStatus.NOT_FOUND);
        } else {

            logger.info("User successfully found: {}", user);
            responseBuilder
                    .setStatus(User.UserResponseStatus.SUCCESS)
                    .setName(user.getResult().getName())
                    .setUserName(user.getResult().getUsername())
                    .setUserId(user.getResult().getId());
        }

        responseObserver.onNext(responseBuilder.build());
        responseObserver.onCompleted();
    }

    @Override
    public void createUser(User.CreateUserRequest request, StreamObserver<User.CreateUserResponse> responseObserver) {
        String phone_number = request.getPhoneNumber();
        String username = request.getUsername();
        String name = request.getName();

        User.CreateUserResponse.Builder responseBuilder = User.CreateUserResponse.newBuilder();

        if (!CheckName(name) || !CheckUsername(username) || !CheckPhoneNumber(phone_number)) {

            RuntimeException runtimeException = new RuntimeException("Username or Phone number is invalid");
            logger.info("Name, Username or Phone number is incorrect {}{}{}{}{}", name, " ", username, " ", phone_number, runtimeException);
            responseBuilder.setStatus(User.CreateUserStatus.VALIDATION_FAILED);

        } else {
            var user = repository.CreateUser(phone_number, username, name);

            if (user.getReturnStatus() == CreateStatus.ALREADY_EXISTS) {

                RuntimeException runtimeException = new RuntimeException("Already exists " + username);
                logger.info("User has already existed", runtimeException );
                responseBuilder.setStatus(User.CreateUserStatus.ALREADY_EXISTS);

            } else if (user.getReturnStatus() == CreateStatus.CREATE_FAILED) {

                RuntimeException runtimeException = new RuntimeException("Exception while creating a user " + username);
                logger.info("Error of creating a new user", runtimeException);
                responseBuilder.setStatus(User.CreateUserStatus.CREATE_FAILED);
            } else {
                logger.info("User successfully created: {}", username);
                responseBuilder
                        .setStatus(User.CreateUserStatus.CREATED)
                        .setName(name)
                        .setUserName(username)
                        .setUserId(user.getResult().getId());
            }
            responseObserver.onNext(responseBuilder.build());
            responseObserver.onCompleted();

        }
    }

    private Boolean CheckUsername(String username) {
        return username.matches("^[a-z][_a-z0-9]{2,19}$");
    }

    private Boolean CheckName(String name) {
        return name.matches("^.{1,50}$");
    }

    private Boolean CheckPhoneNumber(String phone_number) {
        return phone_number.matches("^\\+79\\d{9}$");
    }
}
