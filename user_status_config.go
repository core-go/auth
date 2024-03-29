package auth

type UserStatusConfig struct {
	Activated   string `yaml:"activated" mapstructure:"activated" json:"activated,omitempty" gorm:"column:activated" bson:"activated,omitempty" dynamodbav:"activated,omitempty" firestore:"activated,omitempty"`
	Deactivated string `yaml:"deactivated" mapstructure:"deactivated" json:"deactivated,omitempty" gorm:"column:deactivated" bson:"deactivated,omitempty" dynamodbav:"deactivated,omitempty" firestore:"deactivated,omitempty"`
	Disable     string `yaml:"disable" mapstructure:"disable" json:"disable,omitempty" gorm:"column:disable" bson:"disable,omitempty" dynamodbav:"disable,omitempty" firestore:"disable,omitempty"`
	Suspended   string `yaml:"suspended" mapstructure:"suspended" json:"suspended,omitempty" gorm:"column:suspended" bson:"suspended,omitempty" dynamodbav:"suspended,omitempty" firestore:"suspended,omitempty"`
}
