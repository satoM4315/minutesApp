package main

import (
	"encoding/base64"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	_ "github.com/mattn/go-sqlite3" //DBのパッケージだが、操作はGORMで行うため、importだけして使わない
)

// セッション情報
type TempSession struct {
	gorm.Model
	SessionID string `gorm:"unique;not null"`
	UserID    string `gorm:"not null"`
	ValidTime time.Time
}

// 指定したsessionIDのセッションがあるか確認する
func SessionExist(sessionID string) bool {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(dbGetOne)")
	}
	var session TempSession
	var count int

	db.Where(&TempSession{SessionID: sessionID}).Find(&session).Count(&count)
	if count == 0 {
		return false
	}
	db.Close()
	return true
}

// 指定したuserIDのセッションがあるか確認する
func SessionExistByUserID(userID string) bool {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(dbGetOne)")
	}
	var session TempSession
	var count int
	db.Where(&TempSession{UserID: userID}).Find(&session).Count(&count)
	if count == 0 {
		return false
	}
	db.Close()
	return true
}

//指定したsessionIDのオブジェクトが削除される
func sessionDelete(sessionID string) {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(dbDelete)")
	}
	var session TempSession
	db.Where(&TempSession{SessionID: sessionID}).Limit(1).Find(&session)
	db.Delete(&session)
	db.Close()
}

// sessionを作成。sessionIDとuserIDの組みを格納し、sessionIDを返す
func createSession(userID string) string {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(createUser)")
	}
	defer db.Close()

	sessionID := LongSecureRandomBase64()
	now := time.Now()

	// Insert処理
	if err := db.Create(&TempSession{SessionID: sessionID, UserID: userID, ValidTime: now.Add(5 * time.Hour)}).Error; err != nil {

		return ""
	}
	return sessionID

}

// 指定したsessionIDのuserIDを返す
func getuserIDBySessionID(sessionID string) string {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(getUser)")
	}
	defer db.Close()
	var session TempSession
	db.Where(&TempSession{SessionID: sessionID}).Find(&session)

	return session.UserID
}

// getUserById は、指定されたIDを持つユーザーを一つ返します。
// ユーザーが存在しない場合、空のレコードが返る?(GORMの仕様を要確認)
func getSessionIDByuserID(userID string) string {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(getUserById)")
	}
	defer db.Close()
	var session TempSession
	db.Where(&TempSession{UserID: userID}).Limit(1).Find(&session)
	return session.SessionID
}

//期限切れのセッション情報を削除
func sessionStoreUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := gorm.Open("sqlite3", "minutes.sqlite3")
		if err != nil {
			panic("データベース開ません(getUserById)")
		}

		var session TempSession
		now := time.Now()
		term_date := now.Add(-10 * 24 * time.Hour) //10日間アクセスのないsessionIDは自動で消去される

		db.Where("valid_time <= ?", term_date).Delete(&session)
		db.Close()
		c.Next()
	}
}

//sessionの有効期限をすぎていないか
//すぎていたらtureを返す
func SessionTimeOut(sessionID string) bool {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(dbGetOne)")
	}
	var session TempSession
	now := time.Now()

	db.Where(&TempSession{SessionID: sessionID}).Find(&session)

	if session.ValidTime.Before(now) {
		return true
	}
	db.Close()
	return false
}

//正当なセッションIDを持っているか確認する
//ただし、権限などはここでは保証しない
func sessionCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		session := sessions.Default(ctx)
		sessionID := session.Get("SessionID")

		// セッションがない場合
		if sessionID == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			ctx.Abort()
			return
		} else if !(SessionExist(sessionID.(string))) {
			session.Clear()
			session.Save()
			// 不当なセッション情報によるアクセスの場合
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
			ctx.Abort()
			return
		} else if SessionTimeOut(sessionID.(string)) {
			//セッション有効時間が切れていた場合
			//セッションからデータを破棄する
			sessionDelete(sessionID.(string))
			//session.Clear()
			//session.Save()

			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Session time out"})
			ctx.Abort()
			return
		}
		//ctx.Next()

	}
}

//session IDを生成するための関数群
func SecureRandom() string {
	return uuid.New().String()
}

func SecureRandomBase64() string {
	return base64.StdEncoding.EncodeToString(uuid.New().NodeID())
}

func LongSecureRandomBase64() string {
	return SecureRandomBase64() + SecureRandomBase64()
}

func MultipleSecureRandomBase64(n int) string {
	if n <= 1 {
		return SecureRandomBase64()
	}
	return SecureRandomBase64() + MultipleSecureRandomBase64(n-1)
}

//cookieが正当なものか、セッションIDが正しいものか確認する
//正しいものでなければentranceに追い返す
//cookieからuserIDを返す
func getSessionUserID(ctx *gin.Context) User {
	session := sessions.Default(ctx)
	sessionID := session.Get("SessionID").(string)
	if !(SessionExist(sessionID)) {
		ctx.Redirect(http.StatusSeeOther, "/entrance")
		ctx.Abort()
		return User{}
	}
	userID := getuserIDBySessionID(sessionID)
	return getUser(userID)

}

//セッション全取得、デバッグ用
func sessionGetAll() []TempSession {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(dbGetAll)")
	}
	var sessions []TempSession
	db.Order("created_at desc").Find(&sessions)
	db.Close()
	return sessions
}

//指定したuserIDのレコードのvalid_timeを今にする(今が有効期限になる)
//時間がたったセッション情報がtimeoutになるかテストするとき用
func sessionTimeSetNow(userID string) {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(dbGetAll)")
	}
	var session TempSession
	db.Where(&TempSession{UserID: userID}).Find(&session)

	now := time.Now()
	session.ValidTime = now
	db.Save(&session)
	db.Close()
}

//指定したuserIDのレコードのvalid_timeを今にする(今が有効期限になる)
//時間がたったセッション情報がtimeoutになるかテストするとき用
func sessionTimeSet10daysLater(userID string) {
	db, err := gorm.Open("sqlite3", "minutes.sqlite3")
	if err != nil {
		panic("データベース開ません(dbGetAll)")
	}
	var session TempSession
	db.Where(&TempSession{UserID: userID}).Find(&session)

	now := time.Now()
	session.ValidTime = now.Add(-10 * 24 * time.Hour)
	db.Save(&session)
	db.Close()
}
