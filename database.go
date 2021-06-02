package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" //DBのパッケージだが、操作はGORMで行うため、importだけして使わない
	"golang.org/x/crypto/bcrypt"    // パスワードを暗号化する際に使う
)

/*
gorm.Modelの中身
カラム
  ・id
  ・created_at
  ・updated_at
  ・deleted_at
*/
/*外部からカラムを参照するときは
id → ID
created_at → CreatedAt
updated_at → UpdatedAt
deleted_at → DeletedAt
*/
// テーブル名：messages -->　テーブル名は自動で複数形になる
type Message struct {
	gorm.Model
	Message   string
	MeetingID uint `gorm:"column:meeting_id";`
	UserID    uint
}

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

type Meeting struct {
	gorm.Model
	Name   string `gorm:"unique;not null"`
	UserID uint
}

type Entry struct {
	gorm.Model
	MeetingID int `gorm:"unique;not null;column:meeting_id;"`
	UserID    uint
}

/*
DBの内容
(ID,作成日,更新日,削除日のカラムは全てに入っている)
・ユーザー
  ・ユーザーネーム（ログインID）
  ・パスワード（暗号化したもの）
・会議
  ・会議名
・メッセージ
  ・内容
  ・会議ID
  ・ユーザーID
・エントリー
  ・会議ID
  ・ユーザーID
*/

//DBマイグレート
//main関数の最初でdbInit()を呼ぶことでデータベースマイグレート
func dbInit() {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3") //第一引数：使用するDBのデバイス。第二引数：ファイル名
	if err != nil {
		panic("データベース開ません(dbinit)")
	}
	db.AutoMigrate(&User{}, &Message{}, Meeting{}, &Entry{}, &TempSession{}) //ファイルがなければ、生成を行う。すでにあればマイグレート。すでにあってマイグレートされていれば何も行わない
	defer db.Close()
}

//DB追加
//ミーティングIDを指定するように変更
func messageInsert(message string, meetingID uint, userID uint) {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(dbInsert)")
	}
	db.Create(&Message{
		Message:   message,
		MeetingID: meetingID,
		UserID:    userID,
	})
	defer db.Close()
}

//指定されたミーティングIDを持つメッセージを全取得
func MeetingMessageGetAll(meetingID uint) []Message {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(dbGetAll)")
	}
	var messages []Message
	//カラム名は自動でスネークケース
	//MeetingID => meeting_id
	db.Order("created_at desc").Find(&messages, "meeting_id = ?", meetingID) //db.Find(&messages)で構造体Messageに対するテーブルの要素全てを取得し、それをOrder("created_at desc")で新しいものが上に来るように並び替えている
	db.Close()
	return messages
}

//DB一つ取得
//idを与えることで、該当するMessageオブジェクトが一つ返される
func dbGetOne(id int) Message {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(dbGetOne)")
	}
	var message Message
	db.First(&message, id)
	db.Close()
	return message
}

//DB更新
//idとmessageを与えることで、該当するidのMessageオブジェクトのMessageが更新される
func dbUpdate(id int, update_message string) {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(dgUpdate)")
	}
	var message Message
	db.First(&message, id)
	message.Message = update_message
	db.Save(&message)
	db.Close()
}

//DB削除
//指定したidのMessageオブジェクトが削除される
func dbDelete(id int) {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(dbDelete)")
	}
	var message Message
	db.First(&message, id)
	db.Delete(&message)
	db.Close()
}

// ユーザー登録処理
func createUser(username string, password string) error {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(createUser)")
	}
	defer db.Close()
	// パスワード暗号化
	passwordEncrypt, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// Insert処理
	if err := db.Create(&User{Username: username, Password: string(passwordEncrypt)}).Error; err != nil {
		return err
	}
	return nil
}

// ユーザーネーム(ログインID)を指定してそのユーザーのレコードを取ってくる
// 指定したユーザーのレコードがない場合は、IDが0のレコードを返す
func getUser(username string) User {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(getUser)")
	}
	defer db.Close()
	var user User
	db.Where(&User{Username: username}).Find(&user)
	return user
}

// getUserById は、指定されたIDを持つユーザーを一つ返します。
// ユーザーが存在しない場合、IDが0のレコードが返ります。
func getUserByID(id uint) User {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(getUserById)")
	}
	defer db.Close()
	var user User
	db.First(&user, id)
	return user
}

// パスワードの比較
// dbPasswordはデータベースから取ってきたパスワード（暗号化済み）
// formPasswordはログイン時に入力されたパスワード（平文）
func comparePassword(dbPassword string, formPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(formPassword))
}

//議事録一覧を取得
func getAllMeeting() []Meeting {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(getAllMeeting)")
	}
	defer db.Close()
	meeting := make([]Meeting, 0)
	db.Order("created_at desc").Find(&meeting)
	return meeting
}

//主に議事録名の取得に用いる
func getMeetingByID(meetingID uint) Meeting {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(getUserById)")
	}
	defer db.Close()
	var meeting Meeting
	db.Find(&meeting, "id = ?", meetingID)
	return meeting
}

//議事録追加
func createMeeting(meeting string, userID uint) error {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(dbInsert)")
	}
	if err := db.Create(&Meeting{Name: meeting, UserID: userID}).Error; err != nil {
		return err
	}

	defer db.Close()
	return nil
}
