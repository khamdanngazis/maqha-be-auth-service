-- Update test user passwords to correct bcrypt hash for "rahasia"
-- Hash: $2b$12$vn6/FSLdnCzS8wqJWAjiWOmxhzf5U1LlLdOJdAVM3jGcjilK/8Is.

UPDATE "user" 
SET password = '$2b$12$vn6/FSLdnCzS8wqJWAjiWOmxhzf5U1LlLdOJdAVM3jGcjilK/8Is.'
WHERE username IN ('admin', 'staff', 'loginuser', 'inactiveadmin', 'expiredadmin', 'newuser');

-- Verify update
SELECT username, password FROM "user" LIMIT 10;
