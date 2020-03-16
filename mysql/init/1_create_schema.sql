USE `extreme`;

CREATE TABLE IF NOT EXISTS `rating` (
    id CHAR(36) PRIMARY KEY,
    rating DOUBLE NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `effect_point` (
    name VARCHAR(16) PRIMARY KEY,
    point DOUBLE NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `stamp_relation` (
    id_from CHAR(36),
    id_to CHAR(36),
    point DOUBLE NOT NULL,
    PRIMARY KEY (id_from, id_to),
    KEY `id_to_idx` (id_to)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `stamp` (
    id CHAR(36) PRIMARY KEY,
    used INT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `seen_channel` (
    id CHAR(36) PRIMARY KEY,
    last_processed_message DATETIME(6)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
