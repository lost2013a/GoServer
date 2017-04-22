package mail

import (
	"dbmgo"
)

type TMailMoudle struct {
	PlayerID uint32 `bson:"_id"`
}

func (self *TMailMoudle) InitWriteDB(id uint32) {
	self.PlayerID = id
	dbmgo.InsertSync("Mail", self)
}
func (self *TMailMoudle) LoadFromDB(id uint32) {
	dbmgo.Find("Mail", "_id", id, self)
}
func (self *TMailMoudle) OnLogin()  {}
func (self *TMailMoudle) OnLogout() {}