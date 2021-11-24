CREATE TABLE "session" (
    "id"            VARCHAR(20) PRIMARY KEY,
    "user_id"       VARCHAR(20) NOT NULL,
    "metadata"      JSONB NOT NULL,
    "created_at"    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_seen_at"  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "expired_at"    TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX CONCURRENTLY "session_user_id_idx" ON "session" ("user_id");
