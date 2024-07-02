package sqlstorage

const MySQLTables string = ` 
CREATE TABLE IF NOT EXISTS categories (
  id VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE IF NOT EXISTS rules (
  id VARCHAR(255) NOT NULL,
  pattern VARCHAR(255),
  category_id VARCHAR(255) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_category_id (category_id),
  CONSTRAINT fk_rules_category FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS expenses (
  id VARCHAR(255) NOT NULL,
  amount FLOAT DEFAULT 0.0,
  product VARCHAR(255),
  shop VARCHAR(255),
  expend_date DATETIME DEFAULT NULL,
  category_id VARCHAR(255) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_category_id (category_id),
  CONSTRAINT fk_expenses_category FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Insert into categories table
INSERT IGNORE INTO categories (id, name) VALUES ('unknown', 'Unknown');
INSERT IGNORE INTO categories (id, name) VALUES ('services', 'Services');
INSERT IGNORE INTO categories (id, name) VALUES ('entertainment', 'Entertainment');
INSERT IGNORE INTO categories (id, name) VALUES ('salivery', 'Salivery');
INSERT IGNORE INTO categories (id, name) VALUES ('transportation', 'Transportation');
INSERT IGNORE INTO categories (id, name) VALUES ('home', 'Home');

-- Insert into expenses table
INSERT IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('1', 45.30, 'Groceries', 'Walmart', '2024-06-01 10:30:00', 'home');
INSERT IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('2', 15.00, 'Movie Ticket', 'AMC', '2024-06-02 18:00:00', 'salivery');
INSERT IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('3', 100.50, 'Electricty Bill', 'Electric company', '2024-06-03 09:00:00', 'entertainment');
INSERT IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('4', 60.75, 'Gasoline', 'Shell', '2024-06-04 08:00:00', 'transportation');
INSERT IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('5', 200.00, 'Doctor Visit', 'Health Clinic', '2024-06-05 11:00:00', 'services');
`

const SQLiteTables string = ` 
PRAGMA foreign_keys=ON;

CREATE TABLE IF NOT EXISTS categories (
  id TEXT NOT NULL,
  name TEXT,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS rules (
  id TEXT NOT NULL,
  pattern TEXT,
  category_id TEXT DEFAULT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS expenses (
  id TEXT NOT NULL,
  amount REAL DEFAULT NULL,
  product TEXT,
  shop TEXT,
  expend_date DATETIME DEFAULT NULL,
  category_id TEXT DEFAULT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
);


-- Insert into categories table
INSERT OR IGNORE INTO categories (id, name) VALUES ('unknown', 'Unknown');
INSERT OR IGNORE INTO categories (id, name) VALUES ('services', 'Services');
INSERT OR IGNORE INTO categories (id, name) VALUES ('entertainment', 'Entertainment');
INSERT OR IGNORE INTO categories (id, name) VALUES ('salivery', 'Salivery');
INSERT OR IGNORE INTO categories (id, name) VALUES ('transportation', 'Transportation');
INSERT OR IGNORE INTO categories (id, name) VALUES ('home', 'Home');

-- Insert into expenses table
INSERT OR IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('1', 45.30, 'Groceries', 'Walmart', '2024-06-01 10:30:00', 'home');
INSERT OR IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('2', 15.00, 'Movie Ticket', 'AMC', '2024-06-02 18:00:00', 'salivery');
INSERT OR IGNORE INTO EXPENSES (id, amount, product, shop, expend_date, category_id) VALUES ('3', 100.50, 'ELECTRICITY BILL', 'ELECTRIC COMPANY', '2024-06-03 09:00:00', 'entertainment');
INSERT OR IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('4', 60.75, 'Gasoline', 'Shell', '2024-06-04 08:00:00', 'transportation');
INSERT OR IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('5', 200.00, 'Doctor Visit', 'Health Clinic', '2024-06-05 11:00:00', 'services');
`
