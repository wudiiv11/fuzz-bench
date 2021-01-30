package db

import (
	"log"
	"userservice/db/mysql"
)

type User struct {
	Username     string
	Passwd       string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

// Signup: 注册用户, 这里的密码是已经加密过的密码
func Signup(username string, passwd string) bool {
	stmt, err := mysql.Conn().Prepare("insert ignore into tbl_user (`user_name`, `user_pwd`) values (?,?)")
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	if rowsAffected, err := ret.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}
	return false
}

// LogIn: 判断密码是否一致, 这里的密码是已经加密过的密码
func LogIn(username string, passwd string) bool {
	stmt, err := mysql.Conn().Prepare("select * from tbl_user where user_name=?")
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		log.Println(err.Error())
		return false
	} else if rows == nil {
		log.Println("username not found " + username)
		return false
	}

	pRows := mysql.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == passwd {
		return true
	}
	return false
}

// UpdateToken: 更新用户token
func UpdateToken(username string, token string) bool {
	stmt, err := mysql.Conn().Prepare("replace into tbl_user_token (user_name, user_token) values (?,?)")
	if err != nil {
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, token)
	if err != nil {
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}
	return false
}

// GetUserInfo: 通过用户名得到用户信息
func GetUserInfo(username string) (User, error) {
	user := User{}

	stmt, err := mysql.Conn().Prepare("select user_name, user_pwd, email, phone, signup_at, last_active, status from tbl_user where user_name=?")
	if err != nil {
		return user, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Username, &user.Passwd, &user.Email, &user.Phone, &user.SignupAt, &user.LastActiveAt, &user.Status)
	if err != nil {
		return user, err
	}
	return user, nil
}
