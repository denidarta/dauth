-- "down" means: undo everything the "up" file did.
-- golang-migrate runs this when you do: make migrate-down
-- Drop order matters — tables that REFERENCE others must go first,
-- or Postgres will refuse (a membership card can't exist without the user table).

DROP TABLE IF EXISTS invitations;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS user_products;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;
