BEGIN;

CREATE TABLE IF NOT EXISTS "urls" (
    id UUID PRIMARY KEY NOT NULL,
    key varchar(255) NOT NULL,
    url varchar(512) NOT NULL,
    expiry_time TIMESTAMP WITH TIME ZONE NOT NULL
);

END;
