package cassandra

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/core-go/auth"
	"github.com/gocql/gocql"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson"
	// "github.com/core-go/auth"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
	// "log"
	// "reflect"
	// "strconv"
	// "strings"
	// "time"
)

type AuthenticationRepository struct {
	UserCassandra           *gocql.ClusterConfig
	userTableName           string
	passwordTableName       string
	CheckTwoFactors         func(ctx context.Context, id string) (bool, error)
	activatedStatus         interface{}
	Status                  auth.UserStatusConfig
	IdName                  string
	UserName                string
	UserId                  string
	SuccessTimeName         string
	FailTimeName            string
	FailCountName           string
	LockedUntilTimeName     string
	StatusName              string
	PasswordChangedTimeName string
	PasswordName            string
	ContactName             string
	EmailName               string
	PhoneName               string
	DisplayNameName         string
	MaxPasswordAgeName      string
	UserTypeName            string
	AccessDateFromName      string
	AccessDateToName        string
	AccessTimeFromName      string
	AccessTimeToName        string
	TwoFactorsName          string
}

func NewAuthenticationRepositoryByConfig(db *gocql.ClusterConfig, userTableName, passwordTableName string, activatedStatus string, status auth.UserStatusConfig, c auth.SchemaConfig, options ...func(context.Context, string) (bool, error)) *AuthenticationRepository {
	return NewAuthenticationRepository(db, userTableName, passwordTableName, activatedStatus, status, c.Id, c.Username, c.UserId, c.SuccessTime, c.FailTime, c.FailCount, c.LockedUntilTime, c.Status, c.PasswordChangedTime, c.Password, c.Contact, c.Email, c.Phone, c.DisplayName, c.MaxPasswordAge, c.UserType, c.AccessDateFrom, c.AccessDateTo, c.AccessTimeFrom, c.AccessTimeTo, c.TwoFactors, options...)
}

func NewAuthenticationRepository(db *gocql.ClusterConfig, userTableName, passwordTableName string, activatedStatus string, status auth.UserStatusConfig, idName, userName, userID, successTimeName, failTimeName, failCountName, lockedUntilTimeName, statusName, passwordChangedTimeName, passwordName, contactName, emailName, phoneName, displayNameName, maxPasswordAgeName, userTypeName, accessDateFromName, accessDateToName, accessTimeFromName, accessTimeToName, twoFactorsName string, options ...func(context.Context, string) (bool, error)) *AuthenticationRepository {
	var checkTwoFactors func(context.Context, string) (bool, error)
	if len(options) >= 1 {
		checkTwoFactors = options[0]
	}
	return &AuthenticationRepository{
		UserCassandra:           db,
		userTableName:           strings.ToLower(userTableName),
		passwordTableName:       strings.ToLower(passwordTableName),
		CheckTwoFactors:         checkTwoFactors,
		activatedStatus:         strings.ToLower(activatedStatus),
		Status:                  status,
		IdName:                  strings.ToLower(idName),
		UserName:                strings.ToLower(userName),
		UserId:                  strings.ToLower(userID),
		SuccessTimeName:         strings.ToLower(successTimeName),
		FailTimeName:            strings.ToLower(failTimeName),
		FailCountName:           strings.ToLower(failCountName),
		LockedUntilTimeName:     strings.ToLower(lockedUntilTimeName),
		StatusName:              strings.ToLower(statusName),
		PasswordChangedTimeName: strings.ToLower(passwordChangedTimeName),
		PasswordName:            strings.ToLower(passwordName),
		ContactName:             strings.ToLower(contactName),
		EmailName:               strings.ToLower(emailName),
		PhoneName:               strings.ToLower(phoneName),
		DisplayNameName:         strings.ToLower(displayNameName),
		MaxPasswordAgeName:      strings.ToLower(maxPasswordAgeName),
		UserTypeName:            strings.ToLower(userTypeName),
		AccessDateFromName:      strings.ToLower(accessDateFromName),
		AccessDateToName:        strings.ToLower(accessDateToName),
		AccessTimeFromName:      strings.ToLower(accessTimeFromName),
		AccessTimeToName:        strings.ToLower(accessTimeToName),
		TwoFactorsName:          strings.ToLower(twoFactorsName),
	}
}

