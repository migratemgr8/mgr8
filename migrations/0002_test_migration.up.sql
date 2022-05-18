ALTER TABLE users ADD COLUMN ddi VARCHAR(3);

CREATE VIEW user_phones AS
SELECT name, CONCAT(ddi, phone) AS full_phone FROM users;