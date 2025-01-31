package org.example.service.db;

import org.example.service.entity.User;
import org.example.service.services.CreateStatus;
import org.example.service.services.Result;
import org.example.service.services.GetUserStatus;

import java.util.Date;

public interface UserDAOInterface {
    Result<User, CreateStatus> addUser(User user);
    void updateUserName(String owner_id, String viewer_id, String name);
    void updateUserUsername(String owner_id, String viewer_id, String username);
    void updateUserPhoto(String owner_id, String viewer_id, String photo);
    void updateUserDate(String owner_id, String viewer_id, Date date);
    void updateBirthUserVisibility(String owner_id, String viewer_id, User.VisibilityStatus visibility);
    void updatePhoneUserVisibility(String owner_id, String viewer_id, User.VisibilityStatus visibility);

    void DeleteUser(org.example.service.User.UUID owner_id, org.example.service.User.UUID viewer_id);

    Result<User, GetUserStatus> getUserByUsername(String username);

    Result<User, GetUserStatus> getUserByPhone(String phone);
}
