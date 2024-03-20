DROP TABLE IF EXISTS shelves;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS positions;

CREATE TABLE shelves (
	id SERIAL PRIMARY KEY NOT NULL,
	name TEXT
);

INSERT INTO shelves(name) VALUES ('А'),('Б'),('Ж'),('З'),('В');

CREATE TABLE products (
	id SERIAL PRIMARY KEY NOT NULL,
	name TEXT UNIQUE NOT NULL,
	base INTEGER REFERENCES shelves(id),
	additional TEXT[]
);

INSERT INTO products(id, name, base) VALUES
	(1,'Ноутбук', (SELECT id FROM shelves WHERE name='А')),
	(2,'Телевизор', (SELECT id FROM shelves WHERE name='А')),
	(4,'Системный блок', (SELECT id FROM shelves WHERE name='Ж')),
	(6,'Микрофон', (SELECT id FROM shelves WHERE name='Ж'));
INSERT INTO products(id, name, base, additional) VALUES
	(3,'Телефон', (SELECT id FROM shelves WHERE name='Б'), '{"З","В"}'),
	(5,'Часы', (SELECT id FROM shelves WHERE name='Ж'), '{"А"}');


CREATE TABLE orders (
	id SERIAL PRIMARY KEY NOT NULL
);

CREATE TABLE positions (
	id SERIAL PRIMARY KEY NOT NULL,
	order_id INTEGER REFERENCES orders(id),
	product_id INTEGER REFERENCES products(id),
	count INTEGER NOT NULL
);

BEGIN;
INSERT INTO orders(id) VALUES(10);
INSERT INTO positions(order_id, product_id, count) VALUES
	(10, (SELECT id FROM products WHERE name='Ноутбук'), 2),
	(10, (SELECT id FROM products WHERE name='Телефон'), 1),
	(10, (SELECT id FROM products WHERE name='Микрофон'), 1);
COMMIT;

BEGIN;
INSERT INTO orders(id) VALUES (11);
INSERT INTO positions(order_id, product_id, count) VALUES
	(11, (SELECT id FROM products WHERE name='Телевизор'), 3);
COMMIT;

BEGIN;
INSERT INTO orders(id) VALUES (14);
INSERT INTO positions(order_id, product_id, count) VALUES
	(14, (SELECT id FROM products WHERE name='Ноутбук'), 3),
	(14, (SELECT id FROM products WHERE name='Системный блок'), 4);
COMMIT;

BEGIN;
INSERT INTO orders(id) VALUES (15);
INSERT INTO positions(order_id, product_id, count) VALUES
	(15, (SELECT id FROM products WHERE name='Часы'), 1);
COMMIT;