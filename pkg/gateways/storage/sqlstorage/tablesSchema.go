package sqlstorage

import (
	"expenses-app/pkg/domain/expense"
	"time"
)

var thisMonth = time.Now().Format("2006-01")
var pastMonth = time.Now().AddDate(0, -1, 0).Format("2006-01")

var MySQLInserts string = ` 
INSERT IGNORE INTO categories VALUES('0c202ba7-39a8-4f67-bbe1-9dcb30d2a346','Unknown');
INSERT IGNORE INTO categories VALUES('5e27a713-1a0b-45d5-8c59-2c6aff0cd0ae','Services');
INSERT IGNORE INTO categories VALUES('40a82e30-a9ee-4cc4-8018-de1e98c4d3be','Entertainment');
INSERT IGNORE INTO categories VALUES('15e090a7-0318-40e8-b657-8289708a1ba4','Fishes');
INSERT IGNORE INTO categories VALUES('72ee0c4a-3809-41ef-98a5-5630a9a53242','Construction');
INSERT IGNORE INTO categories VALUES('3b325284-e437-49c1-b311-9dae17c47eed','Home');
INSERT IGNORE INTO categories VALUES('12e78baa-8785-419b-aca0-b625d5fb0b49','Food');
INSERT IGNORE INTO categories VALUES('40373de2-de1f-449d-92ef-29da67322efd','Clubs');
INSERT IGNORE INTO categories VALUES('289a075e-69c5-446e-8afb-e47f6f59d74e','Investments');
INSERT IGNORE INTO categories VALUES('2b0e8ddf-29b6-4b0f-b15e-7a7c68520583','Donations');

INSERT IGNORE INTO rules VALUES('4df66260-5457-4eae-805f-b42bfaa9cdff','BadaBing','40373de2-de1f-449d-92ef-29da67322efd');
INSERT IGNORE INTO rules VALUES('5be2b350-d963-49d7-bcd0-6c2ebe73b9f5','Gardener','3b325284-e437-49c1-b311-9dae17c47eed');INSERT INTO rules VALUES('47805da9-10ce-4f95-bf0c-a8d848e9130e','Nuovo Vesuvio S.R.L','12e78baa-8785-419b-aca0-b625d5fb0b49');
INSERT IGNORE INTO rules VALUES('af11376c-88ba-48bc-aaf4-0aa0de3a31a6','Transfer to Phil Intintola','2b0e8ddf-29b6-4b0f-b15e-7a7c68520583');

INSERT IGNORE INTO expenses VALUES('1f37f2b2-2df0-4774-b52c-3c7cfd5ad7f3',45.2999999999999971,'Groceries','Centanni’s Market','` + pastMonth + `-01 00:00:00+00:00','3b325284-e437-49c1-b311-9dae17c47eed');
INSERT IGNORE INTO expenses VALUES('22f488f9-5d77-434f-80ce-3efa6b245500',100.5,'ELECTRICITY BILL','Jersey Power','` + pastMonth + `-03 00:00:00+00:00','40a82e30-a9ee-4cc4-8018-de1e98c4d3be');
INSERT IGNORE INTO expenses VALUES('52c30029-c0a2-4e17-b830-607c0c0393b8',60.75,'Gasoline','Exxon','` + pastMonth + `-04 00:00:00+00:00','72ee0c4a-3809-41ef-98a5-5630a9a53242');
INSERT IGNORE INTO expenses VALUES('183e01a2-7892-4932-a190-3c8c0946bca6',1.0,'Unknown with/o ppl','Example2','` + pastMonth + `-05 00:00:00+00:00','0c202ba7-39a8-4f67-bbe1-9dcb30d2a346');
INSERT IGNORE INTO expenses VALUES('8e752e9a-3027-4b1a-a225-b270076a7ea9',31.0,'Direct Debit','Excavator Renting','` + thisMonth + `-28 00:00:00+00:00','5e27a713-1a0b-45d5-8c59-2c6aff0cd0ae');
INSERT IGNORE INTO expenses VALUES('2784bcfd-da55-4783-9a04-b198ebb4d9a7',1.0,'Unknown with ppl','Example1','` + thisMonth + `-01 00:00:00+00:00','0c202ba7-39a8-4f67-bbe1-9dcb30d2a346');
INSERT IGNORE INTO expenses VALUES('a4c805cb-77a0-4907-8e18-2abe22bde140',3000.0,'Monthly donation','Transfer to Phil Intintola','` + thisMonth + `-10 00:00:00+00:00','2b0e8ddf-29b6-4b0f-b15e-7a7c68520583');
INSERT IGNORE INTO expenses VALUES('996b8883-4270-4deb-8dce-e14d8d9e5e58',1000.0,'Lend to a friend','Cash','` + thisMonth + `-09 00:00:00+00:00','40373de2-de1f-449d-92ef-29da67322efd');
INSERT IGNORE INTO expenses VALUES('82f37c1b-0e22-4566-b987-d195fb9ab7f4',436.0,'Dinner','Vesuvio','` + thisMonth + `-23 00:00:00+00:00','12e78baa-8785-419b-aca0-b625d5fb0b49');

INSERT IGNORE INTO expense_users VALUES('8e752e9a-3027-4b1a-a225-b270076a7ea9','9e10c58f-adb3-419b-9e20-f4fbb075661e');
INSERT IGNORE INTO expense_users VALUES('52c30029-c0a2-4e17-b830-607c0c0393b8','9e10c58f-adb3-419b-9e20-f4fbb075661e');
INSERT IGNORE INTO expense_users VALUES('22f488f9-5d77-434f-80ce-3efa6b245500','9e10c58f-adb3-419b-9e20-f4fbb075661e');
INSERT IGNORE INTO expense_users VALUES('1f37f2b2-2df0-4774-b52c-3c7cfd5ad7f3','9e10c58f-adb3-419b-9e20-f4fbb075661e');
INSERT IGNORE INTO expense_users VALUES('a4c805cb-77a0-4907-8e18-2abe22bde140','a6f4fe7e-52b6-48fc-ba0b-bef77940168f');
INSERT IGNORE INTO expense_users VALUES('2784bcfd-da55-4783-9a04-b198ebb4d9a7','9e10c58f-adb3-419b-9e20-f4fbb075661e');
INSERT IGNORE INTO expense_users VALUES('996b8883-4270-4deb-8dce-e14d8d9e5e58','2c0e2b2c-1794-4b8c-a43d-1cab7c3a8ea6');
INSERT IGNORE INTO expense_users VALUES('82f37c1b-0e22-4566-b987-d195fb9ab7f4','af657517-64d2-44e3-865b-142bb18296ab');
INSERT IGNORE INTO expense_users VALUES('82f37c1b-0e22-4566-b987-d195fb9ab7f4','efc7d7a8-aeb5-4a8a-aa38-350ca18a8873');

INSERT IGNORE INTO rule_users VALUES('4df66260-5457-4eae-805f-b42bfaa9cdff','9e10c58f-adb3-419b-9e20-f4fbb075661e');
INSERT IGNORE INTO rule_users VALUES('5be2b350-d963-49d7-bcd0-6c2ebe73b9f5','9e10c58f-adb3-419b-9e20-f4fbb075661e');
INSERT IGNORE INTO rule_users VALUES('5be2b350-d963-49d7-bcd0-6c2ebe73b9f5','a6f4fe7e-52b6-48fc-ba0b-bef77940168f');
INSERT IGNORE INTO rule_users VALUES('47805da9-10ce-4f95-bf0c-a8d848e9130e','9e10c58f-adb3-419b-9e20-f4fbb075661e');
INSERT IGNORE INTO rule_users VALUES('af11376c-88ba-48bc-aaf4-0aa0de3a31a6','a6f4fe7e-52b6-48fc-ba0b-bef77940168f');
`
var MySQLTables string = ` 
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

INSERT IGNORE INTO categories (id, name) VALUES ('` + expense.UnknownCategoryID + `', 'Unknown');
`

