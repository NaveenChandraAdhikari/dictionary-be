package sqlconnect

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"restapi/internal/models"
	"restapi/internal/pkg/utils"
)

func AddWordDBHandler(newWord models.Word, createdBy int) (models.Word, error) {

	db, err := ConnectDb()
	if err != nil {
		return models.Word{}, utils.ErrorHandler(err, "Error adding word")
	}
	defer db.Close()
	newWord.CreatedBy = createdBy
	//meaningsJSON, err := json.Marshal(newWord.Meanings)
	if err != nil {
		return models.Word{}, utils.ErrorHandler(err, "error marshalling meanings")
	}
	query := utils.GenerateInsertQuery("words", models.Word{}) + " RETURNING id"
	values := utils.GetStructValues(newWord)
	err = db.QueryRow(query, values...).Scan(&newWord.ID)
	if err != nil {
		return models.Word{}, utils.ErrorHandler(err, "Error adding word")
	}
	//TODO if a user add same word twice please make the error handling specific
	return newWord, nil
}

// TODO USER SOFT DELETE // HARD DELETE FOR THE WORD AND MEANING ......
func DeleteWordDBHandler(ids []int, userID int) ([]int, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error deleting Id")

	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error starting transaction")

	}

	stmt, err := tx.Prepare("DELETE FROM words WHERE id =$1 AND created_by=$2")
	if err != nil {

		tx.Rollback()
		return nil, utils.ErrorHandler(err, "Error preparing delete word")
	}
	defer stmt.Close()

	deletedIds := []int{}

	for _, id := range ids {
		result, err := stmt.Exec(id, userID)

		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "Error executing delete Id")
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "error checking delete result ")
		}
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}
		if rowsAffected < 1 {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, fmt.Sprintf("id %s not found", id))
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error committing delete Id")

	}
	if len(deletedIds) < 1 {
		return nil, utils.ErrorHandler(err, "No Ids found to delete")
	}
	return deletedIds, nil
}
func DeleteOneWordDBHandler(id int, userId int) error {

	db, err := ConnectDb()
	if err != nil {
		return utils.ErrorHandler(err, "Error deleting id")
	}
	result, err := db.Exec("DELETE FROM words WHERE id=$1 AND created_by=$2", id, userId)
	if err != nil {
		return utils.ErrorHandler(err, "Error deleting id")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "Error checking delete result ")
	}
	if rowsAffected == 0 {
		return utils.ErrorHandler(err, "id not found")
	}
	return nil
}
func UpdatedWordDBHandler(id int, updatedWord models.Word, userId int) (models.Word, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Word{}, utils.ErrorHandler(err, "Error updating id")
	}
	defer db.Close()
	//fetch existing word
	var existingWord models.Word

	var meaningsJSON []byte

	query := `SELECT id,word,phonetic,origin,meanings,created_by FROM words WHERE id=$1 AND created_by=$2`
	//fmt.Println("i am the cause", existingWord.Meanings)
	err = db.QueryRow(query, id, userId).Scan(&existingWord.ID, &existingWord.Word, &existingWord.Phonetic, &existingWord.Origin, &meaningsJSON, &existingWord.CreatedBy)
	//fmt.Println("i am the cause 2", existingWord.Meanings)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Word{}, utils.ErrorHandler(err, "id not found or not owned")
		}
		return models.Word{}, utils.ErrorHandler(err, "Error updating id")

	}

	//unmarshal meaningsSON into slice........now we have it as []string
	if err := json.Unmarshal(meaningsJSON, &existingWord.Meanings); err != nil {
		return models.Word{}, utils.ErrorHandler(err, "Error unmarshalling meanings")
	}

	updatedWord.ID = existingWord.ID
	updatedWord.CreatedBy = existingWord.CreatedBy

	//marshal meanings -> convert back before saving
	meaningsToSave, err := json.Marshal(updatedWord.Meanings)
	if err != nil {
		return models.Word{}, utils.ErrorHandler(err, "Error marshalling meanings")
	}

	updateQuery := `UPDATE words SET word=$1,phonetic=$2,origin=$3,meanings=$4 WHERE id=$5 AND created_by=$6`
	_, err = db.Exec(updateQuery, updatedWord.Word, updatedWord.Phonetic, updatedWord.Origin, meaningsToSave, updatedWord.ID, updatedWord.CreatedBy)
	if err != nil {
		return models.Word{}, utils.ErrorHandler(err, "Error updating word(duplicate word?)")
	}
	return updatedWord, nil
}

