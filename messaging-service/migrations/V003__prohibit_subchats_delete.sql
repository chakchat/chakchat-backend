CREATE FUNCTION messaging.raise_cannot_delete_subchat() RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'cannot delete a row from % table because you should delete from messaging.chat table instead',
        TG_TABLE_NAME;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER ensure_cannot_delete_from_personal_chat_t
    BEFORE DELETE
    ON messaging.personal_chat
    FOR EACH STATEMENT
EXECUTE PROCEDURE messaging.raise_cannot_delete_subchat();

CREATE TRIGGER ensure_cannot_delete_from_group_chat_t
    BEFORE DELETE
    ON messaging.group_chat
    FOR EACH STATEMENT
EXECUTE PROCEDURE messaging.raise_cannot_delete_subchat();

CREATE TRIGGER ensure_cannot_delete_from_secret_personal_chat_t
    BEFORE DELETE
    ON messaging.secret_personal_chat
    FOR EACH STATEMENT
EXECUTE PROCEDURE messaging.raise_cannot_delete_subchat();

CREATE TRIGGER ensure_cannot_delete_from_secret_group_chat_t
    BEFORE DELETE
    ON messaging.secret_group_chat
    FOR EACH STATEMENT
EXECUTE PROCEDURE messaging.raise_cannot_delete_subchat();