package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ─── GetStatusText ────────────────────────────────────────

func TestAssetStatus_GetStatusText(t *testing.T) {
	tests := []struct {
		status AssetStatus
		want   string
	}{
		{AssetStatusActive, "ใช้งานอยู่"},
		{AssetStatusDefective, "ชำรุด"},
		{AssetStatusWaitDecom, "รอปลดระวาง"},
		{AssetStatusDecommission, "ปลดระวางแล้ว"},
		{AssetStatusActiveReadyToSell, "พร้อมขาย"},
		{AssetStatusMissing, "สูญหาย"},
		{AssetStatusPlanToReplace, "รอเปลี่ยนใหม่"},
		{AssetStatus("unknown"), "ไม่ทราบสถานะ"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			got := tt.status.GetStatusText()
			assert.Equal(t, tt.want, got)
		})
	}
}

// ─── GetColor ─────────────────────────────────────────────

func TestAssetStatus_GetColor(t *testing.T) {
	tests := []struct {
		status AssetStatus
		want   string
	}{
		{AssetStatusActive, "#4CAF50"},
		{AssetStatusDefective, "#EF5350"},
		{AssetStatusWaitDecom, "#FFA726"},
		{AssetStatusDecommission, "#78909C"},
		{AssetStatusActiveReadyToSell, "#42A5F5"},
		{AssetStatusMissing, "#E53935"},
		{AssetStatusPlanToReplace, "#AB47BC"},
		{AssetStatus("unknown"), "#78909C"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			got := tt.status.GetColor()
			assert.Equal(t, tt.want, got)
		})
	}
}

// ─── TableName ────────────────────────────────────────────

func TestEquipment_TableName(t *testing.T) {
	e := Equipment{}
	assert.Equal(t, "equipments", e.TableName())
}