//	func GetWordsDBHandler(r *http.Request, userID, limit, page int) ([]models.Word, int, error) {
//		db, err := ConnectDb()
//		if err != nil {
//			return nil, 0, utils.ErrorHandler(err, "error retrieving words")
//		}
//		defer db.Close()
//
//		query := "SELECT id, word, phonetic, origin, meanings, created_by, created_at, updated_at FROM words WHERE 1=1 "
//		var args []interface{}
//
//		query, args = utils.AddFilters(r, query, args)
//
//		// Add pagination (after filters, so $n continues)
//		paramCount := len(args) + 1
//		offset := (page - 1) * limit
//		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
//		args = append(args, limit, offset)
//
//		query = utils.AddSorting(r, query) // Sorting after LIMIT/OFFSET (Postgres allows)
//
//		rows, err := db.Query(query, args...)
//		if err != nil {
//			return nil, 0, utils.ErrorHandler(err, "error retrieving words")
//		}
//		defer rows.Close()
//
//		var words []models.Word
//		for rows.Next() {
//			var word models.Word
//			var meaningsBytes []byte
//			err := rows.Scan(&word.ID, &word.Word, &word.Phonetic, &word.Origin, &meaningsBytes, &word.CreatedBy, &word.CreatedAt, &word.UpdatedAt)
//			if err != nil {
//				return nil, 0, utils.ErrorHandler(err, "error scanning words")
//			}
//			err = json.Unmarshal(meaningsBytes, &word.Meanings)
//			if err != nil {
//				return nil, 0, utils.ErrorHandler(err, "error unmarshalling meanings")
//			}
//			words = append(words, word)
//		}
//
//		// Total count (separate query without LIMIT/OFFSET)
//		countQuery := "SELECT COUNT(*) FROM words WHERE 1=1 "
//		countQuery, countArgs := utils.AddFilters(r, countQuery, nil)
//		var totalWords int
//		err = db.QueryRow(countQuery, countArgs...).Scan(&totalWords)
//		if err != nil {
//			return nil, 0, utils.ErrorHandler(err, "error counting words")
//		}
//
//		return words, totalWords, nil
//	}
func GetWordsDBHandler(r *http.Request, userID, limit, page int) ([]models.Word, int, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, 0, utils.ErrorHandler(err, "error retrieving words")
	}
	defer db.Close()

	query := "SELECT id, word, phonetic, origin, meanings, created_by, created_at, updated_at FROM words WHERE 1=1"
	var args []interface{}

	query, args = utils.AddFilters(r, query, args)

	// Add filtering for userID ownership
	args = append(args, userID)
	query += fmt.Sprintf(" AND created_by = $%d", len(args))

	query = utils.AddSorting(r, query)

	paramCount := len(args) + 1
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
	args = append(args, limit, (page-1)*limit)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, utils.ErrorHandler(err, "error retrieving words")
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var word models.Word
		var meaningsBytes []byte
		err := rows.Scan(&word.ID, &word.Word, &word.Phonetic, &word.Origin, &meaningsBytes, &word.CreatedBy, &word.CreatedAt, &word.UpdatedAt)
		if err != nil {
			return nil, 0, utils.ErrorHandler(err, "error scanning words")
		}
		err = json.Unmarshal(meaningsBytes, &word.Meanings)
		if err != nil {
			return nil, 0, utils.ErrorHandler(err, "error unmarshalling meanings-")
		}
		words = append(words, word)
	}

	// Count total words for this user with filters (without limit/offset)
	countQuery := "SELECT COUNT(*) FROM words WHERE 1=1"
	var countArgs []interface{}
	countQuery, countArgs = utils.AddFilters(r, countQuery, countArgs)

	// Add userID filter for count too
	countArgs = append(countArgs, userID)
	countQuery += fmt.Sprintf(" AND created_by = $%d", len(countArgs))

	var totalWords int
	err = db.QueryRow(countQuery, countArgs...).Scan(&totalWords)
	if err != nil {
		return nil, 0, utils.ErrorHandler(err, "error counting words")
	}

	return words, totalWords, nil
}

func GetOneWordDBHandler(word string) (models.Word, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Word{}, utils.ErrorHandler(err, "error retrieving word")
	}
	defer db.Close()

	var w models.Word
	var meaningsBytes []byte
	query := "SELECT id, word, phonetic, origin, meanings, created_by, created_at, updated_at FROM words WHERE word = $1"
	err = db.QueryRow(query, word).Scan(&w.ID, &w.Word, &w.Phonetic, &w.Origin, &meaningsBytes, &w.CreatedBy, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Word{}, utils.ErrorHandler(err, "word not found")
		}
		return models.Word{}, utils.ErrorHandler(err, "error retrieving word")
	}

	err = json.Unmarshal(meaningsBytes, &w.Meanings)
	if err != nil {
		return models.Word{}, utils.ErrorHandler(err, "error unmarshalling meanings")
	}

	return w, nil
}
