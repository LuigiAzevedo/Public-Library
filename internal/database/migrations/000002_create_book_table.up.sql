CREATE TABLE IF NOT EXISTS "books" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "title" varchar NOT NULL,
  "author" varchar NOT NULL,
  "amount" int NOT NULL,
  "updated_at" timestamp,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "books" ("title");