-- +goose Up

CREATE TABLE `users` (
  `id` bigint(20) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `lang` varchar(2) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `links` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `user_id` bigint(20) NOT NULL,
  `source_link` varchar(150) COLLATE utf8mb4_unicode_ci NOT NULL,
  `link` varchar(150) COLLATE utf8mb4_unicode_ci NOT NULL,
  `conversion` int(11) NOT NULL DEFAULT 0,
  `conversion_unique` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `sender_id_idx` (`user_id`),
  KEY `link_idx` (`link`),
  CONSTRAINT `links_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `views` (
  `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `created_at` datetime NOT NULL DEFAULT NOW() ON UPDATE CURRENT_TIMESTAMP,
  `link_id` int(11) NOT NULL,
  `ip` varchar(16) NOT NULL,
  `user-agent` varchar(50) NOT NULL,
  FOREIGN KEY (`link_id`) REFERENCES `links` (`id`)
);