CREATE SCHEMA "users";

CREATE TYPE users.user_field AS ENUM ('date_of_birth', 'phone');

CREATE TYPE users.field_visibility AS ENUM ('everyone', 'specfied', 'only_me');

CREATE TABLE IF NOT EXISTS users.user (
    id          UUID    PRIMARY KEY,
	name        TEXT NOT NULL,
	username    TEXT NOT NULL,
	phone       TEXT,
	date_of_birth TIMESTAMP WITH TIME ZONE,
	photo_url    TEXT,
	created_at   BIGINT NOT NULL,

	date_of_birth_visibility users.field_visibility DEFAULT 'everyone',
	phone_visibility       users.field_visibility  DEFAULT 'everyone'
);

CREATE TABLE IF NOT EXISTS users.field_restrictions (
    owner_user_id UUID REFERENCES users.user (id) ON DELETE CASCADE,
    field_name users.user_field NOT NULL,
    permitted_user_id UUID REFERENCES users.user (id) ON DELETE CASCADE,
    PRIMARY KEY (
        owner_user_id, 
        field_name, 
        permitted_user_id
    )
);