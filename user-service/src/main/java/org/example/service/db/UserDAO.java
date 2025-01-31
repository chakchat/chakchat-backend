package org.example.service.db;

import jakarta.persistence.criteria.CriteriaBuilder;
import jakarta.persistence.criteria.CriteriaQuery;
import jakarta.persistence.criteria.Root;
import org.example.service.entity.User;
import org.example.service.services.CreateStatus;
import org.example.service.services.Result;
import org.example.service.services.GetUserStatus;
import org.hibernate.*;
import org.hibernate.cfg.*;

import java.util.Date;
import java.util.Objects;

public class UserDAO implements UserDAOInterface{
    SessionFactory factory = new Configuration().configure().addAnnotatedClass(User.class).buildSessionFactory();

    @Override
    public Result<User, CreateStatus> addUser(User user) {
        Result<User, CreateStatus> result = new Result<>();
        if (getUserByUsername(user.getUsername()).getReturnStatus() == GetUserStatus.SUCCESS) {
            result.setReturnStatus(CreateStatus.ALREADY_EXISTS);
            return result;
        }
        Session session = factory.getCurrentSession();
        Transaction tx = null;
        try {
            tx = session.beginTransaction();
            org.example.service.User.UUID user_id = (org.example.service.User.UUID) session.save(user);
            user.setCreatedAt(new Date());
            user.setId(user_id);
            tx.commit();
            result.setReturnStatus(CreateStatus.CREATED);
            result.setResult(user);
        } catch (HibernateException ex) {
            if (tx != null) tx.rollback();
            ex.printStackTrace();
            result.setReturnStatus(CreateStatus.CREATE_FAILED);
        }
        session.close();
        return result;
    }

    @Override
    public void updateUserName(String owner_id, String viewer_id, String name) {
        if (Objects.equals(owner_id, viewer_id)) {
            Session session = factory.openSession();
            Transaction tx = null;
            try {
                User user = session.get(User.class, owner_id);
                user.setName(name);
                session.merge(user);
                tx.commit();
            } catch (HibernateException ex) {
                if (tx != null) tx.rollback();
                ex.printStackTrace();
            }
            session.close();
        }
    }

    @Override
    public void updateUserUsername(String owner_id, String viewer_id, String username) {
        if (Objects.equals(owner_id, viewer_id)) {
            Session session = factory.openSession();
            Transaction tx = null;
            try {
                User user = session.get(User.class, owner_id);
                user.setUsername(username);
                session.merge(user);
                tx.commit();
            } catch (HibernateException ex) {
                if (tx != null) tx.rollback();
                ex.printStackTrace();
            }
            session.close();
        }
    }

    @Override
    public void updateUserPhoto(String owner_id, String viewer_id, String photo) {
        if (Objects.equals(owner_id, viewer_id)) {
            Session session = factory.openSession();
            Transaction tx = null;
            try {
                User user = session.get(User.class, owner_id);
                user.setPhoto(photo);
                session.merge(user);
                tx.commit();
            } catch (HibernateException ex) {
                if (tx != null) tx.rollback();
                ex.printStackTrace();
            }
            session.close();
        }
    }

    @Override
    public void updateUserDate(String owner_id, String viewer_id, Date date) {
        if (Objects.equals(owner_id, viewer_id)) {
            Session session = factory.openSession();
            Transaction tx = null;
            try {
                User user = session.get(User.class, owner_id);
                user.setDateOfBirth(date);
                session.merge(user);
                tx.commit();
            } catch (HibernateException ex) {
                if (tx != null) tx.rollback();
                ex.printStackTrace();
            }
            session.close();
        }
    }

    @Override
    public void updateBirthUserVisibility(String owner_id, String viewer_id, User.VisibilityStatus visibility) {
        if (Objects.equals(owner_id, viewer_id)) {
            Session session = factory.openSession();
            Transaction tx = null;
            try {
                User user = session.get(User.class, owner_id);
                user.setBirthVisibilityStatus(visibility);
                session.merge(user);
                tx.commit();
            } catch (HibernateException ex) {
                if (tx != null) tx.rollback();
                ex.printStackTrace();
            }
            session.close();
        }
    }

