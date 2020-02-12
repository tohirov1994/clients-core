package core

import (
	"database/sql"
	"errors"
	"fmt"
	DSN "github.com/tohirov1994/database"
)

var ErrorPassword = errors.New("password is not valid")

type Atm struct {
	Id       int64
	City     string
	District string
	Street   string
}

type Card struct {
	Id         int
	PAN        int
	PIN        int
	Balance    int
	HolderName string
	CVV        int
	Validity   int
}

type ServicesStruct struct {
	Id      int
	Service string
}

func Init(db *sql.DB) (err error) {
	initDDLsDMLs := []string{DSN.ManagersDDL, DSN.ClientsDDL, DSN.ClientsCardsDDL, DSN.AtmsDDL, DSN.ServicesDDL,
		DSN.ManagersDML, DSN.ClientsDML, DSN.ClientsCardsDML, DSN.AtmsDML, DSN.ServicesDML}
	for _, init := range initDDLsDMLs {
		_, err = db.Exec(init)
		if err != nil {
			return err
		}
	}
	return nil
}

func SignIn(loginUsr, passwordUsr string, db *sql.DB) (int, bool, error) {
	var dbLogin, dbPassword string
	var ClientId int
	err := db.QueryRow(DSN.GetLoginPassIdClient, loginUsr).Scan(&dbLogin, &dbPassword, &ClientId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, err
	}
	if dbPassword != passwordUsr {
		return 0, false, ErrorPassword
	}
	return ClientId, true, nil
}

func GetCurrentBalanceClientId(clientCardId int, db *sql.DB) (balance int, err error) {
	var idClient int
	err = db.QueryRow(DSN.GetBalanceClientId, clientCardId).Scan(&idClient, &balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, err
		}
		return 0, err
	}
	return balance, nil
}

func CheckPan(panClient int64, db *sql.DB) (result int64, err error) {
	var checker int64
	err = db.QueryRow(DSN.CheckPAN, panClient).Scan(&checker)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, err
		}
		return 0, err
	}
	return checker, nil
}

func GetCurrentBalanceClientPAN(clientPAN int64, db *sql.DB) (balance int, err error) {
	err = db.QueryRow(DSN.GetBalanceClientPAN, clientPAN).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, err
		}
		return 0, err
	}
	return balance, nil
}

func GetTransferCard(id int, db *sql.DB) (count int, err error) {
	err = db.QueryRow(DSN.GetTransferCard, id).Scan(&count)
	if err != nil {
		err := fmt.Errorf("can't find last Account Number %e", err)
		return 0, err
	}
	return count, nil
}

func SelectCards(id int, panCheck int64, db *sql.DB) (panAccept int64, err error) {
	err = db.QueryRow(DSN.SelectCardWhoHaveManyCards, id, panCheck).Scan(&panAccept)
	if err != nil {
		err := fmt.Errorf("can't select your card %e", err)
		return 0, err
	}
	return panAccept, nil
}

func OneCard(panReceiver int64, idSender, amount int, db *sql.DB) (status bool, err error) {
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	_, err = tx.Exec(
		DSN.OutOneAmmount,
		sql.Named("amount", amount),
		sql.Named("idClient", idSender),
	)
	_, err = tx.Exec(
		DSN.InAmmount,
		sql.Named("amount", amount),
		sql.Named("PANInner", panReceiver),
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func MoreCard(panSender, panReceiver int64, amount int, db *sql.DB) (status bool, err error) {
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	_, err = tx.Exec(
		DSN.OutMoreOneAmmount,
		sql.Named("amount", amount),
		sql.Named("panClient", panSender),
	)
	_, err = tx.Exec(
		DSN.InAmmount,
		sql.Named("amount", amount),
		sql.Named("PANInner", panReceiver),
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func CheckServiceName(Name string, db *sql.DB) (result string, err error) {
	var checker string
	err = db.QueryRow(DSN.CheckServiceName, Name).Scan(&checker)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", err
	}
	return checker, nil
}

func ServicesPayOneCard(nameService string, payerId, amount int, db *sql.DB) (result bool, err error) {
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	_, err = tx.Exec(
		DSN.OutOneAmmount,
		sql.Named("amount", amount),
		sql.Named("idClient", payerId),
	)
	_, err = tx.Exec(
		DSN.PayService,
		sql.Named("amount", amount),
		sql.Named("serviceName", nameService),
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func ServicesPayMoreCard(nameService string, cardPAN int64, amount int, db *sql.DB) (result bool, err error) {
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	_, err = tx.Exec(
		DSN.OutMoreOneAmmount,
		sql.Named("amount", amount),
		sql.Named("panClient", cardPAN),
	)
	_, err = tx.Exec(
		DSN.PayService,
		sql.Named("amount", amount),
		sql.Named("serviceName", nameService),
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func ATMsGet(db *sql.DB) (atms []Atm, err error) {
	rows, err := db.Query(DSN.GetATMData)
	if err != nil {
		return nil, err
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			atms = nil
		}
	}()

	for rows.Next() {
		atm := Atm{}
		err = rows.Scan(&atm.Id, &atm.City, &atm.District, &atm.Street)
		if err != nil {
			return nil, err
		}
		atms = append(atms, atm)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return atms, nil
}

func CardsGet(id int, db *sql.DB) (cards []Card, err error) {
	rows, err := db.Query(DSN.GetCards, id)
	if err != nil {
		return nil, err
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			cards = nil
		}
	}()
	for rows.Next() {
		card := Card{}
		err = rows.Scan(&card.Id, &card.PAN, &card.PIN, &card.Balance, &card.HolderName, &card.CVV, &card.Validity)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return cards, nil
}

func GetAllService(db *sql.DB) (services []ServicesStruct, err error) {
	rows, err := db.Query(DSN.GetServices)
	if err != nil {
		return nil, err
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			services = nil
		}
	}()
	for rows.Next() {
		service := ServicesStruct{}
		err = rows.Scan(&service.Id, &service.Service)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return services, nil
}
