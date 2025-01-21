package org.example.service.db;

import jakarta.persistence.*;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;

import java.util.Date;

@Entity
@Table(name = "USER")
public class User {

    @Id
    @Column(name = "id")
    @GeneratedValue
    private String id;

    @Column(name = "username")
    private String username;

    @Column(name = "name")
    private String name;

    @Column(name = "phone")
    private String phone;

    @Column(name = "photo")
    private String photo;

    @Column(name = "date_birth")
    private Date dateOfBirth;

    @Column(name = "created_at")
    private Date createdAt;

    @Enumerated(EnumType.STRING)
    @Column(name = "birth_visibility")
    private VisibilityStatus birthVisibilityStatus;

    @Enumerated(EnumType.STRING)
    @Column(name = "phone_visibility")
    private VisibilityStatus phoneVisibilityStatus;

    public enum VisibilityStatus {
        All,
        None,
        Some
    }
    public User() {}
    public User(String id, String username, String name, String phone, String photo, Date dateOfBirth) {
        this.id = id;
        this.username = username;
        this.name = name;
        this.phone = phone;
        this.photo = photo;
        this.dateOfBirth = dateOfBirth;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getPhone() {
        return phone;
    }

    public void setPhone(String phone) {
        this.phone = phone;
    }

    public String getPhoto() {
        return photo;
    }

    public void setPhoto(String photo) {
        this.photo = photo;
    }

    public Date getDateOfBirth() {
        return dateOfBirth;
    }

    public void setDateOfBirth(Date dateOfBirth) {
        this.dateOfBirth = dateOfBirth;
    }

    public Date getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(Date createdAt) {
        this.createdAt = createdAt;
    }

    public VisibilityStatus getBirthVisibilityStatus() {
        return birthVisibilityStatus;
    }

    public void setBirthVisibilityStatus(VisibilityStatus birthVisibilityStatus) {
        this.birthVisibilityStatus = birthVisibilityStatus;
    }

    public VisibilityStatus getPhoneVisibilityStatus() {
        return phoneVisibilityStatus;
    }

    public void setPhoneVisibilityStatus(VisibilityStatus phoneVisibilityStatus) {
        this.phoneVisibilityStatus = phoneVisibilityStatus;
    }
}

