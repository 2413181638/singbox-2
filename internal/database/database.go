package database

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	*gorm.DB
}

type Server struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Host        string    `json:"host"`
	Port        int       `json:"port"`
	ServerType  string    `json:"server_type"` // vmess, vless, shadowsocks, trojan
	Protocol    string    `json:"protocol"`
	Settings    string    `json:"settings"` // JSON格式的服务器配置
	Tags        string    `json:"tags"`
	IsActive    bool      `json:"is_active"`
	Ping        int       `json:"ping"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Subscription struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Token       string    `json:"token"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastSync    time.Time `json:"last_sync"`
	ServerCount int       `json:"server_count"`
	IsActive    bool      `json:"is_active"`
}

type ConnectionLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ServerID  uint      `json:"server_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	Upload    int64     `json:"upload"`
	Download  int64     `json:"download"`
	Success   bool      `json:"success"`
	Error     string    `json:"error"`
}

func Init(dbPath string) (*Database, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// 自动迁移数据库表
	err = db.AutoMigrate(&Server{}, &Subscription{}, &ConnectionLog{})
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (db *Database) GetServers() ([]Server, error) {
	var servers []Server
	err := db.Find(&servers).Error
	return servers, err
}

func (db *Database) GetActiveServers() ([]Server, error) {
	var servers []Server
	err := db.Where("is_active = ?", true).Find(&servers).Error
	return servers, err
}

func (db *Database) CreateServer(server *Server) error {
	return db.Create(server).Error
}

func (db *Database) UpdateServer(server *Server) error {
	return db.Save(server).Error
}

func (db *Database) DeleteServer(id uint) error {
	return db.Delete(&Server{}, id).Error
}

func (db *Database) GetSubscriptions() ([]Subscription, error) {
	var subscriptions []Subscription
	err := db.Find(&subscriptions).Error
	return subscriptions, err
}

func (db *Database) CreateSubscription(subscription *Subscription) error {
	return db.Create(subscription).Error
}

func (db *Database) UpdateSubscription(subscription *Subscription) error {
	return db.Save(subscription).Error
}

func (db *Database) DeleteSubscription(id uint) error {
	return db.Delete(&Subscription{}, id).Error
}

func (db *Database) CreateConnectionLog(log *ConnectionLog) error {
	return db.Create(log).Error
}

func (db *Database) GetConnectionLogs(limit int) ([]ConnectionLog, error) {
	var logs []ConnectionLog
	err := db.Order("start_time desc").Limit(limit).Find(&logs).Error
	return logs, err
}