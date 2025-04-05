CREATE TABLE messaging.chat_sequence (
    chat_id UUID PRIMARY KEY,
    last_update_id BIGINT NOT NULL DEFAULT 0
);

CREATE FUNCTION messaging.get_next_update_id()
RETURNS TRIGGER AS $$
DECLARE
    next_id BIGINT;
BEGIN
    PERFORM pg_advisory_xact_lock(hashtext(NEW.chat_id::text));
    
    INSERT INTO messaging.chat_sequence (chat_id, last_update_id)
    VALUES (NEW.chat_id, 0)
    ON CONFLICT (chat_id) DO NOTHING;
    
    UPDATE messaging.chat_sequence
    SET last_update_id = last_update_id + 1
    WHERE chat_id = NEW.chat_id
    RETURNING last_update_id INTO next_id;
    
    NEW.update_id := next_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_update_id_trigger
    BEFORE INSERT ON messaging.update
    FOR EACH ROW
    WHEN (NEW.update_id IS NULL OR NEW.update_id = 0)
EXECUTE FUNCTION messaging.get_next_update_id();
