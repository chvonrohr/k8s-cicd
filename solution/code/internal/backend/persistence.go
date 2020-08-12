package backend

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/model"
)

// InitialisePersistence creates a new database connection using config variables.
//
// The available variables are:
// db.username - username
// db.password - password
// db.host - host
// db.database - database
// db.port - port (numeric)
func InitialisePersistence() (*gorm.DB, error) {
	var (
		username = viper.GetString("db.username")
		password = viper.GetString("db.password")
		host     = viper.GetString("db.host")
		port     = viper.GetInt("db.port")
		database = viper.GetString("db.database")
	)
	// format dsn based on above values
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", host, port, username, database, password)
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.Site{})
	db.AutoMigrate(&model.Crawl{})
	db.AutoMigrate(&model.Page{})

	return db, nil
}

func PersistenceMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, t := NewTransaction(db)
		defer t.Close()
		ctx.Set("tx", t)
		ctx.Next()
	}
}

func GetTx(ctx *gin.Context) *gorm.DB {
	return ctx.MustGet("tx").(*Transaction).Tx
}

func FailTx(ctx *gin.Context) {
	t := ctx.MustGet("tx").(*Transaction)
	t.Fail()
}

type Transaction struct {
	once     sync.Once
	rollback bool
	Tx       *gorm.DB
}

func (t *Transaction) Close() {
	t.once.Do(func() {
		if t.rollback {
			t.Tx.Rollback()
		} else {
			t.Tx.Commit()
		}
	})
}

func (t *Transaction) Fail() {
	t.rollback = true
}

func NewTransaction(db *gorm.DB) (*gorm.DB, *Transaction) {
	db = db.Begin()
	return db, &Transaction{Tx: db}
}
