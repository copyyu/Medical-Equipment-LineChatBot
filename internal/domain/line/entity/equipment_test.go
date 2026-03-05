package entity

import (
	"testing"
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
			if got != tt.want {
				t.Errorf("GetStatusText(%s) = %q, want %q", tt.status, got, tt.want)
			}
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
			if got != tt.want {
				t.Errorf("GetColor(%s) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

// ─── TableName ────────────────────────────────────────────

func TestEquipment_TableName(t *testing.T) {
	e := Equipment{}
	if got := e.TableName(); got != "equipments" {
		t.Errorf("TableName() = %q, want %q", got, "equipments")
	}
}
