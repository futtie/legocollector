package database

import (
	"database/sql"
	"fmt"

	"github.com/futtie/legocollector/structs"

	_ "github.com/go-sql-driver/mysql"
)

const (
	maxOpenConnections = 1000
	maxConRetries      = 3  // Number of retries
	maxConTimeout      = 10 // this is seconds
)

// Database ...
type Database struct {
	ConnString	string
	DB		*sql.DB
}

// Init the config for the connection
func InitDatabase(connstring string) *Database {
	return &Database{
		ConnString: host,
	}
}

// ConnectToDatabase Connect to the database
func (client *Database) ConnectToDatabase(dsn string, maxTries int) error {
	if maxTries == 0 {
		return fmt.Errorf("sql: could not connect to db after %d tries", maxConRetries)
	}
	con, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Error(err)
		return err
	}
	con.SetMaxOpenConns(maxOpenConnections)
	client.DB = con
	if err := client.Ping(); err != nil {
		log.Error(err)
		return client.ConnectToDatabase(dsn, maxTries-1)
	}
	return client.Ping()
}

// IsConnected Check if we're connected to the database
func (client *Database) IsConnected() bool {
	if err := client.Ping(); err != nil {
		return false
	}
	return true
}

// Ping Health check for database connection
func (client *Database) Ping() error {
	if client.DB == nil {
		return fmt.Errorf("sql: connection is nil")
	}
	return client.DB.Ping()
}

// EnsureConnected Ensure we are connected to the database
func (client *Database) EnsureConnected() error {
	if !client.IsConnected() {
		collation := "utf8_general_ci"
		if client.Collation != "" {
			collation = client.Collation
		}
		return client.ConnectToDatabase(
			client.ConnString,
			maxConRetries,
		)
	}
	return nil
}

// CreateDatabase creates the tables for the application
func (client *Database) CreateDatabase() error {

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
		_, err := client.DB.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetSetList gets the list of sets
func (client *Database) GetSetList() ([]LegoSet, error) {

	rows, err := client.DB.Query("SELECT id, name, description from legoset")
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
	return result, nil
}

// GetPartList gets the list of parts for the set
func (client *Database) GetPartList(setID int) ([]LegoPart, error) {
	
	rows, err := client.DB.Query("SELECT partnumber,description,legocolor_id,legoset_id,requiredqty,foundqty,requiredqty=foundqty as lowpri FROM legopart WHERE legoset_id = ? order by lowpri;", setID)
	if err != nil {
		return nil, err
	}

	result := []LegoPart{}

	for rows.Next() {
		var lp LegoPart
		err = rows.Scan(&lp.Partnumber, &lp.Description, &lp.ColorID, &lp.SetID, &lp.RequiredQty, &lp.FoundQty, &lp.LowPriority)
		if err != nil {
			return nil, err
		}
		result = append(result, lp)
	}
	return result, nil
}

// SaveSet saves the set
func (client *Database) SaveSet(set LegoSet) int {

	insert, err := client.DB.Prepare("INSERT INTO legoset(name, description) VALUES (?, ?)")
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
	return int(id64)
}

// SaveParts saves an array of parts
func (client *Database) SaveParts(parts []LegoPart) {

	insert, err := client.DB.Prepare("INSERT INTO legopart(partnumber, description, legocolor_id, legoset_id, requiredqty, foundqty) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	for _, part := range parts {
		_, err := insert.Exec(part.Partnumber, part.Description, part.ColorID, part.SetID, part.RequiredQty, 0)
		if err != nil {
			panic(err.Error())
		}
	}
}

// UpdatePart updates the part
func (client *Database) UpdatePart(part LegoPart) {

	insert, err := client.DB.Prepare("UPDATE legopart SET partnumber=?, description=?, legocolor_id=?, legoset_id=?, requiredqty=?, foundqty=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	_, err = insert.Exec(part.Partnumber, part.Description, part.ColorID, part.SetID, part.RequiredQty, part.FoundQty)
	if err != nil {
		panic(err.Error())
	}
}

func (client *Database) SetPartFoundQuantity(setID, partNumber, colorID, direction string) (int, error) {

	var statement string

	if direction == "up" {
		statement = "UPDATE legopart SET foundqty=foundqty+1 WHERE foundqty<requiredqty AND partnumber=? AND legocolor_id=? AND legoset_id=?"
	} else {
		statement = "UPDATE legopart SET foundqty=foundqty-1 WHERE foundqty>0 AND partnumber=? AND legocolor_id=? AND legoset_id=?"
	}
	update, err := client.DB.Prepare(statement)
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

	rows, err := client.DB.Query("SELECT foundqty FROM legopart WHERE partnumber=? AND legocolor_id=? AND legoset_id=?", partNumber, colorID, setID)
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
func (client *Database) GetLegoColors() (map[int]string, error) {

	rows, err := client.DB.Query("SELECT colornumber, colorname FROM legocolor")
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
	return result, nil
}

// SaveColors saves an array of colors
func (client *Database) SaveColors(colors []LegoColor) {
	
	insert, err := client.DB.Prepare("INSERT INTO legocolor (colornumber, colorname) VALUES (?, ?)")
	if err != nil {
		panic(err.Error())
	}

	for _, color := range colors {
		_, err := insert.Exec(color.Number, color.Name)
		if err != nil {
			panic(err.Error())
		}
	}
}