var SQLiteInserts string = ` 
INSERT OR IGNORE INTO categories VALUES('0c202ba7-39a8-4f67-bbe1-9dcb30d2a346','Unknown');
INSERT OR IGNORE INTO categories VALUES('5e27a713-1a0b-45d5-8c59-2c6aff0cd0ae','Services');
INSERT OR IGNORE INTO categories VALUES('40a82e30-a9ee-4cc4-8018-de1e98c4d3be','Entertainment');
INSERT OR IGNORE INTO categories VALUES('15e090a7-0318-40e8-b657-8289708a1ba4','Fishes');
INSERT OR IGNORE INTO categories VALUES('72ee0c4a-3809-41ef-98a5-5630a9a53242','Construction');
INSERT OR IGNORE INTO categories VALUES('3b325284-e437-49c1-b311-9dae17c47eed','Home');
INSERT OR IGNORE INTO categories VALUES('12e78baa-8785-419b-aca0-b625d5fb0b49','Food');
INSERT OR IGNORE INTO categories VALUES('40373de2-de1f-449d-92ef-29da67322efd','Clubs');
INSERT OR IGNORE INTO categories VALUES('289a075e-69c5-446e-8afb-e47f6f59d74e','Investments');
INSERT OR IGNORE INTO categories VALUES('2b0e8ddf-29b6-4b0f-b15e-7a7c68520583','Donations');

INSERT OR IGNORE INTO rules VALUES('4df66260-5457-4eae-805f-b42bfaa9cdff','BadaBing','40373de2-de1f-449d-92ef-29da67322efd');
INSERT OR IGNORE INTO rules VALUES('5be2b350-d963-49d7-bcd0-6c2ebe73b9f5','Gardener','3b325284-e437-49c1-b311-9dae17c47eed');
INSERT OR IGNORE INTO rules VALUES('47805da9-10ce-4f95-bf0c-a8d848e9130e','Nuovo Vesuvio S.R.L','12e78baa-8785-419b-aca0-b625d5fb0b49');
INSERT OR IGNORE INTO rules VALUES('af11376c-88ba-48bc-aaf4-0aa0de3a31a6','Transfer to Phil Intintola','2b0e8ddf-29b6-4b0f-b15e-7a7c68520583');

INSERT OR IGNORE INTO expenses VALUES('1f37f2b2-2df0-4774-b52c-3c7cfd5ad7f3',45.2999999999999971,'Groceries','Centanni’s Market','` + pastMonth + `-01 00:00:00+00:00','3b325284-e437-49c1-b311-9dae17c47eed');
INSERT OR IGNORE INTO expenses VALUES('22f488f9-5d77-434f-80ce-3efa6b245500',100.5,'ELECTRICITY BILL','Jersey Power','` + pastMonth + `-03 00:00:00+00:00','40a82e30-a9ee-4cc4-8018-de1e98c4d3be');
INSERT OR IGNORE INTO expenses VALUES('52c30029-c0a2-4e17-b830-607c0c0393b8',60.75,'Gasoline','Exxon','` + pastMonth + `-04 00:00:00+00:00','72ee0c4a-3809-41ef-98a5-5630a9a53242');
INSERT OR IGNORE INTO expenses VALUES('183e01a2-7892-4932-a190-3c8c0946bca6',1.0,'Unknown with/o ppl','Example2','` + pastMonth + `-05 00:00:00+00:00','0c202ba7-39a8-4f67-bbe1-9dcb30d2a346');
INSERT OR IGNORE INTO expenses VALUES('8e752e9a-3027-4b1a-a225-b270076a7ea9',31.0,'Direct Debit','Excavator Renting','` + thisMonth + `-28 00:00:00+00:00','5e27a713-1a0b-45d5-8c59-2c6aff0cd0ae');
INSERT OR IGNORE INTO expenses VALUES('2784bcfd-da55-4783-9a04-b198ebb4d9a7',1.0,'Unknown with ppl','Example1','` + thisMonth + `-01 00:00:00+00:00','0c202ba7-39a8-4f67-bbe1-9dcb30d2a346');
INSERT OR IGNORE INTO expenses VALUES('a4c805cb-77a0-4907-8e18-2abe22bde140',3000.0,'Monthly donation','Transfer to Phil Intintola','` + thisMonth + `-10 00:00:00+00:00','2b0e8ddf-29b6-4b0f-b15e-7a7c68520583');
INSERT OR IGNORE INTO expenses VALUES('996b8883-4270-4deb-8dce-e14d8d9e5e58',1000.0,'Lend to a friend','Cash','` + thisMonth + `-09 00:00:00+00:00','40373de2-de1f-449d-92ef-29da67322efd');
INSERT OR IGNORE INTO expenses VALUES('82f37c1b-0e22-4566-b987-d195fb9ab7f4',436.0,'Dinner','Vesuvio','` + thisMonth + `-23 00:00:00+00:00','12e78baa-8785-419b-aca0-b625d5fb0b49');

INSERT OR IGNORE INTO expense_users VALUES('8e752e9a-3027-4b1a-a225-b270076a7ea9','9e10c58f-adb3-419b-9e20-f4fbb075661e');
INSERT OR IGNORE INTO expense_users VALUES('52c30029-c0a2-4e17-b830-607c0c0393b8','9e10c58f-adb3-419b-9e20-f4fbb075661e');
INSERT OR IGNORE INTO expense_users VALUES('22f488f9-5d77-434f-80ce-3efa6b245500','9e10c58f-adb3-419b-9e20-f4fbb075661e');
INSERT OR IGNORE INTO expense_users VALUES('1f37f2b2-2df0-4774-b52c-3c7cfd5ad7f3','9e10c58f-adb3-419b-9e20-f4fbb075661e');
INSERT OR IGNORE INTO expense_users VALUES('a4c805cb-77a0-4907-8e18-2abe22bde140','a6f4fe7e-52b6-48fc-ba0b-bef77940168f');
INSERT OR IGNORE INTO expense_users VALUES('2784bcfd-da55-4783-9a04-b198ebb4d9a7','9e10c58f-adb3-419b-9e20-f4fbb075661e');
INSERT OR IGNORE INTO expense_users VALUES('996b8883-4270-4deb-8dce-e14d8d9e5e58','2c0e2b2c-1794-4b8c-a43d-1cab7c3a8ea6');
INSERT OR IGNORE INTO expense_users VALUES('82f37c1b-0e22-4566-b987-d195fb9ab7f4','af657517-64d2-44e3-865b-142bb18296ab');
INSERT OR IGNORE INTO expense_users VALUES('82f37c1b-0e22-4566-b987-d195fb9ab7f4','efc7d7a8-aeb5-4a8a-aa38-350ca18a8873');

INSERT OR IGNORE INTO rule_users VALUES('4df66260-5457-4eae-805f-b42bfaa9cdff','9e10c58f-adb3-419b-9e20-f4fbb075661e');
INSERT OR IGNORE INTO rule_users VALUES('5be2b350-d963-49d7-bcd0-6c2ebe73b9f5','9e10c58f-adb3-419b-9e20-f4fbb075661e');
INSERT OR IGNORE INTO rule_users VALUES('5be2b350-d963-49d7-bcd0-6c2ebe73b9f5','a6f4fe7e-52b6-48fc-ba0b-bef77940168f');
INSERT OR IGNORE INTO rule_users VALUES('47805da9-10ce-4f95-bf0c-a8d848e9130e','9e10c58f-adb3-419b-9e20-f4fbb075661e');
INSERT OR IGNORE INTO rule_users VALUES('af11376c-88ba-48bc-aaf4-0aa0de3a31a6','a6f4fe7e-52b6-48fc-ba0b-bef77940168f');
`
var SQLiteTables string = ` 
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

INSERT OR IGNORE INTO categories (id, name) VALUES ('` + expense.UnknownCategoryID + `', 'Unknown');
`
