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
 CONSTRAINT fk_rules_category FOREIGN KEY (category_id) REFERENCES categories (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS expenses (
  id VARCHAR(255) NOT NULL,
  amount FLOAT DEFAULT 0.0,
  product VARCHAR(255),
  shop VARCHAR(255),
  expend_date DATETIME DEFAULT NULL,
  category_id VARCHAR(255) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_catexp_id (category_id),
  CONSTRAINT fk_expenses_category FOREIGN KEY (category_id) REFERENCES categories (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS rule_users (
  rule_id VARCHAR(255) NOT NULL,
  user_id VARCHAR(255) NOT NULL,
  PRIMARY KEY (rule_id, user_id),
  FOREIGN KEY (rule_id) REFERENCES rules (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS expense_users (
  expense_id VARCHAR(255) NOT NULL,
  user_id VARCHAR(255) NOT NULL,
  PRIMARY KEY (expense_id, user_id),
  FOREIGN KEY (expense_id) REFERENCES expenses (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT IGNORE INTO categories (id, name) VALUES ('` + expense.UnkownCategoryID + `', 'Unknown');
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

CREATE TABLE IF NOT EXISTS rule_users (
  rule_id TEXT NOT NULL,
  user_id TEXT NOT NULL,
  PRIMARY KEY (rule_id, user_id),
  FOREIGN KEY (rule_id) REFERENCES rules (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS expense_users (
  expense_id TEXT NOT NULL,
  user_id TEXT NOT NULL,
  PRIMARY KEY (expense_id, user_id),
  FOREIGN KEY (expense_id) REFERENCES expenses (id) ON DELETE CASCADE
);

INSERT OR IGNORE INTO categories (id, name) VALUES ('` + expense.UnkownCategoryID + `', 'Unknown');
`
