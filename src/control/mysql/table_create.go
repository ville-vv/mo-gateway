package mysql

import "database/sql"

func tableCreateSql() map[string]string {
	tables := make(map[string]string)
	tables["account"] = `
CREATE TABLE account(
	id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
	user_no BIGINT NOT NULL,
	phone varchar(15) NOT NULL,
	email varchar(50) NOT NULL,
	create_at timestamp NOT NULL default CURRENT_TIMESTAMP,
	update_at timestamp NOT NULL default CURRENT_TIMESTAMP,
	UNIQUE KEY uq_account_user_no (user_no),
	KEY idx_account_phone (phone),
	KEY idx_account_email (email)
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
`
	return tables
}

func getTables(db *sql.DB) (tbList map[string]bool, err error) {
	res, err := db.Query("show tables;")
	if err != nil {
		return
	}
	tbList = make(map[string]bool)
	for res.Next() {
		tb := ""
		if err = res.Scan(&tb); err != nil {
			return
		}
		tbList[tb] = true
	}
	return
}

func getDatabases(db *sql.DB) (dbList map[string]bool, err error) {
	res, err := db.Query("show databases;")
	if err != nil {
		return
	}
	dbList = make(map[string]bool)
	for res.Next() {
		name := ""
		if err = res.Scan(&name); err != nil {
			return
		}
		dbList[name] = true
	}
	return
}
