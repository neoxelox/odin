CREATE TABLE "user" (
    "id"                VARCHAR(20) PRIMARY KEY,
    "phone"             VARCHAR(100) NOT NULL,
    "name"              VARCHAR(100) NOT NULL,
    "email"             VARCHAR(100) NOT NULL,
    "picture"           VARCHAR(1000) NOT NULL, -- TODO: Decrease max chars to 100 TOC
    "birthday"          DATE NOT NULL,
    "language"          VARCHAR(2) NOT NULL,
    "last_session_id"   VARCHAR(20) NULL,
    "is_banned"         BOOLEAN NOT NULL,
    "created_at"        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "deleted_at"        TIMESTAMP WITH TIME ZONE NULL
);

CREATE UNIQUE INDEX CONCURRENTLY "user_soft_delete_cnt" ON "user" ("phone")
    WHERE "deleted_at" IS NULL;

CREATE INDEX CONCURRENTLY "user_phone_idx" ON "user" ("phone");
CREATE INDEX CONCURRENTLY "user_name_idx" ON "user" USING gin ("name" gin_trgm_ops);
CREATE INDEX CONCURRENTLY "user_deleted_at_idx" ON "user" ("deleted_at");
