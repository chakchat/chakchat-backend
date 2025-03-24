CREATE FUNCTION messaging.check_cannot_delete_subchat() RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (SELECT * FROM messaging.chat WHERE chat_id = old.chat_id) THEN 
        RAISE EXCEPTION 'cannot delete a row from % table because you should delete from messaging.chat table instead',
            TG_TABLE_NAME;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER ensure_cannot_delete_from_personal_chat_t
    BEFORE DELETE
    ON messaging.personal_chat
    FOR EACH ROW
EXECUTE PROCEDURE messaging.check_cannot_delete_subchat();

CREATE TRIGGER ensure_cannot_delete_from_group_chat_t
    BEFORE DELETE
    ON messaging.group_chat
    FOR EACH ROW
EXECUTE PROCEDURE messaging.check_cannot_delete_subchat();

CREATE TRIGGER ensure_cannot_delete_from_secret_personal_chat_t
    BEFORE DELETE
    ON messaging.secret_personal_chat
    FOR EACH ROW
EXECUTE PROCEDURE messaging.check_cannot_delete_subchat();

CREATE TRIGGER ensure_cannot_delete_from_secret_group_chat_t
    BEFORE DELETE
    ON messaging.secret_group_chat
    FOR EACH ROW
EXECUTE PROCEDURE messaging.check_cannot_delete_subchat();