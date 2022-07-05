ALTER TABLE users ADD COLUMN ddi varchar(3);

		CREATE VIEW user_phones AS
			SELECT name,
			CONCAT(phone, ddi) AS full_phone
			FROM users;
