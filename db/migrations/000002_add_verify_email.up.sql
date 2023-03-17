CREATE TABLE "verify_emails" (
  "id" bigserial PRIMARY KEY,
  "username" varchar,
  "email" varchar NOT NULL,
  "secret_key" varchar NOT NULL,
  "is_used" boolean NOT NULL DEFAULT false,
  "expired_at" timestamptz NOT NULL DEFAULT (now() + interval '15 minutes'),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");