CREATE TABLE "invitation" (
    "id"            VARCHAR(20) PRIMARY KEY,
    "phone"         VARCHAR(100) NOT NULL,
    "community_id"  VARCHAR(20) NOT NULL,
    "door"          VARCHAR(100) NOT NULL,
    "role"          VARCHAR(100) NOT NULL,
    "created_at"    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "reminded_at"   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "expires_at"    TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX CONCURRENTLY "invitation_phone_community_id_cnt" ON "invitation" ("phone", "community_id");

CREATE INDEX CONCURRENTLY "invitation_phone_idx" ON "invitation" ("phone");
CREATE INDEX CONCURRENTLY "invitation_community_id_idx" ON "invitation" ("community_id");
