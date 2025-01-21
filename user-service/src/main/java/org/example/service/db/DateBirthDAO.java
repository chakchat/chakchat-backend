package org.example.service.db;

import org.hibernate.HibernateException;
import org.hibernate.Session;
import org.hibernate.SessionFactory;
import org.hibernate.Transaction;
import org.hibernate.cfg.Configuration;

import java.util.Date;
import java.util.List;

public class DateBirthDAO {

    SessionFactory factory = new Configuration().configure().addAnnotatedClass(User.class).buildSessionFactory();

    public void AddUserVisible(String owner_id, String viewer_id) {
        Session session = factory.openSession();
        Transaction tx = null;
        try {
            DateBirthVisibility user = session.get(DateBirthVisibility.class, owner_id);
            List<String> viewers_ids = user.getViewerId();
            viewers_ids.add(viewer_id);
            user.setViewerId(viewers_ids);
            session.merge(user);
            tx.commit();
        } catch (HibernateException ex) {
            if (tx != null) tx.rollback();
            ex.printStackTrace();
        }
        session.close();
    }

    public void DeleteUserVisible(String owner_id, String viewer_id) {
        Session session = factory.openSession();
        Transaction tx = null;
        try {
            DateBirthVisibility user = session.get(DateBirthVisibility.class, owner_id);
            List<String> viewers_ids = user.getViewerId();
            viewers_ids.remove(viewer_id);
            user.setViewerId(viewers_ids);
            session.merge(user);
            tx.commit();
        } catch (HibernateException ex) {
            if (tx != null) tx.rollback();
            ex.printStackTrace();
        }
        session.close();
    }

    public List<String> getVisibleViewers(String owner_id) {
        Session session = factory.openSession();
        Transaction tx = null;
        List<String> viewers_ids = null;
        try {
            viewers_ids = session.get(DateBirthVisibility.class, owner_id).getViewerId();
            tx.commit();
        } catch (HibernateException ex) {
            if (tx != null) tx.rollback();
            ex.printStackTrace();
        }
        session.close();
        return viewers_ids;
    }
}
