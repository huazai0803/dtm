package storage

import (
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/yedf/dtm/dtmcli/dtmimp"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("storage: NotFound")
var ErrUniqueConflict = errors.New("storage: UniqueKeyConflict")

type Store interface {
	Ping() error
	PopulateData(skipDrop bool)
	FindTransGlobalStore(gid string) *TransGlobalStore
	ScanTransGlobalStores(position *string, limit int64) []TransGlobalStore
	FindBranches(gid string) []TransBranchStore
	UpdateBranchesSql(branches []TransBranchStore, updates []string) *gorm.DB
	LockGlobalSaveBranches(gid string, status string, branches []TransBranchStore, branchStart int)
	MaySaveNewTrans(global *TransGlobalStore, branches []TransBranchStore) error
	ChangeGlobalStatus(global *TransGlobalStore, newStatus string, updates []string, finished bool)
	TouchCronTime(global *TransGlobalStore, nextCronInterval int64)
	LockOneGlobalTrans(expireIn time.Duration) *TransGlobalStore
}

func wrapError(err error) error {
	if err == gorm.ErrRecordNotFound || err == redis.Nil {
		return ErrNotFound
	}
	dtmimp.E2P(err)
	return err
}
