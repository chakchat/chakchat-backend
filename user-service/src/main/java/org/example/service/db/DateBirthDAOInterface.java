package org.example.service.db;

import org.example.service.User;

import java.util.List;

public interface DateBirthDAOInterface {
    void AddVisibleUser(User.UUID owner_id, User.UUID viewer_id);
    void DeleteVisibleViewer(User.UUID owner_id, User.UUID viewer_id);
    List<User.UUID> getVisibleViewers(User.UUID owner_id);
}
