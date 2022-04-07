CREATE TABLE "users" (
    "id" uuid NULL,
    "name" VARCHAR NULL,
    "email" VARCHAR NULL,
    "phone" VARCHAR NULL,
    "password_hash" VARCHAR NULL,
    "active" BOOLEAN DEFAULT true,
    "created_at" TIMESTAMP NULL,
    "updated_at" TIMESTAMP NULL,
    "deleted_at" TIMESTAMP NULL,
    CONSTRAINT "users_pk" PRIMARY KEY (id)
);
CREATE UNIQUE INDEX users_id_idx ON "users" (id);
CREATE UNIQUE INDEX users_email_idx ON "users" (email);


CREATE TABLE "access_histories" (
    "user_id" uuid NULL,
    "logged_at" TIMESTAMP NULL,
    FOREIGN KEY ("user_id") REFERENCES users ("id") ON DELETE CASCADE
);
CREATE INDEX access_history_user_id_idx ON "access_histories" USING btree (user_id);
