package model

import (
	"fmt"
	"github.com/golang/glog"
	"time"
	"github.com/dolotech/lib/db"
	"github.com/dolotech/lib/utils"
	"math/rand"
	"errors"
)

func (this *User) GetById() (bool, error) {
	return db.C().Engine().Where("uid = ?", this.Uid).Get(this)
}

func (this *User) GetByAccount() (bool, error) {
	return db.C().Engine().Where("account = ?", this.Account).Get(this)
}

func (this *User) GetByUnionId() (bool, error) {
	return db.C().Engine().Where("union_id = ?", this.UnionId).Get(this)
}

type User struct {
	Uid        uint32    `xorm:"'uid' pk autoincr BIGINT"`            // 用户id
	Account    string    `xorm:"'account' index unique  VARCHAR(16)"` // 客户端玩家展示的账号
	DeviceId   string    `xorm:"'device_id' VARCHAR(32)"`             // 设备id
	UnionId    string    `xorm:"'union_id' VARCHAR(32)"`              // 微信联合id
	Nickname   string    `xorm:"'nickname' VARCHAR(32)"`              // 微信昵称
	Sex        uint8     `xorm:"'sex' smallint"`                      // 微信性别 0-未知，1-男，2-女
	Profile    string    `xorm:"'profile' VARCHAR(64)"`               // 微信头像
	Invitecode string    `xorm:"'invitecode' VARCHAR(6)"`             // 绑定的邀请码
	Chips      uint32    `xorm:"'chips'"`                             // 筹码
	Lv         uint8     `xorm:"'lv' smallint"`                       // 等级
	CreatedAt  time.Time `xorm:"'created_at' index  created"`         // 注册时间
	LastTime   time.Time `xorm:"'last_time'"`                         // 上次登录时间
	LastIp     uint32    `xorm:"'last_ip' BIGINT"`                    // 最后登录ip
	Kind       uint8     `xorm:"'kind'  not null smallint"`           // 用户类型
	Disable    bool      `xorm:"'disable'"`                           // 是否禁用
	Signature  string    `xorm:"'signature' VARCHAR(64)"`             // 个性签名
	Gps        string    `xorm:"'gps' VARCHAR(32)"`                   // gps定位数据
	Black      bool      `xorm:"'black'"`                             // 黑名单列表
	RoomID     uint32    `xorm:"'room_id'"`                           // 当前所在房间号，0表示不在房间,用于掉线重连
}

func (u *User) Insert() error {
	_, err := db.C().Engine().InsertOne(u)
	if err != nil {
		glog.Errorln(err)
		return err
	}

return  nil
	//INSERT INTO COMPANY (ID,NAME,AGE,ADDRESS,SALARY,JOIN_DATE) VALUES (1, 'Paul', 32, 'California', 20000.00 ,'2001-07-13');
	//sql := fmt.Sprintf(`INSERT INTO public.user(union_id,nickname,sex,profile)VALUES(%v,%v,%v,%v)`, u.UnionId, u.Nickname, u.Sex, u.Profile)
	//sql:= fmt.Sprintf("insert into public.user(union_id,nickname,sex,profile)values(%v,%v,%v,%v)",u.UnionId,u.Nickname,u.Sex,u.Profile)
	if len(u.UnionId ) == 0{
			return errors.New("unionid can not empty")
	}

	if len(u.Nickname ) == 0{
		u.Nickname= " "
	}

	if len(u.Profile ) == 0{
		u.Profile= " "
	}

	if u.Sex == 0{
		u.Sex = 1
	}

	//sql:= fmt.Sprintf(`INSERT INTO "user"("union_id","nickname","sex","profile") VALUES ("%v","%v","%v","%v")`,u.UnionId,u.Nickname,u.Sex,u.Profile)
	//sql:=`INSERT INTO public.user(union_id,nickname,sex,profile) VALUES ('abcdef','',0,'')`
	//sql:=`INSERT INTO public.user(union_id,nickname,sex,profile) VALUES ('fabcdef','',0,'') return uid`
	sql:= fmt.Sprintf(`INSERT INTO public.user(union_id,nickname,sex,profile) VALUES ('%v','%v','%v','%v') RETURNING uid`,"ffabd",u.Nickname,u.Sex,u.Profile)
	rel, err := db.C().Engine().Exec(sql)
	if err != nil {
		glog.Errorln(err)
		return err
	}
	glog.Errorln(rel)
	id, err := rel.LastInsertId()
	glog.Errorln(id)
	if err != nil {
		glog.Errorln(err)
		return err
	}

	n := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(9999-1000) + 1000

	account := fmt.Sprintf("%v%v", id, n)

	//"select currval('address_address_id_seq')"
	//

	sql = fmt.Sprintf(`UPDATE public.user SET
	account =%s WHERE uid = %d `, account, id)
	_, err = db.C().Engine().Exec(sql)
	if err != nil {
		glog.Errorln(err)
		return err
	}

	return nil
}
func (u *User) UpdateChips(value int32) error {
	sql := fmt.Sprintf(`UPDATE public.user SET
	chips = chips + %d WHERE uid = %d `, value, u.Uid)
	_, err := db.C().Engine().Exec(sql)
	if err != nil {
		glog.Errorln(err)
	}
	return err
}

