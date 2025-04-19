package datab

import (
	"database/sql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Context struct {
	Db     *sql.DB
	User   UserRepo
	Media  MediaRepo
	Bucket BucketObjectRepo
}

func New(databasePath string) Context {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalln(err)
	}

	context := Context{
		Db: db,
	}

	context.User = GetUserRepository(&context)
	context.Media = GetMediaRepository(&context)
	context.Bucket = GetBucketObjectRepository(&context)

	return context
}
