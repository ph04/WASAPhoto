package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetDatabaseComment(commentId uint32, dbUser DatabaseUser) (DatabaseComment, error) {
	dbComment := DatabaseCommentDefault()

	// get the comment from the database
	err := db.c.QueryRow(`
		SELECT id, user, date, photo, comment_body
		FROM Comment
		WHERE id=?
	`, commentId).Scan(&dbComment.Id, &dbComment.User.Id, &dbComment.Date, &dbComment.Photo.Id, &dbComment.CommentBody)

	if errors.Is(err, sql.ErrNoRows) {
		return dbComment, ErrCommentDoesNotExist
	}

	// get the user of the comment
	dbCommentUser, err := db.GetDatabaseUser(dbComment.User.Id)

	if err != nil {
		return dbComment, err
	}

	dbComment.User.Username = dbCommentUser.Username

	// // get the photo of the comment
	dbPhoto, err := db.GetDatabasePhoto(dbComment.Photo.Id, dbUser)

	if err != nil {
		return dbComment, err
	}

	dbComment.Photo = dbPhoto

	return dbComment, err
}

func (db *appdbimpl) InsertComment(dbComment *DatabaseComment) error {
	// insert the comment into the database
	res, err := db.c.Exec(`
		INSERT INTO Comment(user, photo, date, comment_body)
		VALUES (?, ?, ?, ?)
	`, dbComment.User.Id, dbComment.Photo.Id, dbComment.Date, dbComment.CommentBody)

	if err != nil {
		return err
	}

	// get the comment id
	dbCommentId, err := res.LastInsertId()

	if err != nil {
		return err
	}

	dbComment.Id = uint32(dbCommentId)

	return nil
}

func (db *appdbimpl) DeleteComment(dbComment DatabaseComment) error {
	// remove the comment from the database
	res, err := db.c.Exec(`
		DELETE FROM Comment
		WHERE id=?
	`, dbComment.Id)

	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()

	// if there are no affected rows
	// then the photo was not commented
	if aff == 0 {
		return ErrPhotoNotCommented
	}

	return err
}

func (db *appdbimpl) GetCommentList(dbPhoto DatabasePhoto, dbUser DatabaseUser) (DatabaseCommentList, error) {
	dbCommentList := DatabaseCommentListDefault()

	// get the table of the comments under the photo
	// without considering the comments made by users
	// who banned the user performing the action
	rows, err := db.c.Query(`
		SELECT id, user, photo, date, comment_body
		FROM Comment
		WHERE photo=?
		AND user NOT IN (
			SELECT first_user
			FROM ban
			WHERE second_user=?
		)
		ORDER BY date
	`, dbPhoto.Id, dbUser.Id)

	if errors.Is(err, sql.ErrNoRows) {
		return dbCommentList, ErrPhotoDoesNotExist
	}

	if err != nil {
		return dbCommentList, err
	}

	dbCommentPhoto := DatabasePhotoDefault()

	// build the comment list
	for rows.Next() {
		dbComment := DatabaseCommentDefault()

		err = rows.Scan(&dbComment.Id, &dbComment.User.Id, &dbComment.Photo.Id, &dbComment.Date, &dbComment.CommentBody)

		if err != nil {
			return dbCommentList, err
		}

		dbCommentUser, err := db.GetDatabaseUser(dbComment.User.Id)

		if err != nil {
			return dbCommentList, err
		}

		dbComment.User = dbCommentUser

		if dbCommentPhoto.Id == 0 {
			dbCommentPhoto, err = db.GetDatabasePhoto(dbComment.Photo.Id, dbUser)

			if err != nil {
				return dbCommentList, err
			}
		}

		dbComment.Photo = dbCommentPhoto

		dbCommentList.Comments = append(dbCommentList.Comments, dbComment)
	}

	if rows.Err() != nil {
		return dbCommentList, err
	}

	_ = rows.Close()

	return dbCommentList, err
}
