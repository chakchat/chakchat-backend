package org.example.service.db;

import jakarta.persistence.*;

import java.util.List;

@Entity
@Table(name = "PHONE_VISIBILITY")
public class UserPhoneVisibility {

    @Id
    @Column(name = "owner_id")
    @GeneratedValue
    private String ownerId;

    @ElementCollection
    private List<String> viewerId;

    public UserPhoneVisibility() {}
    public UserPhoneVisibility(String ownerId, List<String> viewerId) {
        this.ownerId = ownerId;
        this.viewerId = viewerId;
    }

    public String getOwnerId() {
        return ownerId;
    }

    public void setOwnerId(String ownerId) {
        this.ownerId = ownerId;
    }

    public List<String> getViewerId() {
        return viewerId;
    }

    public void setViewerId(List<String> viewerId) {
        this.viewerId = viewerId;
    }
}