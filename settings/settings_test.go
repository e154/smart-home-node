package settings_test

import (
	"testing"
	"github.com/e154/smart-home-node/settings"
)

func TestSettingsPtr(t *testing.T) {

	s := settings.SettingsPtr()
	if s == nil {
		t.Errorf("Settings pointer is nil")
	}
}

func TestSettings_Init(t *testing.T) {

	s := settings.SettingsPtr()
	if (s.Init() != s) {
		t.Errorf("Settings pointer is nil")
	}
}

func TestSettings_Load(t *testing.T) {

	//s := settings.SettingsPtr()
	//
	//ns, err := s.Load()
	//if err != nil {
	//	t.Errorf("error %s", err.Error())
	//}
	//
	//if (ns != s) {
	//	t.Errorf("Settings pointer is nil")
	//}
}

func TestSettings_Save(t *testing.T) {

	//s := settings.SettingsPtr()
	//
	//sn, err := s.Save()
	//if err != nil {
	//	t.Errorf("error %s", err.Error())
	//}
	//
	//if (sn != s) {
	//	t.Errorf("Settings pointer is nil")
	//}
}