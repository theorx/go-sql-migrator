CREATE TABLE users
(
    `id`            INT(11)            NOT NULL AUTO_INCREMENT,
    `uuid`          VARCHAR(36) UNIQUE NOT NULL,
    `email`         VARCHAR(64)        NULL,
    `steam_id_64`   VARCHAR(17)        NULL,
    `name`          VARCHAR(32)        NULL,
    `tos`           BOOLEAN            NOT NULL DEFAULT FALSE,
    `last_login_at` INT(11)            NOT NULL,
    `is_admin`      BOOLEAN            NOT NULL DEFAULT FALSE,
    `is_active`     BOOLEAN            NOT NULL DEFAULT TRUE,
    `created_at`    INT(11)            NOT NULL,
    `updated_at`    INT(11)            NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `email_idx` (`email`),
    INDEX `steam_idx` (`steam_id_64`)
) ENGINE = InnoDB;