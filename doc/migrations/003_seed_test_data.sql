-- Seed test data for auth service

-- Insert test client if not exists
INSERT INTO client (company_name, email, phone_number, address, owner_name, is_active, token, created_at)
VALUES ('Test Company', 'test@company.com', '081234567890', 'Jl. Test No. 1', 'Test Owner', true, 'client_token_test_12345', NOW())
ON CONFLICT DO NOTHING;

-- Insert test users (passwords hashed with bcrypt)
-- Password hashing: All test users use password "Test@123456" hashed with bcrypt
INSERT INTO "user" (client_id, username, password, full_name, role, token, token_expired, is_active, created_at)
SELECT 
  (SELECT id FROM client WHERE token = 'client_token_test_12345' LIMIT 1),
  'admin',
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36jbMFS',
  'Admin User',
  1,
  'admin_token_12345',
  NOW() + INTERVAL '30 days',
  true,
  NOW()
WHERE NOT EXISTS (SELECT 1 FROM "user" WHERE username = 'admin')
UNION ALL
SELECT 
  (SELECT id FROM client WHERE token = 'client_token_test_12345' LIMIT 1),
  'staff',
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36jbMFS',
  'Staff User',
  2,
  'staff_token_12345',
  NOW() + INTERVAL '30 days',
  true,
  NOW()
WHERE NOT EXISTS (SELECT 1 FROM "user" WHERE username = 'staff')
UNION ALL
SELECT 
  (SELECT id FROM client WHERE token = 'client_token_test_12345' LIMIT 1),
  'loginuser',
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36jbMFS',
  'Login Test User',
  3,
  'loginuser_token_12345',
  NOW() + INTERVAL '30 days',
  true,
  NOW()
WHERE NOT EXISTS (SELECT 1 FROM "user" WHERE username = 'loginuser')
UNION ALL
SELECT 
  (SELECT id FROM client WHERE token = 'client_token_test_12345' LIMIT 1),
  'inactiveadmin',
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36jbMFS',
  'Inactive Admin User',
  1,
  'inactiveadmin_token_12345',
  NOW() + INTERVAL '30 days',
  false,
  NOW()
WHERE NOT EXISTS (SELECT 1 FROM "user" WHERE username = 'inactiveadmin')
UNION ALL
SELECT 
  (SELECT id FROM client WHERE token = 'client_token_test_12345' LIMIT 1),
  'expiredadmin',
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36jbMFS',
  'Expired Admin User',
  1,
  'expiredadmin_token_12345',
  NOW() - INTERVAL '1 days',
  true,
  NOW()
WHERE NOT EXISTS (SELECT 1 FROM "user" WHERE username = 'expiredadmin')
UNION ALL
SELECT 
  (SELECT id FROM client WHERE token = 'client_token_test_12345' LIMIT 1),
  'newuser',
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36jbMFS',
  'New User',
  3,
  'newuser_token_12345',
  NOW() + INTERVAL '30 days',
  true,
  NOW()
WHERE NOT EXISTS (SELECT 1 FROM "user" WHERE username = 'newuser');

-- Verify data inserted
SELECT 'Clients inserted:' as status, COUNT(*) as count FROM client;
SELECT 'Users inserted:' as status, COUNT(*) as count FROM "user";
