CREATE TABLE packages
(
    `id`                     INT(11) NOT NULL AUTO_INCREMENT,
    `name`                   VARCHAR(64)   NOT NULL,
    `description`            TEXT          NOT NULL,
    `price_eur`              DECIMAL(6, 2) NOT NULL,
    `data_retention_minutes` INT(11) NOT NULL,
    `period_in_days`         INT(11) NOT NULL,
    `feature_basic`          BOOLEAN       NOT NULL DEFAULT FALSE,
    `feature_pro`            BOOLEAN       NOT NULL DEFAULT FALSE,
    `feature_expert`         BOOLEAN       NOT NULL DEFAULT FALSE,
    `is_active`              BOOLEAN       NOT NULL DEFAULT TRUE,
    `created_at`             INT(11) NOT NULL,
    `updated_at`             INT(11) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB;