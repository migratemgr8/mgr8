CREATE TABLE users (
    social_number VARCHAR(9) PRIMARY KEY,
    name VARCHAR(15),
    phone VARCHAR(11),
    ddi VARCHAR(3)
);

CREATE TABLE products (
   name VARCHAR(256),
   qnt INTEGER
);
