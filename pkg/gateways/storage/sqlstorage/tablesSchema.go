package sqlstorage

import "expenses-app/pkg/domain/expense"

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
  CONSTRAINT fk_expenses_category FOREIGN KEY (category_id) REFERENCES categories (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS expense_users (
  expense_id VARCHAR(255) NOT NULL,
  user_id VARCHAR(255) NOT NULL,
  PRIMARY KEY (expense_id, user_id),
  FOREIGN KEY (expense_id) REFERENCES expenses (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT IGNORE INTO categories (id, name) VALUES ('` + expense.UnkownCategoryID + `', 'Unknown');
INSERT IGNORE INTO categories (id, name) VALUES ('5e27a713-1a0b-45d5-8c59-2c6aff0cd0ae', 'Services');
INSERT IGNORE INTO categories (id, name) VALUES ('40a82e30-a9ee-4cc4-8018-de1e98c4d3be', 'Entertainment');
INSERT IGNORE INTO categories (id, name) VALUES ('15e090a7-0318-40e8-b657-8289708a1ba4', 'Salivery');
INSERT IGNORE INTO categories (id, name) VALUES ('72ee0c4a-3809-41ef-98a5-5630a9a53242', 'Transportation');
INSERT IGNORE INTO categories (id, name) VALUES ('3b325284-e437-49c1-b311-9dae17c47eed', 'Home');

INSERT IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('1f37f2b2-2df0-4774-b52c-3c7cfd5ad7f3', 45.30, 'Groceries', 'Walmart', '2024-06-01 10:30:00', '3b325284-e437-49c1-b311-9dae17c47eed');
INSERT IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('5bc8b29b-e379-4e1b-96a6-279e32c603ea', 15.00, 'Movie Ticket', 'AMC', '2024-06-02 18:00:00', '15e090a7-0318-40e8-b657-8289708a1ba4');
INSERT IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('22f488f9-5d77-434f-80ce-3efa6b245500', 100.50, 'Electricty Bill', 'Electric company', '2024-06-03 09:00:00', '5e27a713-1a0b-45d5-8c59-2c6aff0cd0ae');
INSERT IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('52c30029-c0a2-4e17-b830-607c0c0393b8', 60.75, 'Gasoline', 'Shell', '2024-06-04 08:00:00', '72ee0c4a-3809-41ef-98a5-5630a9a53242');
INSERT IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('fed8622c-9f5b-4207-8b83-2db5bb42103f', 200.00, 'Doctor Visit', 'Health Clinic', '2024-06-05 11:00:00', '3b325284-e437-49c1-b311-9dae17c47eed');

INSERT IGNORE INTO expense_users (expense_id, user_id) VALUES ('1f37f2b2-2df0-4774-b52c-3c7cfd5ad7f3', 'a46927b9-7da6-440a-bae1-694a13c4cdc2');
INSERT IGNORE INTO expense_users (expense_id, user_id) VALUES ('5bc8b29b-e379-4e1b-96a6-279e32c603ea', 'b2f77ba9-c7e2-40dc-abf3-e6748611c3e0');
INSERT IGNORE INTO expense_users (expense_id, user_id) VALUES ('22f488f9-5d77-434f-80ce-3efa6b245500', 'b2f77ba9-c7e2-40dc-abf3-e6748611c3e0');
INSERT IGNORE INTO expense_users (expense_id, user_id) VALUES ('22f488f9-5d77-434f-80ce-3efa6b245500', 'a46927b9-7da6-440a-bae1-694a13c4cdc2');
INSERT IGNORE INTO expense_users (expense_id, user_id) VALUES ('52c30029-c0a2-4e17-b830-607c0c0393b8', 'b2f77ba9-c7e2-40dc-abf3-e6748611c3e0');
INSERT IGNORE INTO expense_users (expense_id, user_id) VALUES ('52c30029-c0a2-4e17-b830-607c0c0393b8', 'a46927b9-7da6-440a-bae1-694a13c4cdc2');
INSERT IGNORE INTO expense_users (expense_id, user_id) VALUES ('fed8622c-9f5b-4207-8b83-2db5bb42103f', 'a46927b9-7da6-440a-bae1-694a13c4cdc2');
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

CREATE TABLE IF NOT EXISTS expense_users (
  expense_id TEXT NOT NULL,
  user_id TEXT NOT NULL,
  PRIMARY KEY (expense_id, user_id),
  FOREIGN KEY (expense_id) REFERENCES expenses (id) ON DELETE CASCADE
);

INSERT OR IGNORE INTO categories (id, name) VALUES ('` + expense.UnkownCategoryID + `', 'Unknown');
INSERT OR IGNORE INTO categories (id, name) VALUES ('5e27a713-1a0b-45d5-8c59-2c6aff0cd0ae', 'Services');
INSERT OR IGNORE INTO categories (id, name) VALUES ('40a82e30-a9ee-4cc4-8018-de1e98c4d3be', 'Entertainment');
INSERT OR IGNORE INTO categories (id, name) VALUES ('15e090a7-0318-40e8-b657-8289708a1ba4', 'Salivery');
INSERT OR IGNORE INTO categories (id, name) VALUES ('72ee0c4a-3809-41ef-98a5-5630a9a53242', 'Transportation');
INSERT OR IGNORE INTO categories (id, name) VALUES ('3b325284-e437-49c1-b311-9dae17c47eed', 'Home');

INSERT OR IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('1f37f2b2-2df0-4774-b52c-3c7cfd5ad7f3', 45.30, 'Groceries', 'Walmart', '2024-06-01 10:30:00', '3b325284-e437-49c1-b311-9dae17c47eed');
INSERT OR IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('5bc8b29b-e379-4e1b-96a6-279e32c603ea', 15.00, 'Movie Ticket', 'AMC', '2024-06-02 18:00:00', '15e090a7-0318-40e8-b657-8289708a1ba4');
INSERT OR IGNORE INTO EXPENSES (id, amount, product, shop, expend_date, category_id) VALUES ('22f488f9-5d77-434f-80ce-3efa6b245500', 100.50, 'ELECTRICITY BILL', 'ELECTRIC COMPANY', '2024-06-03 09:00:00', '40a82e30-a9ee-4cc4-8018-de1e98c4d3be');
INSERT OR IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('52c30029-c0a2-4e17-b830-607c0c0393b8', 60.75, 'Gasoline', 'Shell', '2024-06-04 08:00:00', '72ee0c4a-3809-41ef-98a5-5630a9a53242');
INSERT OR IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('fed8622c-9f5b-4207-8b83-2db5bb42103f', 200.00, 'Doctor Visit', 'Health Clinic', '2024-06-05 11:00:00', '5e27a713-1a0b-45d5-8c59-2c6aff0cd0ae');
INSERT OR IGNORE INTO expenses (id, amount, product, shop, expend_date, category_id) VALUES ('183e01a2-7892-4932-a190-3c8c0946bca6', 200.00, 'Categorize me', 'Who knows?', '2024-06-05 11:00:00', '` + expense.UnkownCategoryID + `');

INSERT OR IGNORE INTO expense_users (expense_id, user_id) VALUES ('1f37f2b2-2df0-4774-b52c-3c7cfd5ad7f3', 'a46927b9-7da6-440a-bae1-694a13c4cdc2');
INSERT OR IGNORE INTO expense_users (expense_id, user_id) VALUES ('5bc8b29b-e379-4e1b-96a6-279e32c603ea', 'b2f77ba9-c7e2-40dc-abf3-e6748611c3e0');
INSERT OR IGNORE INTO expense_users (expense_id, user_id) VALUES ('22f488f9-5d77-434f-80ce-3efa6b245500', 'a46927b9-7da6-440a-bae1-694a13c4cdc2');
INSERT OR IGNORE INTO expense_users (expense_id, user_id) VALUES ('52c30029-c0a2-4e17-b830-607c0c0393b8', 'b2f77ba9-c7e2-40dc-abf3-e6748611c3e0');
INSERT OR IGNORE INTO expense_users (expense_id, user_id) VALUES ('fed8622c-9f5b-4207-8b83-2db5bb42103f', 'a46927b9-7da6-440a-bae1-694a13c4cdc2');
INSERT OR IGNORE INTO expense_users (expense_id, user_id) VALUES ('fed8622c-9f5b-4207-8b83-2db5bb42103f', 'b2f77ba9-c7e2-40dc-abf3-e6748611c3e0');

`
