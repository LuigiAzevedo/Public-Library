CREATE TABLE IF NOT EXISTS "loans" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "user_id" int NOT NULL,
  "book_id" int NOT NULL,
  "is_returned" boolean NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "loans" ("user_id");

CREATE INDEX ON "loans" ("user_id", "book_id");

ALTER TABLE "loans" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "loans" ADD FOREIGN KEY ("book_id") REFERENCES "books" ("id");