ALTER TABLE users ADD COLUMN ddi VARCHAR(3);

CREATE VIEW user_numbers AS
SELECT name, CONCAT(ddi, phone) FROM users;