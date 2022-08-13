package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/theorx/go-sql-migrator/pkg/migrator"
	"log"
)

/*
Explicit inline migrations example

The code can be executed multiple times, the migrator will keep track of the state of migrations by the migration id
Each new migration must have incremented migration id, for example next in this file would be with id 4.
*/
func main() {
	db, err := sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test_db")

	if err != nil {
		log.Println(err)
		return
	}

	m := migrator.NewMigrator(db.DB, log.Println)

	log.Println(m.Apply([]migrator.Migration{
		migrator.NewMigration(1, "Create users table", func(db migrator.SQLClient) error {
			if _, err := db.Exec("CREATE TABLE users ( " +
				"`id` INT(11) NOT NULL AUTO_INCREMENT ," +
				" `uuid` VARCHAR(36) UNIQUE NOT NULL ," +
				" `email` VARCHAR(64) NULL ," +
				" `steam_id_64` VARCHAR(17) NULL ," +
				" `name` VARCHAR(32) NULL ," +
				" `tos` BOOLEAN NOT NULL DEFAULT FALSE ," +
				" `last_login_at` INT(11) NOT NULL ," +
				" `is_admin` BOOLEAN NOT NULL DEFAULT FALSE ," +
				" `is_active` BOOLEAN NOT NULL DEFAULT TRUE ," +
				" `created_at` INT(11) NOT NULL ," +
				" `updated_at` INT(11) NOT NULL ," +
				" PRIMARY KEY (`id`), INDEX `email_idx` (`email`), INDEX `steam_idx` (`steam_id_64`)) ENGINE = InnoDB;"); err != nil {
				return err
			}

			return nil
		}),
		migrator.NewMigration(2, "Create organizations table", func(db migrator.SQLClient) error {
			if _, err := db.Exec("CREATE TABLE organizations ( " +
				"`id` INT(11) NOT NULL AUTO_INCREMENT ," +
				" `uuid` VARCHAR(36) UNIQUE NOT NULL ," +
				" `owner_id` INT(11) NOT NULL ," +
				" `name` VARCHAR(64) NOT NULL ," +
				" `is_read_only` BOOLEAN NOT NULL DEFAULT FALSE ," +
				" `is_active` BOOLEAN NOT NULL DEFAULT TRUE ," +
				" `created_at` INT(11) NOT NULL ," +
				" `updated_at` INT(11) NOT NULL ," +
				" PRIMARY KEY (`id`), INDEX `uuid_idx` (`uuid`)) ENGINE = InnoDB;"); err != nil {
				return err
			}

			return nil
		}),
		migrator.NewMigration(3, "Create packages table", func(db migrator.SQLClient) error {
			if _, err := db.Exec("CREATE TABLE packages ( " +
				"`id` INT(11) NOT NULL AUTO_INCREMENT ," +
				" `name` VARCHAR(64) NOT NULL ," +
				" `description` TEXT NOT NULL ," +
				" `price_eur` DECIMAL(6,2) NOT NULL ," +
				" `data_retention_minutes` INT(11) NOT NULL ," +
				" `period_in_days` INT(11) NOT NULL ," +
				" `feature_basic` BOOLEAN NOT NULL DEFAULT FALSE ," +
				" `feature_pro` BOOLEAN NOT NULL DEFAULT FALSE ," +
				" `feature_expert` BOOLEAN NOT NULL DEFAULT FALSE ," +
				" `is_active` BOOLEAN NOT NULL DEFAULT TRUE ," +
				" `created_at` INT(11) NOT NULL ," +
				" `updated_at` INT(11) NOT NULL ," +
				" PRIMARY KEY (`id`)) ENGINE = InnoDB;"); err != nil {
				return err
			}

			return nil
		}),
	}))
}
