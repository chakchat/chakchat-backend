package org.example.service.services;

import org.example.service.entity.User;
import org.example.service.db.UserDAOInterface;
import org.example.service.port.UserServiceRepository;

public class UserService implements UserServiceRepository {


    private final UserDAOInterface userDao;


    public UserService(UserDAOInterface userDao) {

        this.userDao = userDao;
    }

    public Result<User, GetUserStatus> GetUserByPhone(String phone) {
        return userDao.getUserByPhone(phone);
    }

    public Result<User, CreateStatus> CreateUser(String phone, String name, String username) {
        User user = new User(phone, name, username);
        return userDao.addUser(user);
    }
}
