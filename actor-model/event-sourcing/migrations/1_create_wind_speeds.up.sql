CREATE TABLE `wind_speeds` (
  `id` varchar(64) NOT NULL,
  `value` float NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
