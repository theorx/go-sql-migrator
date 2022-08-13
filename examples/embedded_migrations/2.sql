CREATE TABLE organizations
(
    `id`           INT(11)            NOT NULL AUTO_INCREMENT,
    `uuid`         VARCHAR(36) UNIQUE NOT NULL,
    `owner_id`     INT(11)            NOT NULL,
    `name`         VARCHAR(64)        NOT NULL,
    `is_read_only` BOOLEAN            NOT NULL DEFAULT FALSE,
    `is_active`    BOOLEAN            NOT NULL DEFAULT TRUE,
    `created_at`   INT(11)            NOT NULL,
    `updated_at`   INT(11)            NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `uuid_idx` (`uuid`)
) ENGINE = InnoDB;