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

-- Trigger for document_vector
CREATE TRIGGER trg_after_delete_document_vector
AFTER DELETE ON document_vector
FOR EACH ROW
EXECUTE FUNCTION delete_vector_if_unreferenced();

-- Trigger for website_page_vector
CREATE TRIGGER trg_after_delete_website_page_vector
AFTER DELETE ON website_page_vector
FOR EACH ROW
EXECUTE FUNCTION delete_vector_if_unreferenced();