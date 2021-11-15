CREATE TABLE "post" (
    "id"                VARCHAR(20) PRIMARY KEY,
    "thread_id"         VARCHAR(20) NULL,
    "creator_id"        VARCHAR(20) NOT NULL,
    "last_history_id"   VARCHAR(20) NULL,
    "type"              VARCHAR(100) NOT NULL,
    "priority"          INTEGER NULL,
    "recipient_ids"     VARCHAR(20) ARRAY NULL,
    "voter_ids"         VARCHAR(20) ARRAY NOT NULL,
    "created_at"        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at"        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "deleted_at"        TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX CONCURRENTLY "post_thread_id_idx" ON "post" ("thread_id");
CREATE INDEX CONCURRENTLY "post_creator_id_idx" ON "post" ("creator_id");
CREATE INDEX CONCURRENTLY "post_type_idx" ON "post" ("type");
CREATE INDEX CONCURRENTLY "post_priority_idx" ON "post" ("priority");
CREATE INDEX CONCURRENTLY "post_recipient_ids_idx" ON "post" USING gin ("recipient_ids");
CREATE INDEX CONCURRENTLY "post_created_at_idx" ON "post" ("created_at");
