CREATE USER mockdb WITH PASSWORD 'password';
ALTER USER mockdb WITH CREATEDB;
\c - mockdb
CREATE DATABASE mockdb;
GRANT ALL PRIVILEGES ON DATABASE mockdb TO mockdb;
\c mockdb
\i '/docker-entrypoint-initdb.d/init_db.sql'

