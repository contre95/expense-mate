CREATE TABLE `categories` (
`id` varchar(255),
`name` text,
  PRIMARY KEY (`id`)
);

CREATE TABLE `expenses` (
`id` bigint(64),
`price` float(16),
`product` text,
`currency` text,
`shop` text,
`city` text,
`people` text,
`date` datetime,
`created_at` datetime,
`updated_at` datetime,
`category_id` varchar(255),
  PRIMARY KEY (`id`),
  CONSTRAINT`fk_expenses_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`)
);
