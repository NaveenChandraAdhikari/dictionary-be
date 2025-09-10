package sqlconnect

import (
	"database/sql"
	"restapi/internal/models"
	"restapi/internal/pkg/utils"
)

func FindOrCreateOAuthUser(provider, oauthID, email, firstName, lastName, username string) (*models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "database connection error")
	}
	defer db.Close()

	user := &models.Exec{}
	selectQuery := `SELECT id, first_name, last_name, email, username, oauth_provider, oauth_id, is_oauth_user, role, inactive_status 
                    FROM execs WHERE oauth_provider=$1 AND oauth_id=$2`

	err = db.QueryRow(selectQuery, provider, oauthID).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username,
		&user.OAuthProvider, &user.OAuthID, &user.IsOAuthUser, &user.Role, &user.InactiveStatus,
	)

	if err == nil {
		return user, nil
	}

	if err != sql.ErrNoRows {
		return nil, utils.ErrorHandler(err, "database query error")
	}

	oauthPasswordPlaceholder := "NO_PASSWORD_OAUTH_USER"

	insertQuery := `INSERT INTO execs (first_name, last_name, email, hashed_password, username, oauth_provider, oauth_id, is_oauth_user, role, inactive_status) 
                    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	newUser := &models.Exec{
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		Username:       username,
		IsOAuthUser:    true,
		Role:           "user",
		InactiveStatus: false,
	}

	err = db.QueryRow(insertQuery,
		newUser.FirstName,
		newUser.LastName,
		newUser.Email,
		oauthPasswordPlaceholder,
		newUser.Username,
		provider,
		oauthID,
		newUser.IsOAuthUser,
		newUser.Role,
		newUser.InactiveStatus,
	).Scan(&newUser.ID)
	if err != nil {
		return nil, utils.ErrorHandler(err, "error creating OAuth user")
	}

	newUser.OAuthProvider = sql.NullString{String: provider, Valid: true}
	newUser.OAuthID = sql.NullString{String: oauthID, Valid: true}

	return newUser, nil
}
