CREATE TABLE "otp" (
    "id"            VARCHAR(20) PRIMARY KEY,
    "asset"         VARCHAR(100) NOT NULL,
    "type"          VARCHAR(100) NOT NULL,
    "code"          VARCHAR(6) NOT NULL,
    "attempts"      INTEGER NOT NULL,
    "expires_at"    TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX CONCURRENTLY "otp_asset_cnt" ON "otp" ("asset");

CREATE INDEX CONCURRENTLY "otp_asset_idx" ON "otp" ("asset");
