CREATE TYPE messaging.update_type AS ENUM (
    'text_message',
    'text_message_edited',
    'file_message',
    'reaction',
    'update_deleted',
    'secret_update'
);

CREATE TABLE messaging.update (
    chat_id UUID NOT NULL REFERENCES messaging.chat (chat_id) ON DELETE CASCADE,
    update_id BIGINT NOT NULL,
    update_type messaging.update_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sender_id UUID NOT NULL,

    PRIMARY KEY (chat_id, update_id)
);

CREATE TABLE messaging.text_message_update (
    chat_id UUID NOT NULL,
    update_id BIGINT NOT NULL,
    text TEXT NOT NULL,
    reply_to_id BIGINT,

    PRIMARY KEY (chat_id, update_id),
    FOREIGN KEY (chat_id, update_id)
        REFERENCES messaging.update (chat_id, update_id) 
        ON DELETE CASCADE,
    FOREIGN KEY (chat_id, reply_to_id)
        REFERENCES messaging.update (chat_id, update_id) 
        ON DELETE CASCADE
);

CREATE TABLE messaging.text_message_edited_update (
    chat_id UUID NOT NULL,
    update_id BIGINT NOT NULL,
    new_text TEXT NOT NULL,
    message_id BIGINT NOT NULL,

    PRIMARY KEY (chat_id, update_id),
    FOREIGN KEY (chat_id, update_id) 
        REFERENCES messaging.update (chat_id, update_id) 
        ON DELETE CASCADE,
    FOREIGN KEY (chat_id, message_id)
        REFERENCES messaging.update (chat_id, update_id) 
        ON DELETE CASCADE
);

CREATE TABLE messaging.file_message_update (
    chat_id UUID NOT NULL,
    update_id BIGINT NOT NULL,
    file_id UUID NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_mime_type VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    file_url TEXT NOT NULL,
    file_created_at BIGINT NOT NULL,
    reply_to_id BIGINT,

    PRIMARY KEY (chat_id, update_id),
    FOREIGN KEY (chat_id, update_id) 
        REFERENCES messaging.update (chat_id, update_id) 
        ON DELETE CASCADE,
    FOREIGN KEY (chat_id, reply_to_id)
        REFERENCES messaging.update (chat_id, update_id) 
        ON DELETE CASCADE
);

CREATE TABLE messaging.reaction_update (
    chat_id UUID NOT NULL,
    update_id BIGINT NOT NULL,
    reaction VARCHAR(255) NOT NULL,
    message_id BIGINT NOT NULL,

    PRIMARY KEY (chat_id, update_id),
    FOREIGN KEY (chat_id, update_id) 
        REFERENCES messaging.update (chat_id, update_id) 
        ON DELETE CASCADE,
    FOREIGN KEY (chat_id, message_id)
        REFERENCES messaging.update (chat_id, update_id) 
        ON DELETE CASCADE
);

CREATE TYPE messaging.delete_mode AS ENUM (
    'for_all',
    'for_deletion_sender'
);

CREATE TABLE messaging.update_deleted_update (
    chat_id UUID NOT NULL,
    update_id BIGINT NOT NULL,
    deleted_update_id BIGINT NOT NULL,
    mode messaging.delete_mode NOT NULL,

    PRIMARY KEY (chat_id, update_id),
    FOREIGN KEY (chat_id, update_id) 
        REFERENCES messaging.update (chat_id, update_id) 
        ON DELETE CASCADE,
    -- Actually updates should not be deleted physically
    FOREIGN KEY (chat_id, deleted_update_id)
        REFERENCES messaging.update (chat_id, update_id) 
        ON DELETE CASCADE 
);

CREATE TABLE messaging.secret_update (
    chat_id UUID NOT NULL,
    update_id BIGINT NOT NULL,
    payload BYTEA NOT NULL,
    key_hash BYTEA NOT NULL,
    initialization_vector BYTEA NOT NULL,

    PRIMARY KEY (chat_id, update_id),
    FOREIGN KEY (chat_id, update_id) 
        REFERENCES messaging.update (chat_id, update_id) 
        ON DELETE CASCADE
);
