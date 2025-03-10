CREATE SCHEMA IF NOT EXISTS "messaging";

CREATE TYPE messaging.chat_type AS ENUM ('personal', 'group', 'secret_personal', 'secret_group');

CREATE TABLE IF NOT EXISTS messaging.chat (
    chat_id UUID PRIMARY KEY,
    chat_type messaging.chat_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS messaging.personal_chat (
    chat_id UUID NOT NULL PRIMARY KEY REFERENCES messaging.chat (chat_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS messaging.group_chat (
    chat_id UUID NOT NULL PRIMARY KEY REFERENCES messaging.chat (chat_id) ON DELETE CASCADE,
    admin_id UUID NOT NULL,
    group_name VARCHAR(255) NOT NULL,
    group_photo VARCHAR(255),
    group_description TEXT
);

CREATE TABLE IF NOT EXISTS messaging.secret_personal_chat (
    chat_id UUID NOT NULL PRIMARY KEY REFERENCES messaging.chat(chat_id) ON DELETE CASCADE,
    expiration_seconds BIGINT
);

CREATE TABLE IF NOT EXISTS messaging.secret_group_chat (
    chat_id UUID NOT NULL PRIMARY KEY REFERENCES messaging.chat (chat_id) ON DELETE CASCADE,
    admin_id UUID NOT NULL,
    group_name VARCHAR(255) NOT NULL,
    group_photo VARCHAR(255),
    group_description TEXT,
    expiration_seconds BIGINT
);

CREATE TABLE IF NOT EXISTS messaging.membership (
    user_id UUID NOT NULL,
    chat_id UUID NOT NULL REFERENCES messaging.chat (chat_id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, chat_id)
);

CREATE INDEX IF NOT EXISTS membership_user_id_idx 
    ON messaging.membership ("chat_id");

CREATE TABLE IF NOT EXISTS messaging.blocking (
    user_id UUID NOT NULL,
    chat_id UUID NOT NULL REFERENCES messaging.chat (chat_id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, chat_id)
);

CREATE INDEX IF NOT EXISTS blocking_user_id_idx 
    ON messaging.blocking ("chat_id");

--------------------------------------------------------------------------------

CREATE FUNCTION messaging.check_personal_chat_has_exactly_two_members() RETURNS TRIGGER AS $$
DECLARE
    num_members INT;
BEGIN

    num_members := (SELECT COUNT(*) FROM messaging.membership WHERE chat_id = NEW.chat_id);
    IF (num_members != 2) THEN
        RAISE EXCEPTION 
            'Personal chat must have exactly two members but it has % members', num_members;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE CONSTRAINT TRIGGER ensure_personal_chat_has_exact_two_members
    AFTER INSERT OR UPDATE 
    ON messaging.membership
    FOR EACH ROW
    DEFERRABLE INITIALLY DEFERRED
    WHEN (EXISTS (SELECT * FROM messaging.personal_chat WHERE chat_id = NEW.chat_id))
EXECUTE PROCEDURE messaging.check_personal_chat_has_exactly_two_members();

--------------------------------------------------------------------------------

CREATE FUNCTION messaging.check_blocking_user_is_member() RETURNS TRIGGER AS $$
BEGIN
    IF NOT EXISTS (SELECT 1
                   FROM messaging.membership 
                   WHERE user_id = NEW.user_id AND chat_id = NEW.chat_id) 
    THEN
        RAISE EXCEPTION 'Blocking chat user must be a member of the chat';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER ensure_blocking_user_is_not_member
    AFTER INSERT OR UPDATE ON messaging.blocking
    FOR EACH ROW
EXECUTE PROCEDURE messaging.check_blocking_user_is_member();