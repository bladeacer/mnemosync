/*
Copyright © 2025 bladeacer <wg.nick.exe@gmail.com>
*/

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type DirData struct {
	TargetPath string `json:"target_path"`
	Alias string `json:"alias"`
}
type DataStore struct {
	CurrentId int64 `json:"current_id"`
	TrackedDirs map[string]DirData `json:"tracked_dirs"`
}
func GetDataStore() *DataStore {
	return &DataStore{
		CurrentId: 0,
		TrackedDirs: make(map[string]DirData),
	}
}
func LoadDataStore() (*DataStore, error) {
	dbPath := ResolveDbPath()
	ds := GetDataStore()

	data, err := os.ReadFile(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ds, nil
		}
		return nil, fmt.Errorf("error reading database file %s: %w", dbPath, err)
	}
	
	if err := json.Unmarshal(data, ds); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON data from %s: %w", dbPath, err)
	}

	return ds, nil
}
func (ds *DataStore) AddDir(data DirData) string {
    ds.CurrentId += 1
    
    newIDStr := strconv.FormatInt(ds.CurrentId, 10)
    ds.TrackedDirs[newIDStr] = data
    return newIDStr
}
func (ds *DataStore) SaveData(targetPath string) error {
	jsonData, err := json.MarshalIndent(ds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal DataStore to JSON: %w", err)
	}

	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory structure for %s: %w", targetPath, err)
	}
	if err := os.WriteFile(targetPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON data to file %s: %w", targetPath, err)
	}
	return nil
}
