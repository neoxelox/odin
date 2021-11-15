CREATE TABLE "community" (
    "id"            VARCHAR(20) PRIMARY KEY,
    "address"       VARCHAR(100) NOT NULL,
    "name"          VARCHAR(100) NOT NULL,
    "categories"    VARCHAR(100) ARRAY NOT NULL,
    "pinned_ids"    VARCHAR(20) ARRAY NOT NULL,
    "created_at"    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at"    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "deleted_at"    TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX CONCURRENTLY "community_address_idx" ON "community" USING gin ("address" gin_trgm_ops);
CREATE INDEX CONCURRENTLY "community_name_idx" ON "community" USING gin ("name" gin_trgm_ops);
