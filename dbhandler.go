package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "legouser:legopassword@dbserver/legoparts")
	if err != nil {
		panic(err.Error())
	}
	return db
}

// CreateDatabase creates the tables for the application
func createDatabase() error {
	db := dbConn()

	statements := [4]string{
		`CREATE TABLE IF NOT EXISTS legoset (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    description TEXT
)  ENGINE=INNODB;
`,
		`CREATE TABLE IF NOT EXISTS legopart (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    partnumber VARCHAR(32) NOT NULL,
    description TEXT,
    legocolor_id INT NOT NULL,
    legoset_id INT NOT NULL,
    requiredqty INT,
    foundqty INT
)  ENGINE=INNODB;
`,
		`CREATE TABLE IF NOT EXISTS legocolor (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    colornumber INT NOT NULL,
    colorname VARCHAR(255) NOT NULL
)  ENGINE=INNODB;
`}

	for _, stmt := range statements {
		fmt.Println(stmt)
		_, err := db.Exec(stmt)
		if err != nil {
			return err
		}
	}
	defer db.Close()
	return nil
}

// GetSetList gets the list of sets
func getSetList() ([]LegoSet, error) {
	db := dbConn()

	rows, err := db.Query("SELECT id, name, description from legoset")
	if err != nil {
		return nil, err
	}

	result := []LegoSet{}

	for rows.Next() {
		var ls LegoSet
		err = rows.Scan(&ls.ID, &ls.Name, &ls.Description)
		if err != nil {
			return nil, err
		}
		result = append(result, ls)
	}
	defer db.Close()
	return result, nil
}

// GetPartList gets the list of parts for the set
func getPartList(setID int) ([]LegoPart, error) {
	db := dbConn()

	rows, err := db.Query("SELECT partnumber,description,legocolor_id,legoset_id,requiredqty,foundqty FROM legopart WHERE legoset_id = ?", setID)
	if err != nil {
		return nil, err
	}

	result := []LegoPart{}

	for rows.Next() {
		var lp LegoPart
		err = rows.Scan(&lp.Partnumber, &lp.Description, &lp.ColorID, &lp.SetID, &lp.RequiredQty, &lp.FoundQty)
		if err != nil {
			return nil, err
		}
		result = append(result, lp)
	}
	defer db.Close()
	return result, nil
}

// SaveSet saves the set
func saveSet(set LegoSet) int {
	db := dbConn()

	insert, err := db.Prepare("INSERT INTO legoset(name, description) VALUES (?, ?)")
	if err != nil {
		panic(err.Error())
	}
	res, err := insert.Exec(set.Name, set.Description)
	if err != nil {
		panic(err.Error())
	}
	id64, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	return int(id64)
}

// SaveParts saves an array of parts
func saveParts(parts []LegoPart) {
	db := dbConn()

	insert, err := db.Prepare("INSERT INTO legopart(partnumber, description, legocolor_id, legoset_id, requiredqty, foundqty) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	for _, part := range parts {
		_, err := insert.Exec(part.Partnumber, part.Description, part.ColorID, part.SetID, part.RequiredQty, 0)
		if err != nil {
			panic(err.Error())
		}
	}
	defer db.Close()
}

// UpdatePart updates the part
func updatePart(part LegoPart) {
	db := dbConn()

	insert, err := db.Prepare("UPDATE legopart SET partnumber=?, description=?, legocolor_id=?, legoset_id=?, requiredqty=?, foundqty=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	_, err = insert.Exec(part.Partnumber, part.Description, part.ColorID, part.SetID, part.RequiredQty, part.FoundQty)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
}

func setPartFoundQuantity(setID, partNumber, colorID, direction string) (int, error) {
	db := dbConn()

	var statement string

	if direction == "up" {
		statement = "UPDATE legopart SET foundqty=foundqty+1 WHERE foundqty<requiredqty AND partnumber=? AND legocolor_id=? AND legoset_id=?"
	} else {
		statement = "UPDATE legopart SET foundqty=foundqty-1 WHERE foundqty>0 AND partnumber=? AND legocolor_id=? AND legoset_id=?"
	}
	update, err := db.Prepare(statement)
	if err != nil {
		return 0, err
	}

	res, err := update.Exec(partNumber, colorID, setID)
	if err != nil {
		return 0, err
	}
	_, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT foundqty FROM legopart WHERE partnumber=? AND legocolor_id=? AND legoset_id=?", partNumber, colorID, setID)
	//rows, err := db.Query("SELECT requiredqty,foundqty FROM legopart WHERE partnumber=? AND legocolor_id=? AND legoset_id=?", partNumber, colorID, setID)
	if err != nil {
		return 0, err
	}

	//var requiredqty int
	var foundqty int
	if rows.Next() {
		// err = rows.Scan(&requiredqty, &foundqty)
		err = rows.Scan(&foundqty)
		if err != nil {
			return 0, err
		}
	}

	return foundqty, nil
}

// GetLegoColors gets the list of colors
func getLegoColors() (map[int]string, error) {
	db := dbConn()

	rows, err := db.Query("SELECT colornumber, colorname FROM legocolor")
	if err != nil {
		return nil, err
	}

	result := make(map[int]string)

	for rows.Next() {
		var lc LegoColor
		err = rows.Scan(&lc.Number, &lc.Name)
		if err != nil {
			continue
		}
		result[lc.Number] = lc.Name
	}
	defer db.Close()
	return result, nil
}

// SaveColors saves an array of colors
func saveColors(colors []LegoColor) {
	db := dbConn()

	insert, err := db.Prepare("INSERT INTO legocolor (colornumber, colorname) VALUES (?, ?)")
	if err != nil {
		panic(err.Error())
	}

	for _, color := range colors {
		_, err := insert.Exec(color.Number, color.Name)
		if err != nil {
			panic(err.Error())
		}
	}
	defer db.Close()
}
