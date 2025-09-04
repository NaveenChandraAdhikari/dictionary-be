package sqlconnect

import (
	"database/sql"
	"fmt"
	"restapi/internal/models"
	"restapi/internal/pkg/utils"
)

//
//func AddExecDBHandler(newExec models.Exec) (models.Exec, error) {
//	db, err := ConnectDb()
//	if err != nil {
//		return models.Exec{}, utils.ErrorHandler(err, "error adding data")
//
//	}
//	defer db.Close()
//
//	stmt, err := db.Prepare(utils.GenerateInsertQuery("execs", models.Exec{}))
//	{
//		if err != nil {
//			return models.Exec{}, utils.ErrorHandler(err, "error adding data")
//		}
//	}
//	defer stmt.Close()
//
//	newExec.HashedPassword, err = utils.HashPassword(newExec.HashedPassword)
//	if err != nil {
//		return models.Exec{}, utils.ErrorHandler(err, "error adding exec into database")
//	}
//
//	values := utils.GetStructValues(newExec)
//	res, err := stmt.Exec(values...)
//	if err != nil {
//		return models.Exec{}, utils.ErrorHandler(err, "error adding data")
//
//	}
//	lastId, err := res.LastInsertId()
//	if err != nil {
//		return models.Exec{}, utils.ErrorHandler(err, "error adding data")
//	}
//	newExec.ID = int(lastId)
//	return newExec, nil
//
//}

func AddExecDBHandler(newExec models.Exec) (models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "error adding data")
	}
	defer db.Close()

	//  RETURNING id to the query
	query := utils.GenerateInsertQuery("execs", models.Exec{}) + " RETURNING id"

	newExec.HashedPassword, err = utils.HashPassword(newExec.HashedPassword)
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "error adding exec into database")
	}

	values := utils.GetStructValues(newExec)

	// Execute and scan the returned ID,queryrow returns 1 row
	err = db.QueryRow(query, values...).Scan(&newExec.ID)
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "error adding data")
	}

	return newExec, nil
}

func GetUserByUsername(username string) (*models.Exec, error) {
	db, err := ConnectDb()

	//if no user found or what so error then i return nil here
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}
	defer db.Close()

	user := &models.Exec{}
	query := `SELECT id,first_name,last_name,username,inactive_status,role,hashed_password FROM execs WHERE username=$1`
	//fmt.Println("THIS IS THE QUERY", query)
	err = db.QueryRow(query, username).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.InactiveStatus, &user.Role, &user.HashedPassword)
	fmt.Println("userrrr", user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrorHandler(err, "user not found")

		}
		return nil, utils.ErrorHandler(err, "database error")
	}
	return user, nil

}
