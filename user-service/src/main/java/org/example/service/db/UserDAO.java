package org.example.service.db;

import org.hibernate.*;
import org.hibernate.cfg.*;

import java.util.Date;
import java.util.Objects;

public class UserDAO {
    SessionFactory factory = new Configuration().configure().addAnnotatedClass(User.class).buildSessionFactory();

    public void addUser(User user) {
        Session session = factory.getCurrentSession();
        Transaction tx = null;

        try {
            tx = session.beginTransaction();
            session.persist(user);
            user.setCreatedAt(new Date());
            tx.commit();
        } catch (HibernateException ex) {
            if (tx != null) tx.rollback();
            ex.printStackTrace();
        }
        session.close();
    }

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
}