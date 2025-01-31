package org.example.service.db;

import org.example.service.entity.DateBirthVisibility;
import org.example.service.entity.User;
import org.hibernate.HibernateException;
import org.hibernate.Session;
import org.hibernate.SessionFactory;
import org.hibernate.Transaction;
import org.hibernate.cfg.Configuration;

import java.util.List;

public class DateBirthDAO implements DateBirthDAOInterface{

    SessionFactory factory = new Configuration().configure().addAnnotatedClass(User.class).buildSessionFactory();

    public void AddVisibleUser(org.example.service.User.UUID owner_id, org.example.service.User.UUID viewer_id) {
        Session session = factory.openSession();
        Transaction tx = null;
        try {
            DateBirthVisibility user = session.get(DateBirthVisibility.class, owner_id);
            List<org.example.service.User.UUID> viewers_ids = user.getViewerId();
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

    public void DeleteVisibleViewer(org.example.service.User.UUID owner_id, org.example.service.User.UUID viewer_id) {
        Session session = factory.openSession();
        Transaction tx = null;
        try {
            DateBirthVisibility user = session.get(DateBirthVisibility.class, owner_id);
            List<org.example.service.User.UUID> viewers_ids = user.getViewerId();
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

    public List<org.example.service.User.UUID> getVisibleViewers(org.example.service.User.UUID owner_id) {
        Session session = factory.openSession();
        Transaction tx = null;
        List<org.example.service.User.UUID> viewers_ids = null;
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
