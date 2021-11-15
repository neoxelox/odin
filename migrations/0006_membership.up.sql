CREATE TABLE "membership" (
    "id"            VARCHAR(20) PRIMARY KEY,
    "phone"         VARCHAR(100) NULL,
    "user_id"       VARCHAR(20) NULL,
    "community_id"  VARCHAR(20) NOT NULL,
    "state"         VARCHAR(100) NOT NULL,
    "door"          VARCHAR(100) NOT NULL,
    "role"          VARCHAR(100) NOT NULL,
    "created_at"    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at"    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "deleted_at"    TIMESTAMP WITH TIME ZONE NULL
);

CREATE UNIQUE INDEX "membership_user_id_community_id_cnt" ON "membership" ("user_id", "community_id");

CREATE INDEX CONCURRENTLY "membership_phone_idx" ON "membership" ("phone");
CREATE INDEX CONCURRENTLY "membership_user_id_idx" ON "membership" ("user_id");
CREATE INDEX CONCURRENTLY "membership_community_id_idx" ON "membership" ("community_id");
