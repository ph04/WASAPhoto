/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// Ban
	InsertBan(dbUser DatabaseUser, bannedDbUser DatabaseUser) error             // DONE
	DeleteBan(dbUser DatabaseUser, bannedDbUser DatabaseUser) error             // DONE
	CheckBan(firstDbUser DatabaseUser, secondDbUser DatabaseUser) (bool, error) // DONE

	// Follow
	InsertFollow(dbUser DatabaseUser, followedDbUser DatabaseUser) error                          // DONE
	DeleteFollow(dbUser DatabaseUser, followedDbUser DatabaseUser) error                          // DONE
	GetFollowersCount(profileDbUser DatabaseUser, dbUser DatabaseUser) (int, error)               // DONE
	GetFollowingCount(profileDbUser DatabaseUser, dbUser DatabaseUser) (int, error)               // DONE
	GetFollowersList(followersDbUser DatabaseUser, dbUser DatabaseUser) (DatabaseUserList, error) // DONE
	GetFollowingList(followingDbUser DatabaseUser, dbUser DatabaseUser) (DatabaseUserList, error) // DONE
	GetFollowStatus(firstDbUser DatabaseUser, secondDbUser DatabaseUser) (bool, error)            // DONE

	// Photo
	GetDatabasePhoto(photoId uint32, dbUser DatabaseUser) (DatabasePhoto, error) // DONE
	InsertPhoto(dbPhoto *DatabasePhoto) error                                    // DONE
	DeletePhoto(dbPhoto DatabasePhoto) error                                     // DONE
	GetPhotoLikeCount(dbPhoto *DatabasePhoto, dbUser DatabaseUser) error         // DONE
	GetPhotoCommentCount(dbPhoto *DatabasePhoto, dbUser DatabaseUser) error      // DONE
	GetPhotoLikeStatus(dbPhoto *DatabasePhoto, dbUser DatabaseUser) error        // DONE
	GetPhotos(dbProfile *DatabaseProfile, dbUser DatabaseUser) error             // DONE
	GetPhotoCount(dbUser DatabaseUser) (int, error)                              // DONE

	// Like
	InsertLike(dbUser DatabaseUser, dbPhoto DatabasePhoto) error                      // DONE
	DeleteLike(dbUser DatabaseUser, dbPhoto DatabasePhoto) error                      // DONE
	GetLikeList(dbPhoto DatabasePhoto, dbUser DatabaseUser) (DatabaseUserList, error) // DONE

	// Comment
	GetDatabaseComment(commentId uint32, dbUser DatabaseUser) (DatabaseComment, error)      // DONE
	InsertComment(dbComment *DatabaseComment) error                                         // DONE
	DeleteComment(dbComment DatabaseComment) error                                          // DONE
	GetCommentList(dbPhoto DatabasePhoto, dbUser DatabaseUser) (DatabaseCommentList, error) // DONE

	// Stream
	GetDatabaseStream(dbUser DatabaseUser) (DatabaseStream, error) // DONE

	// User
	GetDatabaseUser(userId uint32) (DatabaseUser, error)                              // DONE
	GetDatabaseUserFromDatabaseLogin(dbLogin DatabaseLogin) (DatabaseUser, error)     // DONE
	InsertUser(dbUser *DatabaseUser) error                                            // DONE
	UpdateUser(oldDbUser DatabaseUser, newDbUser DatabaseUser) error                  // DONE
	GetUserList(dbUser DatabaseUser, dbLogin DatabaseLogin) (DatabaseUserList, error) // DONE

	// Liveness
	Ping() error // DONE
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	var err error

	// enable checks for foreign keys
	_, err = db.Exec("PRAGMA foreign_key=ON")

	if err != nil {
		return nil, err
	}

	userTable := `
		CREATE TABLE IF NOT EXISTS User (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE
		);
	`
	photoTable := `
		CREATE TABLE IF NOT EXISTS Photo (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			user INTEGER NOT NULL,
			url TEXT NOT NULL,
			date TEXT NOT NULL,
			FOREIGN KEY (user) REFERENCES User(name)
		);
	`
	commentTable := `
		CREATE TABLE IF NOT EXISTS Comment (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			user INTEGER NOT NULL,
			photo INTEGER NOT NULL,
			date TEXT NOT NULL,
			comment_body TEXT NOT NULL,
			FOREIGN KEY (user) REFERENCES User(name),
			FOREIGN KEY (photo) REFERENCES Photo(id)
		);
	`
	followTable := `
		CREATE TABLE IF NOT EXISTS follow (
			first_user INTEGER NOT NULL,
			second_user INTEGER NOT NULL,
			PRIMARY KEY (first_user, second_user),
			FOREIGN KEY (first_user) REFERENCES User(name),
			FOREIGN KEY (second_user) REFERENCES User(name)
		);
	`
	banTable := `
		CREATE TABLE IF NOT EXISTS ban (
			first_user INTEGER NOT NULL,
			second_user INTEGER NOT NULL,
			PRIMARY KEY (first_user, second_user),
			FOREIGN KEY (first_user) REFERENCES User(name),
			FOREIGN KEY (second_user) REFERENCES User(name)
		);
	`
	likeTable := `
		CREATE TABLE IF NOT EXISTS like (
			user INTEGER NOT NULL,
			photo INTEGER NOT NULL,
			PRIMARY KEY (user, photo),
			FOREIGN KEY (user) REFERENCES User(name),
			FOREIGN KEY (photo) REFERENCES Photo(id)
		);
	`

	_, err = db.Exec(userTable)

	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(photoTable)

	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(commentTable)

	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(followTable)

	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(banTable)

	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(likeTable)

	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
