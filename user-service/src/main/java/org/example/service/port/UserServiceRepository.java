package org.example.service.port;

import org.example.service.entity.User;
import org.example.service.services.CreateStatus;
import org.example.service.services.GetUserStatus;
import org.example.service.services.Result;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;

public interface UserServiceRepository {
    Result<User, GetUserStatus> GetUserByPhone(String phone);
    Result<User, CreateStatus> CreateUser(String phone, String name, String username);
}
