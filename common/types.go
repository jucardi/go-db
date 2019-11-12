package common

import (
	"time"
)

type MigrationInfo struct {
	//IdInt     uint      `gorm:"primary_key"`
	//IdString  string    `bson:"_id"`
	ScriptId  string    `bson:"script_id"`
	Hash      string    `bson:"hash"`
	Timestamp time.Time `bson:"timestamp"`
}
