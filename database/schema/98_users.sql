/*
############################################################
USERS
############################################################
*/

CREATE ROLE schema_spy LOGIN PASSWORD 'schema_spy';
GRANT CONNECT ON DATABASE aicontent TO schema_spy;
GRANT USAGE ON SCHEMA public TO schema_spy;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO schema_spy;