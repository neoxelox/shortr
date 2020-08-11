CREATE SEQUENCE "urls_id_seq";

CREATE TABLE "urls" (
    "id"            INTEGER PRIMARY KEY DEFAULT NEXTVAL('urls_id_seq'),
    "name"          VARCHAR(100) UNIQUE NOT NULL DEFAULT CURRVAL('urls_id_seq'),
    "url"           TEXT NOT NULL,
    "hits"          INTEGER NOT NULL DEFAULT 0,
    "last_hit_at"   TIMESTAMP WITH TIME ZONE NULL,
    "created_at"    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "modified_at"   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX "name_idx" ON "urls" ("name");
