CREATE TABLE users (
    social_number VARCHAR(9) PRIMARY KEY,
    nombre VARCHAR(20) NOT NULL,
    phone VARCHAR(20),
    ddi VARCHAR(20),
    num DECIMAL(20, 10)
);

CREATE TABLE products (
   nombre VARCHAR(256),
   qnt INTEGER
);