func (r *AuthenticationRepository) GetUserInfo(ctx context.Context, user string) (*auth.UserInfo, error) {
	userInfo := auth.UserInfo{}
	session, er0 := r.UserCassandra.CreateSession()
	if er0 != nil {
		return nil, er0
	}
	query := "SELECT * FROM " + r.userTableName + " WHERE " + r.UserName + " = ? ALLOW FILTERING"
	raws := session.Query(query, user).Iter()
	for {
		// New map each iteration
		row := make(map[string]interface{})
		if !raws.MapScan(row) {
			break
		}
		// Do things with row
		if id, ok := row["userid"]; ok {
			userInfo.Id = id.(string)
		}
		if len(r.StatusName) > 0 {
			// rawStatus := raw.Lookup(r.StatusName)
			if status, ok := row[r.StatusName]; ok {
				r.StatusName = status.(string)
			}
		}
		if len(r.ContactName) > 0 {
			if contact, ok := row[r.ContactName]; ok {
				userInfo.Contact = contact.(string)
			}
		}
		if len(r.EmailName) > 0 {
			if email, ok := row[r.EmailName]; ok {
				userInfo.Email = email.(string)
			}
		}
		if len(r.PhoneName) > 0 {
			if phone, ok := row[r.PhoneName]; ok {
				userInfo.Phone = phone.(string)
			}
		}
		if len(r.DisplayNameName) > 0 {
			if displayName, ok := row[r.DisplayNameName]; ok {
				userInfo.DisplayName = displayName.(string)
			}
		}
		if len(r.MaxPasswordAgeName) > 0 {
			if maxPasswordAgeName, ok := row[r.MaxPasswordAgeName]; ok {
				userInfo.MaxPasswordAge = int32(maxPasswordAgeName.(int))
			}
		}
		if len(r.UserTypeName) > 0 {
			if userType, ok := row[r.UserTypeName]; ok {
				userInfo.UserType = userType.(string)
			}
		}
		if len(r.AccessDateFromName) > 0 {
			if accessDateFrom, ok := row[r.AccessDateFromName]; ok {
				userInfo.AccessDateFrom = accessDateFrom.(*time.Time)
			}
		}
		if len(r.AccessDateToName) > 0 {
			if accessDateTo, ok := row[r.AccessDateToName]; ok {
				userInfo.AccessDateTo = accessDateTo.(*time.Time)
			}
		}
		if len(r.AccessTimeFromName) > 0 {
			if accessTimeFrom, ok := row[r.AccessTimeFromName]; ok {
				userInfo.AccessTimeFrom = accessTimeFrom.(*time.Time)
			}
		}
		if len(r.AccessTimeToName) > 0 {
			if accessTimeTo, ok := row[r.AccessTimeToName]; ok {
				userInfo.AccessTimeTo = accessTimeTo.(*time.Time)
			}
		}
	}
	queryPasswordTable := "Select * From " + r.passwordTableName + " WHERE userid = ? ALLOW FILTERING"
	rawPassword := session.Query(queryPasswordTable, userInfo.Id).Iter()
	for {
		row := make(map[string]interface{})
		if !rawPassword.MapScan(row) {
			break
		}
		if len(r.PasswordName) > 0 {
			if pass, ok := row[r.PasswordName]; ok {
				userInfo.Password = pass.(string)
			}
		}
		if len(r.LockedUntilTimeName) > 0 {
			if lockedUntilTime, ok := row[r.LockedUntilTimeName]; ok {
				a1 := lockedUntilTime.(time.Time)
				userInfo.LockedUntilTime = &a1
			}
		}
		if len(r.SuccessTimeName) > 0 {
			if successTime, ok := row[r.SuccessTimeName]; ok {
				a2 := successTime.(time.Time)
				userInfo.SuccessTime = &a2
			}
		}
		if len(r.FailTimeName) > 0 {
			if failTime, ok := row[r.FailTimeName]; ok {
				a3 := failTime.(time.Time)
				userInfo.FailTime = &a3
			}
		}

		if len(r.FailCountName) > 0 {
			if failCountName, ok := row[r.FailCountName]; ok {
				userInfo.FailCount = failCountName.(int)
			}
		}

		if len(r.PasswordChangedTimeName) > 0 {
			if passwordChangedTime, ok := row[r.PasswordChangedTimeName]; ok {
				a4 := passwordChangedTime.(time.Time)
				userInfo.PasswordChangedTime = &a4
			}
		}
	}
	return &userInfo, nil
}

func (r *AuthenticationRepository) Pass(ctx context.Context, userId string) (int64, error) {
	return r.passAuthenticationAndActivate(ctx, userId, false)
}
func (r *AuthenticationRepository) PassAndActivate(ctx context.Context, userId string) (int64, error) {
	return r.passAuthenticationAndActivate(ctx, userId, true)
}

