CREATE OR REPLACE FUNCTION messaging.check_chat_type() RETURNS TRIGGER AS $$
DECLARE
    must_chat_type messaging.chat_type;
BEGIN
    IF TG_TABLE_NAME = 'personal_chat' THEN
        must_chat_type := 'personal';
    ELSIF TG_TABLE_NAME = 'group_chat' THEN
        must_chat_type := 'group';
    ELSIF TG_TABLE_NAME = 'secret_personal_chat' THEN
        must_chat_type := 'personal';
    ELSIF TG_TABLE_NAME = 'secret_group_chat' THEN
        must_chat_type := 'group';
    ELSE
        RAISE EXCEPTION 'Unknown chat relation %', TG_TABLE_NAME;
    END IF;

    IF (SELECT chat_type != must_chat_type FROM messaging.chat WHERE chat_id = NEW.chat_id) THEN
        RAISE EXCEPTION 'The created chat must be of type %', must_chat_type;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER ensure_personal_chat_type
    BEFORE INSERT
    ON messaging.personal_chat
    FOR EACH ROW
EXECUTE PROCEDURE messaging.check_chat_type();

CREATE TRIGGER ensure_group_chat_type
    BEFORE INSERT
    ON messaging.group_chat
    FOR EACH ROW
EXECUTE PROCEDURE messaging.check_chat_type();

CREATE TRIGGER ensure_secret_personal_chat_type
    BEFORE INSERT
    ON messaging.secret_personal_chat
    FOR EACH ROW
EXECUTE PROCEDURE messaging.check_chat_type();

CREATE TRIGGER ensure_secret_group_chat_type
    BEFORE INSERT
    ON messaging.secret_group_chat
    FOR EACH ROW
EXECUTE PROCEDURE messaging.check_chat_type();