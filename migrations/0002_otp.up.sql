CREATE TABLE "otp" (
    "id"            VARCHAR(20) PRIMARY KEY,
    "asset"         VARCHAR(100) NOT NULL,
    "type"          VARCHAR(100) NOT NULL,
    "code"          VARCHAR(6) NOT NULL,
    "attempts"      INTEGER NOT NULL,
    "created_at"    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at"    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "deleted_at"    TIMESTAMP WITH TIME ZONE NULL
);

CREATE UNIQUE INDEX "otp_soft_delete_cnt" ON "otp" ("asset")
    WHERE "deleted_at" IS NOT NULL;

CREATE INDEX CONCURRENTLY "otp_asset_idx" ON "otp" ("asset");
CREATE INDEX CONCURRENTLY "otp_deleted_at_idx" ON "otp" ("deleted_at");