func (r *AuthenticationRepository) passAuthenticationAndActivate(ctx context.Context, userId string, updateStatus bool) (int64, error) {
	if len(r.SuccessTimeName) == 0 && len(r.FailCountName) == 0 && len(r.LockedUntilTimeName) == 0 {
		if !updateStatus {
			return 0, nil
		} else if len(r.StatusName) == 0 {
			return 0, nil
		}
	}
	pass := make(map[string]interface{})
	if len(r.SuccessTimeName) > 0 {
		pass[r.SuccessTimeName] = time.Now()
	}
	if len(r.FailCountName) > 0 {
		pass[r.FailCountName] = 0
	}
	if len(r.LockedUntilTimeName) > 0 {
		pass[r.LockedUntilTimeName] = nil
	}
	query := map[string]interface{}{
		r.IdName: userId,
	}
	if !updateStatus {
		return patch(ctx, r.UserCassandra, r.passwordTableName, pass, query)
	}

	if r.userTableName == r.passwordTableName {
		pass[r.StatusName] = r.activatedStatus
		return patch(ctx, r.UserCassandra, r.passwordTableName, pass, query)
	}

	k1, err := patch(ctx, r.UserCassandra, r.passwordTableName, pass, query)
	if err != nil {
		return k1, err
	}

	user := make(map[string]interface{})
	user[r.IdName] = userId
	user[r.StatusName] = r.activatedStatus
	k2, err1 := patch(ctx, r.UserCassandra, r.userTableName, user, query)
	return k1 + k2, err1
}

func (r *AuthenticationRepository) Fail(ctx context.Context, userId string, failCount int, lockedUntil *time.Time) error {
	if len(r.FailTimeName) == 0 && len(r.FailCountName) == 0 && len(r.LockedUntilTimeName) == 0 {
		return nil
	}
	pass := make(map[string]interface{})
	pass[r.IdName] = userId
	if len(r.FailTimeName) > 0 {
		pass[r.FailTimeName] = time.Now()
	}
	if len(r.FailCountName) > 0 {
		pass[r.FailCountName] = failCount
		if len(r.LockedUntilTimeName) > 0 {
			pass[r.LockedUntilTimeName] = lockedUntil
		}
	}
	query := map[string]interface{}{
		r.IdName: userId,
	}
	_, err := patch(ctx, r.UserCassandra, r.passwordTableName, pass, query)
	return err
}

func patch(ctx context.Context, db *gocql.ClusterConfig, table string, model map[string]interface{}, query map[string]interface{}) (int64, error) {
	session, er0 := db.CreateSession()
	if er0 != nil {
		return 0, er0
	}
	keyUpdate := ""
	keyValue := ""
	for k, v := range query {
		keyUpdate = k
		keyValue = fmt.Sprintf("%v", v)
	}
	str := "SELECT * FROM " + table + " WHERE " + keyUpdate + " = ? ALLOW FILTERING"
	rows := session.Query(str, keyValue).Iter()
	for k, _ := range model {
		flag := false
		for row := range rows.Columns() {
			if rows.Columns()[row].Name == k {
				flag = true
			}
		}
		if !flag {
			if k == "failtime" || k == "lockeduntiltime" || k == "successtime" {
				queryAddCol := "ALTER TABLE " + table + " ADD " + k + " timestamp"
				er0 = session.Query(queryAddCol).Exec()
				if er0 != nil {
					return 0, er0
				}
			} else {
				queryAddCol := "ALTER TABLE " + table + " ADD " + k + " int"
				er0 = session.Query(queryAddCol).Exec()
				if er0 != nil {
					return 0, er0
				}
			}
		}
	}
	objectUpdate := make([]string, 0)
	objectUpdateValue := make([]interface{}, 0)
	for k, v := range model {
		objectUpdate = append(objectUpdate, fmt.Sprintf("%s = ? ", k))
		objectUpdateValue = append(objectUpdateValue, v)
	}
	for k, v := range query {
		keyUpdate = k
		keyValue = fmt.Sprintf("'%v'", v)
	}
	strSql := `UPDATE ` + table + ` SET ` + strings.Join(objectUpdate, ",") + ` WHERE ` + keyUpdate + " = " + keyValue
	result := session.Query(strSql, objectUpdateValue...)
	if result.Exec() != nil {
		log.Println(result.Exec())
		return 0, result.Exec()
	}
	return 1, nil
}
