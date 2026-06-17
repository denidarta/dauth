-- "up" means: apply these changes to create the schema.
-- golang-migrate runs this when you do: make migrate-up

-- Every user in the system. One identity spans all products.
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- gen_random_uuid() = Postgres auto-generates a unique ID so we never pick one manually
    email         TEXT NOT NULL UNIQUE,                       -- UNIQUE = no two rows can have the same email
    password_hash TEXT NOT NULL,                              -- we never store the real password, only a scrambled version (bcrypt)
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),         -- TIMESTAMPTZ = timestamp with timezone, DEFAULT now() = filled automatically
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Products this auth service protects (e.g. "app-a", "dashboard").
-- Rows are seeded manually; no API to create products in v1.
CREATE TABLE products (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Which users belong to which products, and what role they have.
-- This is a "junction table" — it links two other tables together.
-- Think of it as a membership card: user X is a "member" of product Y.
CREATE TABLE user_products (
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,    -- REFERENCES = this must match a real row in users; ON DELETE CASCADE = if the user is deleted, this row is deleted too
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    role       TEXT NOT NULL DEFAULT 'member',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, product_id)                                   -- PRIMARY KEY on two columns = one user can only have one role per product
);

-- Opaque refresh tokens. We only store the SHA-256 hash, never the raw token.
-- Storing the hash means even if the DB leaks, tokens can't be reused.
CREATE TABLE refresh_tokens (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Invitations sent by admins. Also stores only the hash of the invite token.
-- accepted_at is NULL until the invite is redeemed — NULL means "not yet done".
CREATE TABLE invitations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       TEXT NOT NULL,
    product_id  UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    role        TEXT NOT NULL DEFAULT 'member',
    token_hash  TEXT NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    accepted_at TIMESTAMPTZ,   -- no DEFAULT, no NOT NULL — starts as NULL on purpose
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
