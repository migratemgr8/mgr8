ALTER TABLE users ADD COLUMN ddi varchar(20);
ALTER TABLE users ALTER COLUMN num TYPE numeric(20, 10);
ALTER TABLE users ALTER COLUMN nombre TYPE varchar(20);
ALTER TABLE users ALTER COLUMN nombre SET NOT NULL;
ALTER TABLE users ALTER COLUMN phone TYPE varchar(20);
CREATE TABLE products (
nombre varchar(256),
qnt int4
);
