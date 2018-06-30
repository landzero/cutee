package types

import "landzero.net/x/database/orm"

// User user model
type User struct {
	orm.Model
	Avatar  string `orm:"avatar"`
	Account string `orm:"account;unique_index"`
	Name    string `orm:"name;index"`
}