func (u *User) UpdateLogin(ip string) error {
	sql := fmt.Sprintf(`UPDATE public.user SET
	last_time =  %d ,last_ip =  %d WHERE uid = %d `, time.Now(), utils.InetToaton(ip), u.Uid)
	_, err := db.C().Engine().Exec(sql)
	if err != nil {
		glog.Errorln(err)
	}
	return err
}

func (u *User) UpdateRoomId() error {
	sql := fmt.Sprintf(`UPDATE public.user SET
	room_id = chips + %d WHERE uid = %d `, u.RoomID, u.Uid)
	_, err := db.C().Engine().Exec(sql)
	if err != nil {
		glog.Errorln(err)
	}
	return err
}

/*

func (data *UserData) initValue() error {
	userID, err := mongoDBNextSeq(USERDB)
	if err != nil {
		return fmt.Errorf("get next users id error: %v", err)
	}

	data.UserID = userID
	data.AccountID = time.Now().Format("0102") + strconv.Itoa(int(data.UserID))
	data.CreatedAt = uint32(time.Now().Unix())
	return nil
}

func (data *UserData) GetByWechat() error {
	db := mongoDB.Ref()
	defer mongoDB.UnRef(db)
	err := db.DB(DB).C(USERDB).Find(bson.M{"unionid": data.UnionId}).One(data)
	if err != nil {
		glog.Errorln(err)
		return err
	}
	return nil
}

func (data *UserData) IncChips(change int) error {
	db := mongoDB.Ref()
	defer mongoDB.UnRef(db)
	err := db.DB(DB).C(USERDB).UpdateId(data.UserID, bson.M{"$inc": bson.M{"chips": change}})
	if err != nil {
		glog.Errorln(err)
		return err
	}
	return nil
}

func (data *UserData) UpdateSex() error {
	db := mongoDB.Ref()
	defer mongoDB.UnRef(db)
	err := db.DB(DB).C(USERDB).UpdateId(data.UserID, bson.M{"$set": bson.M{"sex": data.Sex}})
	if err != nil {
		glog.Errorln(err)
		return err
	}
	return nil
}

func (data *UserData) Insert() error { //注册
	db := mongoDB.Ref()
	defer mongoDB.UnRef(db)
	err := db.DB(DB).C(USERDB).Insert(data)
	if err != nil {
		glog.Errorln(err)
		return err
	}
	return nil
}

func (data *UserData) Register() error { //注册
	db := mongoDB.Ref()
	defer mongoDB.UnRef(db)

	err := data.Insert()
	if err != nil {
		glog.Errorln(err)
		return err
	}

	return nil
}

func (data *UserData) Login(user *msg.UserLoginInfo) error {
	var result UserData
	db := mongoDB.Ref()
	defer mongoDB.UnRef(db)
	err := db.DB(DB).C(USERDB).Find(bson.M{"name": user.Name, "pwd": user.Pwd}).One(&result)
	if err != nil {
		glog.Errorln(err)
		return err
	}
	return nil
}

//检查用户是否已注册过
func (data *UserData) ExistByAccountID() bool {
	db := mongoDB.Ref()
	defer mongoDB.UnRef(db)
	count, _ := db.DB(DB).C(USERDB).Find(bson.M{"accountid": data.AccountID}).Count()
	return count > 0
}
*/
