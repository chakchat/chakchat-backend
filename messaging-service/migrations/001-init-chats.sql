CREATE SCHEMA IF NOT EXISTS "messaging";

SET search_path TO "messaging";

CREATE TABLE IF NOT EXISTS "chat" (
    chat_id UUID PRIMARY KEY,
    chat_type VARCHAR(31) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
);

CREATE TABLE IF NOT EXISTS "personal_chat" (
    chat_id UUID NOT NULL PRIMARY KEY REFERENCES "chat"(chat_id),
);

CREATE TABLE IF NOT EXISTS "group_chat" (
    chat_id UUID NOT NULL PRIMARY KEY REFERENCES "chat"(chat_id),
    admin_id UUID NOT NULL,
    group_name CHAR(255) NOT NULL,
    group_description TEXT,
);

CREATE TABLE IF NOT EXISTS "secret_personal_chat" (
    chat_id UUID NOT NULL PRIMARY KEY REFERENCES "chat"(chat_id),
    expiration_seconds BIGINT,
);

CREATE TABLE IF NOT EXISTS "secret_group_chat" (
    chat_id UUID NOT NULL PRIMARY KEY REFERENCES "chat"(chat_id),
    group_name CHAR(255) NOT NULL,
    group_description TEXT,
);

CREATE TABLE IF NOT EXISTS "membership" (
    user_id UUID NOT NULL,
    chat_id UUID NOT NULL,
);

CREATE UNIQUE INDEX IF NOT EXISTS "membership_user_id_chat_id_idx" 
    ON "membership" ("user_id", "chat_id");

CREATE INDEX IF NOT EXISTS "membership_user_id_idx" 
    ON "membership" ("user_id");

CREATE INDEX IF NOT EXISTS "membership_user_id_idx" 
    ON "membership" ("chat_id");
