package sqlstorage

const Tables string = ` 
CREATE TABLE IF NOT EXISTS categories (
  id TEXT NOT NULL,
  name TEXT,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS expenses (
  id TEXT NOT NULL,
  price REAL DEFAULT NULL,
  product TEXT,
  currency TEXT,
  shop TEXT,
  city TEXT,
  people TEXT,
  expend_date DATETIME DEFAULT NULL,
  category_id TEXT DEFAULT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (category_id) REFERENCES categories (id)
);

-- Insert into categories table
INSERT OR IGNORE INTO categories (id, name) VALUES ('unknown', 'Unknown');
INSERT OR IGNORE INTO categories (id, name) VALUES ('services', 'Services');
INSERT OR IGNORE INTO categories (id, name) VALUES ('entertainment', 'Entertainment');
INSERT OR IGNORE INTO categories (id, name) VALUES ('salivery', 'Salivery');
INSERT OR IGNORE INTO categories (id, name) VALUES ('transportation', 'Transportation');
INSERT OR IGNORE INTO categories (id, name) VALUES ('home', 'Home');

-- Insert into expenses table
INSERT OR IGNORE INTO expenses (id, price, product, currency, shop, city, people, expend_date, category_id) VALUES ('1', 45.30, 'Groceries', 'USD', 'Walmart', 'New York', 'John', '2024-06-01 10:30:00', 'home');
INSERT OR IGNORE INTO expenses (id, price, product, currency, shop, city, people, expend_date, category_id) VALUES ('2', 15.00, 'Movie Ticket', 'USD', 'AMC', 'Los Angeles', 'Jane', '2024-06-02 18:00:00', 'salivery');
INSERT OR IGNORE INTO EXPENSES (ID, PRICE, PRODUCT, CURRENCY, SHOP, CITY, PEOPLE, EXPEND_DATE, CATEGORY_ID) VALUES ('3', 100.50, 'ELECTRICITY BILL', 'USD', 'ELECTRIC COMPANY', 'CHICAGO', 'ALICE', '2024-06-03 09:00:00', 'entertainment');
INSERT OR IGNORE INTO expenses (id, price, product, currency, shop, city, people, expend_date, category_id) VALUES ('4', 60.75, 'Gasoline', 'USD', 'Shell', 'Houston', 'Bob', '2024-06-04 08:00:00', 'transportation');
INSERT OR IGNORE INTO expenses (id, price, product, currency, shop, city, people, expend_date, category_id) VALUES ('5', 200.00, 'Doctor Visit', 'USD', 'Health Clinic', 'Phoenix', 'Charlie', '2024-06-05 11:00:00', 'services');

`
