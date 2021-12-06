CREATE TABLE "post_history" (
    "id"            VARCHAR(20) PRIMARY KEY,
    "post_id"       VARCHAR(20) NOT NULL,
    "updator_id"    VARCHAR(20) NOT NULL,
    "message"       VARCHAR(280) NOT NULL,
    "categories"    VARCHAR(100) ARRAY NOT NULL,
    "state"         VARCHAR(100) NULL,
    "media"         VARCHAR(100) ARRAY NOT NULL,
    "widgets"       JSONB NOT NULL,
    "created_at"    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX CONCURRENTLY "post_history_post_id_idx" ON "post_history" ("post_id");
CREATE INDEX CONCURRENTLY "post_history_categories_idx" ON "post_history" USING gin ("categories");
CREATE INDEX CONCURRENTLY "post_history_state_idx" ON "post_history" ("state");
