DROP TABLE IF EXISTS urls;
CREATE TABLE IF NOT EXISTS urls (
    id BLOB PRIMARY KEY NOT NULL,
    key BLOB NOT NULL,
    url varchar(255) NOT NULL,
    expiry_time TIMESTAMP NOT NULL
);
