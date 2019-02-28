package storage

import (
	"database/sql"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/go-testfixtures/testfixtures"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"minibox.ai/pkg/utils"
)

var (
	db       *sql.DB
	fixtures *testfixtures.Context
)

type User struct {
	ID        int
	Name      string
	Namespace string

	JoinedProjects []*Project `gorm:"many2many:user_projects;"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Project struct {
	// gorm.JoinTableHandler
	ID        int
	Name      string
	Namespace string
	Members   []*User `gorm:"many2many:user_projects;"`
	Author    User
	AuthorID  int
}

type UserProject struct {
	UserID    int
	ProjectID int
	Open      bool
}

func TestMain(m *testing.M) {
	var err error
	gdb, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Open sqlite3 failed: %v", err)
	}
	gdb.LogMode(true)

	gdb.AutoMigrate(&Project{})
	gdb.AutoMigrate(&User{})
	gdb.AutoMigrate(&UserProject{})

	gdb.Close()
	// Open connection with the test database.
	// Do NOT import fixtures in a production database!
	// Existing data would be deleted
	db, err = sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	// creating the context that hold the fixtures
	// see about all compatible databases in this page below
	fixtures, err = testfixtures.NewFolder(db, &testfixtures.SQLite{}, "fixtures/")

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func prepareTestDatabase() *gorm.DB {
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}

	gdb, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Open sqlite3 failed: %v", err)
	}
	// gdb.SetJoinTableHandler(&User{}, "Projects", &Project{})

	return gdb
}

func TestLoad(t *testing.T) {
	db := prepareTestDatabase()
	// s := NewStorage(db)
	defer db.Close()
	db.LogMode(true)

	var (
		usr = User{ID: 1}
		// prjs = []testProject{}
		prj = Project{AuthorID: 1}
	)

	db.Preload("projects").Find(&usr)
	db.Find(&prj, "id = ?", 1)
	// if err := db.Preload("test_users").First(&prj, "id = ?", 1).Error; err != nil {
	// 	t.Fatalf("reated error %s", err)
	// }
	// t.Logf("")
}

func pluckIds(vals interface{}, field string) (ids []int) {
	v := reflect.ValueOf(vals)
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Slice {
		panic("must slice")
	}

	ids = make([]int, v.Len())
	for i := 0; i < v.Len(); i++ {
		ele := v.Index(i)
		ev := ele.FieldByName(field)
		ids[i] = int(ev.Int())
	}

	return
}

func TestJointable(t *testing.T) {
	db := prepareTestDatabase()
	// s := NewStorage(db)
	defer db.Close()
	db.LogMode(true)

	var (
		usr    = []User{User{ID: 1}, User{ID: 2}}
		prjs   = []Project{}
		proids = []int{}
	)

	// db.First(&prjs, "id = ?", 1)

	if err := db.Model(&usr).Where("user_projects.open = ?", true).Order("name desc").Limit(10).Related(&prjs, "JoinedProjects").Error; err != nil {
		t.Fatalf("%s", err)
	}

	t.Logf("User Projects: %s\n", utils.Prettify(prjs))
	proids = pluckIds(&prjs, "ID")

	// if err := db.Model(&prjs).Related(&usr, "Members").Error; err != nil {
	// 	t.Fatalf("%s", err)
	// }

	db.Preload("Members").Find(&prjs, proids)
	t.Logf("Members: %s\n", utils.Prettify(prjs[0].Members))

	// db.Find(&prj, "id = ?", 1)
}

func TestMax(t *testing.T) {
	db := prepareTestDatabase()
	// s := NewStorage(db)
	defer db.Close()
	db.LogMode(true)

	var (
		usr User
	)

	// db.First(&prjs, "id = ?", 1)

	if err := db.Select("max(id) as id, *").First(&usr).Error; err != nil {
		t.Fatalf("%s", err)
	}

	t.Logf("usr max %#v", usr)
}

func TestLevel(t *testing.T) {
	DB, err := gorm.Open("sqlite3", "test.db")
	DB.LogMode(true)

	if err != nil {
		log.Fatalf("Open sqlite3 failed: %v", err)
	}
	type (
		Level1 struct {
			ID       uint
			Value    string
			Level2ID uint
		}
		Level2 struct {
			ID       uint
			Level1s  []Level1
			Level3ID uint
		}
		Level3 struct {
			ID     uint
			Name   string
			Level2 Level2
		}
	)
	DB.DropTableIfExists(&Level3{})
	DB.DropTableIfExists(&Level2{})
	DB.DropTableIfExists(&Level1{})
	if err := DB.AutoMigrate(&Level3{}, &Level2{}, &Level1{}).Error; err != nil {
		t.Error(err)
	}

	want := make([]Level3, 2)
	want[0] = Level3{
		Level2: Level2{
			Level1s: []Level1{
				{Value: "value1"},
				{Value: "value2"},
			},
		},
	}
	if err := DB.Create(&want[0]).Error; err != nil {
		t.Error(err)
	}
	want[1] = Level3{
		Level2: Level2{
			Level1s: []Level1{
				{Value: "value3"},
				{Value: "value4"},
			},
		},
	}
	if err := DB.Create(&want[1]).Error; err != nil {
		t.Error(err)
	}

	var got []Level3
	if err := DB.Preload("Level2.Level1s").Find(&got).Error; err != nil {
		t.Error(err)
	}
	log.Printf("got %s\n", utils.Prettify(&got))
}
