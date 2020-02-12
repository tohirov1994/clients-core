package core

import (
	"database/sql"
	"errors"
	_ "errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	DSN "github.com/tohirov1994/database"
	"testing"
)

const dbDriver = "sqlite3"
const dbMemory = ":memory:"

func TestSignIn_QueryError(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, _, err = SignIn("", "", db)
	if err == nil {
		t.Errorf("can't execute SignIn: %v", err)
	}

}

func TestInit_CanNotApply(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_ = db.Close()
	err = Init(db)
	if err == nil {
		t.Errorf("Init just not be nil: %v", err)
	}
}

func TestInit_Apply(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	err = Init(db)
	initDDLsDMLs := []string{DSN.ManagersDDL, DSN.ClientsDDL, DSN.ClientsCardsDDL, DSN.AtmsDDL, DSN.ServicesDDL,
		DSN.ManagersDML, DSN.ClientsDML, DSN.ClientsCardsDML, DSN.AtmsDML, DSN.ServicesDML}
	for _, init := range initDDLsDMLs {
		_, err = db.Exec(init)
		if err != nil {
			t.Errorf("can't init db: %v", err)
		}
	}
	if err != nil {
		t.Errorf("Init apply, just nil: %v", err)
	}
}

func TestSignIn_NoValidData(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
	CREATE TABLE clients (
    Id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL)`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	_, err = db.Exec(`INSERT INTO clients(Id, login, password) VALUES (1, 'jack', 'password'),
(2, 'max', 'password'),
(3, 'nilson', 'password'),
(4, 'paterson', 'password')`)
	if err != nil {
		t.Errorf("can't execute insert login and password to DB: %v", err)
	}
	id, result, err := SignIn("kayla", "secret", db)
	if err != nil {
		t.Errorf("can't execute SignIn: %v", err)
	}
	if result == true {
		t.Error("signIn result must be false with incorrect")
	}
	if id != 0 {
		t.Errorf("must be 0: %d", id)
	}
}

func TestSignIn_LoginNotOkForInvalidPassword(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
	CREATE TABLE clients (
    Id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL)`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	_, err = db.Exec(`INSERT INTO clients(Id, login, password) VALUES (1, 'jack', 'password')`)
	if err != nil {
		t.Errorf("can't execute insert login and password to DB: %v", err)
	}
	_, _, err = SignIn("jack", "12345", db)
	if !errors.Is(err, ErrorPassword) {
		t.Errorf("Error for invalid pass: %v", err)
	}
}

func TestSignIn_OK(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`
	CREATE TABLE clients (
    Id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL)`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}

	_, err = db.Exec(`INSERT INTO clients(Id, login, password) VALUES (1, 'jack', 'password'),
(2, 'max', 'password'),
(3, 'nilson', 'password'),
(4, 'paterson', 'password')`)
	if err != nil {
		t.Errorf("can't execute insert login and password to DB: %v", err)
	}

	id, result, err := SignIn("nilson", "password", db)
	if err != nil {
		t.Errorf("can't execute SignIn: %v", err)
	}
	if result != true {
		t.Error("signIn result must be true for existing account")
	}
	if id != 3 {
		t.Errorf("must be 3: %d", id)
	}
}

func TestGetCurrentBalanceClientId_NoTable(t *testing.T)  {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	result, err := GetCurrentBalanceClientId(159, db)
	if err == nil {
		t.Errorf("can't execute query from emty base: %v", err)
	}
	if result != 0 {
		t.Errorf("Result just be zero: %d", result)
	}
}

