CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "hashed_password" varchar,
  "avatar_url" varchar,
  "contact" varchar,
  "security_key" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "verified_at" timestamptz,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "twitter_social" boolean NOT NULL DEFAULT FALSE,
  "google_social" boolean NOT NULL DEFAULT FALSE,
  "apple_social" boolean NOT NULL DEFAULT FALSE
);

CREATE TABLE "requests" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "type" varchar NOT NULL,
  "prefered_nurse" bigint,
  "start_date" timestamptz NOT NULL,
  "end_date" timestamptz NOT NULL,
  "location" point NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "nurses" (
  "id" bigserial PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "contact" varchar NOT NULL,
  "profile_picture" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "approvals" (
  "id" bigserial,
  "request_id" bigint NOT NULL UNIQUE,
  "user_id" bigint NOT NULL,
  "assigned_nurse" bigint NOT NULL,
  "approved_by" bigint NOT NULL,
  "status" varchar NOT NULL,
  "notes" varchar,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("id", "request_id")
);

CREATE TABLE "admins" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "is_super" boolean NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "requests" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "requests" ADD FOREIGN KEY ("prefered_nurse") REFERENCES "nurses" ("id");

ALTER TABLE "approvals" ADD FOREIGN KEY ("request_id") REFERENCES "requests" ("id");

ALTER TABLE "approvals" ADD FOREIGN KEY ("assigned_nurse") REFERENCES "nurses" ("id");

ALTER TABLE "approvals" ADD FOREIGN KEY ("approved_by") REFERENCES "admins" ("id");

