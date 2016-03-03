package cmd

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/go-ini/ini"
	"github.com/luksen/maildir"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
	"strings"
)

type CachedMaildir struct {
	dir         *maildir.Dir
	flagsBucket string
}

type MailAccount struct {
	name  string
	boxes map[string]CachedMaildir
}

func NewCachedMailDir(mdirPath string, db *bolt.DB) (dir CachedMaildir) {
	mdir := maildir.Dir(mdirPath)
	bucketName := path.Base(mdirPath)

	tx, err := db.Begin(true)
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.CreateBucketIfNotExists([]byte(bucketName + "_flags"))
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()

	dir = CachedMaildir{dir: &mdir, flagsBucket: bucketName + "_flags"}
	SynchronizeInfo(dir, db)
	return
}

// Make sure that mail info in db and maildir are identical
// TODO Synchronize messages removed in maildir to db
func SynchronizeInfo(cdir CachedMaildir, db *bolt.DB) {
	cdir.dir.Unseen()
	keys, err := cdir.dir.Keys()
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin(true)
	if err != nil {
		log.Fatal(err)
	}
	bucket := tx.Bucket([]byte(cdir.flagsBucket))

	for _, key := range keys {
		ret := bucket.Get([]byte(key))
		if ret == nil {
			// We don't know about this message
			flags, err := cdir.dir.Flags(key)
			if err != nil {
				println(err)
			}
			bucket.Put([]byte(key), []byte(flags))
		}
	}
	tx.Commit()
}

func GetMaildirPath(key string) string {
	home := os.Getenv("HOME")
	cfg, err := ini.Load(home + "/.msgconfig")
	if err != nil {
		log.Fatal(err)
	}
	section, err := cfg.GetSection("mail")
	if err != nil {
		log.Fatal(err)
	}
	maildir, err := section.GetKey(key)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(maildir.String(), "~", home, -1)
}

func Status() {
	db, err := bolt.Open("msg.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	path := GetMaildirPath("drafts")
	dir := NewCachedMailDir(path, db)

	var drafts []string

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dir.flagsBucket))
		return b.ForEach(func(k, v []byte) error {
			drafts = append(drafts, string(k))
			return nil
		})
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(drafts) > 0 {
		println("## drafts")
		for i, draft := range drafts {
			val, _ := dir.dir.Header(draft)
			fmt.Printf("[%d] To: %s // Subject: %s\n", i+1, val.Get("To"), val.Get("Subject"))
		}
	} else {
		println("No drafts")
	}
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show draft status",
	Long:  `Show draft status`,
	Run: func(cmd *cobra.Command, args []string) {
		Status()
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