    @Override
    public void updatePhoneUserVisibility(String owner_id, String viewer_id, User.VisibilityStatus visibility) {
        if (Objects.equals(owner_id, viewer_id)) {
            Session session = factory.openSession();
            Transaction tx = null;
            try {
                User user = session.get(User.class, owner_id);
                user.setPhoneVisibilityStatus(visibility);
                session.merge(user);
                tx.commit();
            } catch (HibernateException ex) {
                if (tx != null) tx.rollback();
                ex.printStackTrace();
            }
            session.close();
        }
    }

    @Override
    public void DeleteUser(org.example.service.User.UUID owner_id, org.example.service.User.UUID viewer_id) {
        if (owner_id.equals(viewer_id)) {
            Session session = factory.openSession();
            Transaction tx = null;
            try {
                tx = session.beginTransaction();
                User user = session.get(User.class, owner_id);
                session.delete(user);
                tx.commit();
            } catch (HibernateException ex) {
                if (tx != null) tx.rollback();
                ex.printStackTrace();
            }
            session.close();
        }
    }

    @Override
    public Result<User, GetUserStatus> getUserByPhone(String phone) {
        Session session = factory.openSession();
        User user = null;
        Result<User, GetUserStatus> result = new Result<User, GetUserStatus>();
        try {

            // var owner = session.get(User.class, id_owner);

            CriteriaBuilder userCriteria = session.getCriteriaBuilder();
            CriteriaQuery<User> query = userCriteria.createQuery(User.class);
            Root<User> root = query.from(User.class);
            query.select(query.from(User.class)).where(userCriteria.equal(root.get("phone"), phone));
            user = session.createQuery(query).uniqueResult();
//            if (owner.getPhoneVisibilityStatus() == User.VisibilityStatus.None) {
//                user.setPhone(null);
//            } else if (owner.getBirthVisibilityStatus() == User.VisibilityStatus.Some) {
//                var viewers = session.get(UserPhoneVisibility.class, id_owner);
//
//                if (!viewers.getViewerId().contains(user.getId())) {
//                    result.setReturnStatus(ReturnStatus.FAILED);
//                }
//            }
//
//            if (owner.getPhoneVisibilityStatus() != User.VisibilityStatus.None) {
//                if (owner.getBirthVisibilityStatus() == User.VisibilityStatus.None) {
//                    user.setDateOfBirth(null);
//                } else if (owner.getBirthVisibilityStatus() == User.VisibilityStatus.Some) {
//                    var date_viewers = session.get(DateBirthVisibility.class, id_owner);
//                    if (!date_viewers.getViewerId().contains(user.getId())) {
//                        user.setDateOfBirth(null);
//                    }
//                }
//            }
            result.setResult(user);
            result.setReturnStatus(user == null ? GetUserStatus.NOT_FOUND : GetUserStatus.SUCCESS);
        } catch (HibernateException ex) {
            result.setReturnStatus(GetUserStatus.FAILED);
        } finally {
            session.close();
        }
        return result;
    }

    @Override
    public Result<User,GetUserStatus> getUserByUsername(String username) {
        Session session = factory.openSession();
        User user = null;
        Result<User,GetUserStatus> result = new Result<User, GetUserStatus>();
        try {
            CriteriaBuilder userCriteria = session.getCriteriaBuilder();
            CriteriaQuery<User> query = userCriteria.createQuery(User.class);
            Root<User> root = query.from(User.class);
            query.select(query.from(User.class)).where(userCriteria.equal(root.get("username"), username));
            user = session.createQuery(query).uniqueResult();
            result.setResult(user);
            result.setReturnStatus(user == null ? GetUserStatus.NOT_FOUND : GetUserStatus.SUCCESS);
        } catch (HibernateException ex) {
            result.setReturnStatus(GetUserStatus.FAILED);
        } finally {
            session.close();
        }
        return result;
    }
}