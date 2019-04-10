package framework

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

/*会话*/
type Session struct {
	mSessionID        string                      //唯一id
	mLastTimeAccessed time.Time                   //最后访问时间
	mValues           map[interface{}]interface{} //其它对应值(保存用户所对应的一些值，比如用户权限之类)
}

type SessionMgr struct {
	mCookieName string
	mLock sync.RWMutex
	mMaxLifeTime int64

	mSessions map[string]*Session
}

func NewSessionMgr(cookieName string, maxLeftTime int64) *SessionMgr {
	mgr := &SessionMgr{
		mCookieName:cookieName,
		mLock:sync.RWMutex{},
		mMaxLifeTime:maxLeftTime,
		mSessions: make(map[string]*Session),
	}

	go mgr.Gc()

	return mgr
}

func (mgr *SessionMgr) StartSession(w http.ResponseWriter, r *http.Request) string {

	newSessionID := url.QueryEscape(NewSessionID())

	session := &Session{
		mSessionID:newSessionID,
		mLastTimeAccessed:time.Now(),
		mValues: map[interface{}]interface{}{},
	}

	mgr.mSessions[newSessionID] = session

	cookie := &http.Cookie{
		Name:mgr.mCookieName,
		Value:newSessionID,
		HttpOnly:true,
		Path:"/",
		MaxAge:int(mgr.mMaxLifeTime),
	}

	http.SetCookie(w, cookie)

	return newSessionID
}

//创建唯一ID
func NewSessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		nano := time.Now().UnixNano() //微秒
		return strconv.FormatInt(nano, 10)
	}

	return base64.URLEncoding.EncodeToString(b)
}

func (mgr *SessionMgr) Gc() {
	mgr.mLock.Lock()
	defer mgr.mLock.Unlock()

	for sessionId, session := range mgr.mSessions {
		if session.mLastTimeAccessed.Unix() + mgr.mMaxLifeTime < time.Now().Unix() {
			delete(mgr.mSessions, sessionId)
		}
	}

	time.AfterFunc(time.Duration(mgr.mMaxLifeTime) * time.Second, mgr.Gc)
}


