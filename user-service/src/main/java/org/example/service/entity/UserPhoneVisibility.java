package org.example.service.entity;

import jakarta.persistence.*;
import org.example.service.User;

import java.util.List;

@Entity
@Table(name = "PHONE_VISIBILITY")
public class UserPhoneVisibility {

    @Id
    @Column(name = "owner_id")
    @GeneratedValue
    private User.UUID ownerId;

    @ElementCollection
    private List<User.UUID> viewerId;

    public UserPhoneVisibility() {}
    public UserPhoneVisibility(User.UUID ownerId, List<User.UUID> viewerId) {
        this.ownerId = ownerId;
        this.viewerId = viewerId;
    }

    public User.UUID getOwnerId() {
        return ownerId;
    }

    public void setOwnerId(User.UUID ownerId) {
        this.ownerId = ownerId;
    }

    public List<User.UUID> getViewerId() {
        return viewerId;
    }

    public void setViewerId(List<User.UUID> viewerId) {
        this.viewerId = viewerId;
    }
}