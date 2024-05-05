/*
############################################################
CUSTOM FUNCTION AND TRIGGERS
############################################################
*/

--
-- deleting old vector records when the joining table gets deleted
CREATE OR REPLACE FUNCTION delete_vector_if_unreferenced()
RETURNS TRIGGER AS $$
BEGIN
    -- Check if there are no more references in document_vector
    IF (SELECT COUNT(*) FROM document_vector WHERE vector_store_id = OLD.vector_store_id) = 0 THEN
        -- Check if there are no more references in website_page_vector
        IF (SELECT COUNT(*) FROM website_page_vector WHERE vector_store_id = OLD.vector_store_id) = 0 THEN
            -- Delete from vector_store if there are no references
            DELETE FROM vector_store WHERE id = OLD.vector_store_id;
        END IF;
    END IF;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_after_delete_document_vector
AFTER DELETE ON document_vector
FOR EACH ROW
EXECUTE FUNCTION delete_vector_if_unreferenced();

CREATE TRIGGER trg_after_delete_website_page_vector
AFTER DELETE ON website_page_vector
FOR EACH ROW
EXECUTE FUNCTION delete_vector_if_unreferenced();

--
-- Set llm default field to false when another record is set to be a default for the customer
CREATE OR REPLACE FUNCTION set_is_default_false()
RETURNS TRIGGER AS $$
BEGIN
    -- Check if the new or updated row is marked as default
    IF NEW.is_default THEN
        -- Special handling for NULL customer_id (global default)
        IF NEW.customer_id IS NULL THEN
            -- Update other rows that are global defaults
            UPDATE llm
            SET is_default = false
            WHERE customer_id IS NULL AND id != NEW.id AND is_default = true;
        ELSE
            -- Update other rows for the same customer
            UPDATE llm
            SET is_default = false
            WHERE customer_id = NEW.customer_id AND id != NEW.id AND is_default = true;
        END IF;
    END IF;

    -- Proceed with the insert or update
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_is_default_before_insert_or_update
BEFORE INSERT OR UPDATE ON llm
FOR EACH ROW EXECUTE FUNCTION set_is_default_false();
