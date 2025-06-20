package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sony/sonyflake"

	"time"
)

var sf *sonyflake.Sonyflake

func init() {
	sf = sonyflake.NewSonyflake(sonyflake.Settings{
		MachineID: func() (uint16, error) {
			// 返回一个唯一的机器 ID
			return 1, nil
		},
		StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	if sf == nil {
		log.Fatal("Failed to initialize Sonyflake")
	}
}

func GenerateID() (uint64, error) {
	if sf == nil {
		return 0, fmt.Errorf("Sonyflake instance is nil")
	}
	return sf.NextID()
}