func TestGetCurrentBalanceClientId_ErrorRows(t *testing.T)  {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
	CREATE TABLE clients_cards (
	client_id integer NOT NULL,
	balance INTEGER);`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	_, err = db.Exec(`
	INSERT INTO clients_cards VALUES (166, 5000);`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	result, err := GetCurrentBalanceClientId(159, db)
	if err == nil {
		t.Errorf("can't execute query from emty base: %v", err)
	}
	if result != 0 {
		t.Errorf("balance just be zero: %d", result)
	}
}

func TestGetCurrentBalanceClientId_OK(t *testing.T)  {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
	CREATE TABLE clients_cards (
	client_id integer NOT NULL,
	balance INTEGER);`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	_, err = db.Exec(`
	INSERT INTO clients_cards VALUES (166, 5000);`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	balance, err := GetCurrentBalanceClientId(166, db)
	if err != nil {
		t.Errorf("can't query GetBalanceFromClientId: %v", err)
	}
	if balance != 5000 {
		t.Errorf("balance just be 5000: %d", balance)
	}
}

func TestCheckPAN_NoTable(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	result, err := CheckPan(123, db)
	if err == nil {
		t.Errorf("can't execute query from emty base: %v", err)
	}
	if result != 0 {
		t.Errorf("Result just be zero: %d", result)
	}
}

func TestCheckPAN_RowsError(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
	CREATE TABLE clients_cards (
	pan integer NOT NULL UNIQUE);`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	res, err := CheckPan(159, db)
	if err == nil {
		t.Errorf("error just not be nil: %d", err)
	}
	if res != 0 {
		t.Errorf("just be zero: %d", res)
	}
}

func TestCheckPAN_OK(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
	CREATE TABLE clients_cards (
	pan integer NOT NULL UNIQUE);`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	_, err = db.Exec(`
INSERT INTO clients_cards VALUES (222);`)
	if err != nil {
		t.Errorf("can't execute insert login and password to DB: %v", err)
	}
	result, err := CheckPan(222, db)
	if err != nil {
		t.Errorf("can't execute checkPAN: %v", err)
	}
	if result != 222 {
		t.Errorf("just be 222 : %d", result)
	}
}

func TestGetCurrentBalanceClientPAN_NoTable(t *testing.T)  {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	result, err := GetCurrentBalanceClientPAN(992, db)
	if err == nil {
		t.Errorf("can't execute query from emty base: %v", err)
	}
	if result != 0 {
		t.Errorf("Result just be zero: %d", result)
	}
}

func TestGetCurrentBalanceClientPAN_ErrorRows(t *testing.T)  {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
	CREATE TABLE clients_cards (
	pan integer NOT NULL,
	balance INTEGER);`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	_, err = db.Exec(`
	INSERT INTO clients_cards VALUES (645, 5000);`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	result, err := GetCurrentBalanceClientPAN(159, db)
	if err == nil {
		t.Errorf("can't execute query from emty base: %v", err)
	}
	if result != 0 {
		t.Errorf("balance just be zero: %d", result)
	}
}

func TestGetCurrentBalanceClientPAN_OK(t *testing.T)  {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
	CREATE TABLE clients_cards (
	pan integer NOT NULL,
	balance INTEGER);`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	_, err = db.Exec(`
	INSERT INTO clients_cards VALUES (333, 100000);`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	balance, err := GetCurrentBalanceClientPAN(333, db)
	if err != nil {
		t.Errorf("can't query GetBalanceFromClientPAN: %v", err)
	}
	if balance != 100000 {
		t.Errorf("balance just be 100000: %d", balance)
	}
}

func TestGetTransferCard_NoTable(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	result, err := GetTransferCard(1, db)
	if err == nil {
		t.Errorf("can't execute get card: %v", err)
	}
	if result != 0 {
		t.Errorf("We just have zero cards: %v", result)
	}
}

func TestGetTransferCard_OK(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS clients_cards
(
    Id         INTEGER PRIMARY KEY AUTOINCREMENT,
    pan        INTEGER NOT NULL UNIQUE,
    balance    INTEGER NOT NULL,
    client_id  INTEGER NOT NULL REFERENCES clients
);`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	_, err = db.Exec(`INSERT INTO clients_cards(Id, pan, balance, client_id) 
VALUES (1, 1111, 1000000, 1),
(2, 2222, 2000000, 1),
(3, 3333, 3000000, 3),
(4, 4444, 352, 4),
(5, 5555, 5000000, 5);`)
	if err != nil {
		t.Errorf("can't execute insert card to DB: %v", err)
	}
	result, err := GetTransferCard(1, db)
	if err != nil {
		t.Errorf("can't execute get card: %v", err)
	}
	if result != 2 {
		t.Errorf("We just have 2 cards: %v", result)
	}
}

func TestSelectCards_NoSuchData(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	result, err := SelectCards(1, 2222, db)
	if err == nil {
		t.Errorf("can't execute get card: %v", err)
	}
	if result != 0 {
		t.Errorf("card not such, PAN just be zero: %v", result)
	}
}

func TestSelectCards_OK(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS clients_cards
(
    Id         INTEGER PRIMARY KEY AUTOINCREMENT,
    pan        INTEGER NOT NULL UNIQUE,
    balance    INTEGER NOT NULL,
    client_id  INTEGER NOT NULL REFERENCES clients
);`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	_, err = db.Exec(`INSERT INTO clients_cards(Id, pan, balance, client_id) 
VALUES (1, 1111, 1000000, 1),
(2, 2222, 2000000, 1),
(3, 3333, 3000000, 3),
(4, 4444, 352, 4),
(5, 5555, 5000000, 5);`)
	if err != nil {
		t.Errorf("can't execute insert card to DB: %v", err)
	}
	result, err := SelectCards(1, 2222, db)
	if err != nil {
		t.Errorf("can't execute get card: %v", err)
	}
	if result != 2222 {
		t.Errorf("PAN card just be 2222: %v", result)
	}
}

func TestOneCard_DbClose(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_ = db.Close()
	result, err := OneCard(4444, 2222, 2000000, db)
	if err == nil {
		t.Errorf("transefer money just have error: %v", err)
	}
	if result == true {
		t.Errorf("trasfer just be false: %v", result)
	}
}

func TestOneCard_CanNotBegin(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	result, err := OneCard(4444, 2222, 2000000, db)
	if err == nil {
		t.Errorf("transefer money just have error: %v", err)
	}
	if result == true {
		t.Errorf("trasfer just be false: %v", result)
	}
}

func TestOneCard_OK(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS clients_cards
(
    Id         INTEGER PRIMARY KEY AUTOINCREMENT,
    pan        INTEGER NOT NULL UNIQUE,
    balance    INTEGER NOT NULL,
    client_id  INTEGER NOT NULL REFERENCES clients
);`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	_, err = db.Exec(`INSERT INTO clients_cards(Id, pan, balance, client_id) 
VALUES (1, 1111, 1000000, 1),
(2, 2222, 2000000, 1),
(3, 3333, 3000000, 3),
(4, 4444, 352, 4),
(5, 5555, 5000000, 5);`)
	if err != nil {
		t.Errorf("can't execute insert card to DB: %v", err)
	}
	result, err := OneCard(4444, 2222, 2000000, db)
	if err != nil {
		t.Errorf("can't execute transefer money: %v", err)
	}
	if result != true {
		t.Errorf("trasfer just be true: %v", result)
	}
}

func TestMoreCard_DbClose(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_ = db.Close()
	result, err := MoreCard(4444, 2222, 2000000, db)
	if err == nil {
		t.Errorf("transefer money just have error: %v", err)
	}
	if result == true {
		t.Errorf("trasfer just be false: %v", result)
	}
}

func TestMoreCard_CanNotBegin(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	result, err := MoreCard(4444, 2222, 2000000, db)
	if err == nil {
		t.Errorf("transefer money just have error: %v", err)
	}
	if result == true {
		t.Errorf("trasfer just be false: %v", result)
	}
}

func TestMoreCard_OK(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS clients_cards
(
    Id         INTEGER PRIMARY KEY AUTOINCREMENT,
    pan        INTEGER NOT NULL UNIQUE,
    balance    INTEGER NOT NULL,
    client_id  INTEGER NOT NULL REFERENCES clients
);`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	_, err = db.Exec(`INSERT INTO clients_cards(Id, pan, balance, client_id) 
VALUES (1, 1111, 1000000, 1),
(2, 2222, 2000000, 1),
(3, 3333, 3000000, 3),
(4, 4444, 352, 4),
(5, 5555, 5000000, 5);`)
	if err != nil {
		t.Errorf("can't execute insert card to DB: %v", err)
	}
	result, err := MoreCard(5555, 4444, 200000, db)
	if err != nil {
		t.Errorf("can't execute transefer money: %v", err)
	}
	if result != true {
		t.Errorf("trasfer just be true: %v", result)
	}
}

func TestCheckServiceName_NoTable(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	result, err := CheckServiceName("name", db)
	if err == nil {
		t.Errorf("can't execute query from emty base: %v", err)
	}
	if result != "" {
		t.Error("Result signIn no be true, when records account is empty")
	}
}

func TestCheckServiceName_RowsError(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
	CREATE TABLE services (
	service TEXT NOT NULL UNIQUE);`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	_, err = CheckServiceName("internet", db)
	if err == nil {
		t.Errorf("Not ErrorPassword error for invalid pass: %v", err)
	}
}

func TestCheckServiceName_OK(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
	CREATE TABLE services (
	service TEXT NOT NULL UNIQUE);`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	_, err = db.Exec(`
INSERT INTO services
VALUES ('internet')
ON CONFLICT DO NOTHING;`)
	if err != nil {
		t.Errorf("can't execute insert login and password to DB: %v", err)
	}
	result, err := CheckServiceName("internet", db)
	if err != nil {
		t.Errorf("can't execute signIn: %v", err)
	}
	if result != "internet" {
		t.Error("signIn result must be true for existing account")
	}
}

func TestServicesPayOneCard_DbClose(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_ = db.Close()
	result, err := ServicesPayOneCard("phone", 2222, 2000000, db)
	if err == nil {
		t.Errorf("pay of service just have error: %v", err)
	}
	if result == true {
		t.Errorf("pay just be false: %v", result)
	}
}

func TestServicesPayOneCard_CanNotBegin(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	result, err := ServicesPayOneCard("internet", 2222, 2000000, db)
	if err == nil {
		t.Errorf("pay of service just have error: %v", err)
	}
	if result == true {
		t.Errorf("pay just be false: %v", result)
	}
}

func TestServicesPayOneCard_OK(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS services
(service TEXT    NOT NULL,
balance integer    NOT NULL);`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	_, err = db.Exec(`
INSERT INTO services
VALUES ('phone', 5000);`)
	if err != nil {
		t.Errorf("can't execute insert card to DB: %v", err)
	}
	result, err := ServicesPayOneCard("water", 4444, 200000, db)
	if err != nil {
		t.Errorf("can't execute pay service: %v", err)
	}
	if result != true {
		t.Errorf("trasfer just be true: %v", result)
	}
}

func TestServicesPayMoreCard_DbClose(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_ = db.Close()
	result, err := ServicesPayMoreCard("phone", 2222, 2000000, db)
	if err == nil {
		t.Errorf("pay of service just have error: %v", err)
	}
	if result == true {
		t.Errorf("pay just be false: %v", result)
	}
}

func TestServicesPayMoreCard_CanNotBegin(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	result, err := ServicesPayMoreCard("internet", 2222, 2000000, db)
	if err == nil {
		t.Errorf("pay of service just have error: %v", err)
	}
	if result == true {
		t.Errorf("pay just be false: %v", result)
	}
}

func TestServicesPayMoreCard_OK(t *testing.T) {
	db, err := sql.Open(dbDriver, dbMemory)
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS services
(service TEXT    NOT NULL,
balance integer    NOT NULL);`)
	if err != nil {
		t.Errorf("can't execute query to base: %v", err)
	}
	_, err = db.Exec(`
INSERT INTO services
VALUES ('phone', 5000);`)
	if err != nil {
		t.Errorf("can't execute insert card to DB: %v", err)
	}
	result, err := ServicesPayMoreCard("water", 4444, 200000, db)
	if err != nil {
		t.Errorf("can't execute pay service: %v", err)
	}
	if result != true {
		t.Errorf("trasfer just be true: %v", result)
	}
}

func ExampleATMsGet_WithoutData() {
	db, _ := sql.Open(dbDriver, dbMemory)
	result, _ := ATMsGet(db)
	fmt.Println(result)
	//Output: []
}

func ExampleATMsGet_OK() {
	db, _ := sql.Open(dbDriver, dbMemory)
	_, _ = db.Exec(`
CREATE TABLE IF NOT EXISTS atms
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    city     TEXT NOT NULL,
    district TEXT NOT NULL,
    street   TEXT NOT NULL
);`)
	_, _ = db.Exec(`
INSERT INTO atms
VALUES (1, 'Dushanbe', 'Somoni', 'Foteh51')
ON CONFLICT DO NOTHING;`)
	result, _ := ATMsGet(db)
	fmt.Println(result)
	//Output: [{1 Dushanbe Somoni Foteh51}]
}

func ExampleCardsGet_WithoutData() {
	db, _ := sql.Open(dbDriver, dbMemory)
	result, _ := CardsGet(0, db)
	fmt.Println(result)
	//Output: []
}

func ExampleCardsGet_OK() {
	db, _ := sql.Open(dbDriver, dbMemory)
	_, _ = db.Exec(`
CREATE TABLE clients_cards
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    pan        INTEGER NOT NULL UNIQUE,
    pin        INTEGER NOT NULL,
    balance    INTEGER NOT NULL,
    holderName TEXT    NOT NULL,
    cvv        INTEGER NOT NULL,
    validity   INTEGER NOT NULL,
    client_id  INTEGER NOT NULL REFERENCES clients
);`)
	_, _ = db.Exec(`
INSERT INTO clients_cards
VALUES (1, 2021600000000000, 1994, 1000000, 'ADMIN CLIENT', 333, 0222, 1);`)
	result, _ := CardsGet(1, db)
	fmt.Println(result)
	//Output: [{1 2021600000000000 1994 1000000 ADMIN CLIENT 333 222}]
}

func ExampleGetAllService_WithoutData() {
	db, _ := sql.Open(dbDriver, dbMemory)
	result, _ := GetAllService(db)
	fmt.Println(result)
	//Output: []
}

func ExampleGetAllService_OK() {
	db, _ := sql.Open(dbDriver, dbMemory)
	_, _ = db.Exec(`
CREATE TABLE services
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    service TEXT    NOT NULL
);`)
	_, _ = db.Exec(`
INSERT INTO services
VALUES (1, 'internet');`)
	result, _ := GetAllService(db)
	fmt.Println(result)
	//Output: [{1 internet}]
}




