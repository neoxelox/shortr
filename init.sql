CREATE TABLE "urls" (
    "id"            SERIAL PRIMARY KEY,
    "name"          VARCHAR(100) UNIQUE NULL,
    "url"           TEXT NOT NULL,
    "hits"          INTEGER NOT NULL DEFAULT 0,
    "last_hit_at"   TIMESTAMP WITH TIME ZONE NULL,
    "created_at"    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "modified_at"   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX "name_idx" ON "urls" ("name");
