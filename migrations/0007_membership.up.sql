CREATE TABLE "membership" (
    "id"            VARCHAR(20) PRIMARY KEY,
    "user_id"       VARCHAR(20) NOT NULL,
    "community_id"  VARCHAR(20) NOT NULL,
    "door"          VARCHAR(100) NOT NULL,
    "role"          VARCHAR(100) NOT NULL,
    "created_at"    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "deleted_at"    TIMESTAMP WITH TIME ZONE NULL
);

CREATE UNIQUE INDEX CONCURRENTLY "membership_user_id_community_id_cnt" ON "membership" ("user_id", "community_id");

CREATE INDEX CONCURRENTLY "membership_user_id_idx" ON "membership" ("user_id");
CREATE INDEX CONCURRENTLY "membership_community_id_idx" ON "membership" ("community_id");
CREATE INDEX CONCURRENTLY "membership_role_idx" ON "membership" ("role");
