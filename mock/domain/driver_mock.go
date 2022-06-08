// Code generated by MockGen. DO NOT EDIT.
// Source: domain/driver.go

// Package domain_mock is a generated GoMock package.
package domain_mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/migratemgr8/mgr8/domain"
)

// MockDriver is a mock of Driver interface.
type MockDriver struct {
	ctrl     *gomock.Controller
	recorder *MockDriverMockRecorder
}

// MockDriverMockRecorder is the mock recorder for MockDriver.
type MockDriverMockRecorder struct {
	mock *MockDriver
}

// NewMockDriver creates a new mock instance.
func NewMockDriver(ctrl *gomock.Controller) *MockDriver {
	mock := &MockDriver{ctrl: ctrl}
	mock.recorder = &MockDriverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDriver) EXPECT() *MockDriverMockRecorder {
	return m.recorder
}

// CreateAppliedMigrationsTable mocks base method.
func (m *MockDriver) CreateAppliedMigrationsTable() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAppliedMigrationsTable")
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAppliedMigrationsTable indicates an expected call of CreateAppliedMigrationsTable.
func (mr *MockDriverMockRecorder) CreateAppliedMigrationsTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAppliedMigrationsTable", reflect.TypeOf((*MockDriver)(nil).CreateAppliedMigrationsTable))
}

// CreateMigrationsLogsTable mocks base method.
func (m *MockDriver) CreateMigrationsLogsTable() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMigrationsLogsTable")
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateMigrationsLogsTable indicates an expected call of CreateMigrationsLogsTable.
func (mr *MockDriverMockRecorder) CreateMigrationsLogsTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMigrationsLogsTable", reflect.TypeOf((*MockDriver)(nil).CreateMigrationsLogsTable))
}

// Deparser mocks base method.
func (m *MockDriver) Deparser() domain.Deparser {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Deparser")
	ret0, _ := ret[0].(domain.Deparser)
	return ret0
}

// Deparser indicates an expected call of Deparser.
func (mr *MockDriverMockRecorder) Deparser() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deparser", reflect.TypeOf((*MockDriver)(nil).Deparser))
}

// DropAppliedMigrationsTable mocks base method.
func (m *MockDriver) DropAppliedMigrationsTable() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DropAppliedMigrationsTable")
	ret0, _ := ret[0].(error)
	return ret0
}

// DropAppliedMigrationsTable indicates an expected call of DropAppliedMigrationsTable.
func (mr *MockDriverMockRecorder) DropAppliedMigrationsTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DropAppliedMigrationsTable", reflect.TypeOf((*MockDriver)(nil).DropAppliedMigrationsTable))
}

// DropMigrationsLogsTable mocks base method.
func (m *MockDriver) DropMigrationsLogsTable() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DropMigrationsLogsTable")
	ret0, _ := ret[0].(error)
	return ret0
}

// DropMigrationsLogsTable indicates an expected call of DropMigrationsLogsTable.
func (mr *MockDriverMockRecorder) DropMigrationsLogsTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DropMigrationsLogsTable", reflect.TypeOf((*MockDriver)(nil).DropMigrationsLogsTable))
}

// Execute mocks base method.
func (m *MockDriver) Execute(statements []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", statements)
	ret0, _ := ret[0].(error)
	return ret0
}

// Execute indicates an expected call of Execute.
func (mr *MockDriverMockRecorder) Execute(statements interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockDriver)(nil).Execute), statements)
}

// ExecuteTransaction mocks base method.
func (m *MockDriver) ExecuteTransaction(url string, f func() error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecuteTransaction", url, f)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExecuteTransaction indicates an expected call of ExecuteTransaction.
func (mr *MockDriverMockRecorder) ExecuteTransaction(url, f interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecuteTransaction", reflect.TypeOf((*MockDriver)(nil).ExecuteTransaction), url, f)
}

// GetLatestMigrationVersion mocks base method.
func (m *MockDriver) GetLatestMigrationVersion() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestMigrationVersion")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestMigrationVersion indicates an expected call of GetLatestMigrationVersion.
func (mr *MockDriverMockRecorder) GetLatestMigrationVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestMigrationVersion", reflect.TypeOf((*MockDriver)(nil).GetLatestMigrationVersion))
}

// GetVersionHashing mocks base method.
func (m *MockDriver) GetVersionHashing(version int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVersionHashing", version)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVersionHashing indicates an expected call of GetVersionHashing.
func (mr *MockDriverMockRecorder) GetVersionHashing(version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVersionHashing", reflect.TypeOf((*MockDriver)(nil).GetVersionHashing), version)
}

// HasAppliedMigrationsTable mocks base method.
func (m *MockDriver) HasAppliedMigrationsTable() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasAppliedMigrationsTable")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasAppliedMigrationsTable indicates an expected call of HasAppliedMigrationsTable.
func (mr *MockDriverMockRecorder) HasAppliedMigrationsTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasAppliedMigrationsTable", reflect.TypeOf((*MockDriver)(nil).HasAppliedMigrationsTable))
}

// HasMigrationLogsTable mocks base method.
func (m *MockDriver) HasMigrationLogsTable() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasMigrationLogsTable")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasMigrationLogsTable indicates an expected call of HasMigrationLogsTable.
func (mr *MockDriverMockRecorder) HasMigrationLogsTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasMigrationLogsTable", reflect.TypeOf((*MockDriver)(nil).HasMigrationLogsTable))
}

// InsertIntoAppliedMigrations mocks base method.
func (m *MockDriver) InsertIntoAppliedMigrations(version int, username, currentDate, hash string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertIntoAppliedMigrations", version, username, currentDate, hash)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertIntoAppliedMigrations indicates an expected call of InsertIntoAppliedMigrations.
func (mr *MockDriverMockRecorder) InsertIntoAppliedMigrations(version, username, currentDate, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertIntoAppliedMigrations", reflect.TypeOf((*MockDriver)(nil).InsertIntoAppliedMigrations), version, username, currentDate, hash)
}

// InsertIntoMigrationLog mocks base method.
func (m *MockDriver) InsertIntoMigrationLog(migrationNum int, migrationType, username, currentDate string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertIntoMigrationLog", migrationNum, migrationType, username, currentDate)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertIntoMigrationLog indicates an expected call of InsertIntoMigrationLog.
func (mr *MockDriverMockRecorder) InsertIntoMigrationLog(migrationNum, migrationType, username, currentDate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertIntoMigrationLog", reflect.TypeOf((*MockDriver)(nil).InsertIntoMigrationLog), migrationNum, migrationType, username, currentDate)
}

// InstallTool mocks base method.
func (m *MockDriver) InstallTool() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstallTool")
	ret0, _ := ret[0].(error)
	return ret0
}

// InstallTool indicates an expected call of InstallTool.
func (mr *MockDriverMockRecorder) InstallTool() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstallTool", reflect.TypeOf((*MockDriver)(nil).InstallTool))
}

// IsToolInstalled mocks base method.
func (m *MockDriver) IsToolInstalled() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsToolInstalled")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsToolInstalled indicates an expected call of IsToolInstalled.
func (mr *MockDriverMockRecorder) IsToolInstalled() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsToolInstalled", reflect.TypeOf((*MockDriver)(nil).IsToolInstalled))
}

// ParseMigration mocks base method.
func (m *MockDriver) ParseMigration(scriptFile string) (*domain.Schema, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseMigration", scriptFile)
	ret0, _ := ret[0].(*domain.Schema)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseMigration indicates an expected call of ParseMigration.
func (mr *MockDriverMockRecorder) ParseMigration(scriptFile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseMigration", reflect.TypeOf((*MockDriver)(nil).ParseMigration), scriptFile)
}

// RemoveAppliedMigration mocks base method.
func (m *MockDriver) RemoveAppliedMigration(version int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAppliedMigration", version)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAppliedMigration indicates an expected call of RemoveAppliedMigration.
func (mr *MockDriverMockRecorder) RemoveAppliedMigration(version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAppliedMigration", reflect.TypeOf((*MockDriver)(nil).RemoveAppliedMigration), version)
}

// UninstallTool mocks base method.
func (m *MockDriver) UninstallTool() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UninstallTool")
	ret0, _ := ret[0].(error)
	return ret0
}

// UninstallTool indicates an expected call of UninstallTool.
func (mr *MockDriverMockRecorder) UninstallTool() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UninstallTool", reflect.TypeOf((*MockDriver)(nil).UninstallTool))
}

// MockDeparser is a mock of Deparser interface.
type MockDeparser struct {
	ctrl     *gomock.Controller
	recorder *MockDeparserMockRecorder
}

// MockDeparserMockRecorder is the mock recorder for MockDeparser.
type MockDeparserMockRecorder struct {
	mock *MockDeparser
}

// NewMockDeparser creates a new mock instance.
func NewMockDeparser(ctrl *gomock.Controller) *MockDeparser {
	mock := &MockDeparser{ctrl: ctrl}
	mock.recorder = &MockDeparserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeparser) EXPECT() *MockDeparserMockRecorder {
	return m.recorder
}

// AddColumn mocks base method.
func (m *MockDeparser) AddColumn(tableName, columnName string, column *domain.Column) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddColumn", tableName, columnName, column)
	ret0, _ := ret[0].(string)
	return ret0
}

// AddColumn indicates an expected call of AddColumn.
func (mr *MockDeparserMockRecorder) AddColumn(tableName, columnName, column interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddColumn", reflect.TypeOf((*MockDeparser)(nil).AddColumn), tableName, columnName, column)
}

// CreateTable mocks base method.
func (m *MockDeparser) CreateTable(table *domain.Table) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTable", table)
	ret0, _ := ret[0].(string)
	return ret0
}

// CreateTable indicates an expected call of CreateTable.
func (mr *MockDeparserMockRecorder) CreateTable(table interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTable", reflect.TypeOf((*MockDeparser)(nil).CreateTable), table)
}

// DropColumn mocks base method.
func (m *MockDeparser) DropColumn(tableName, columnName string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DropColumn", tableName, columnName)
	ret0, _ := ret[0].(string)
	return ret0
}

// DropColumn indicates an expected call of DropColumn.
func (mr *MockDeparserMockRecorder) DropColumn(tableName, columnName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DropColumn", reflect.TypeOf((*MockDeparser)(nil).DropColumn), tableName, columnName)
}

// DropTable mocks base method.
func (m *MockDeparser) DropTable(tableName string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DropTable", tableName)
	ret0, _ := ret[0].(string)
	return ret0
}

// DropTable indicates an expected call of DropTable.
func (mr *MockDeparserMockRecorder) DropTable(tableName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DropTable", reflect.TypeOf((*MockDeparser)(nil).DropTable), tableName)
}

// MakeColumnNotNull mocks base method.
func (m *MockDeparser) MakeColumnNotNull(tableName, columnName string, column *domain.Column) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeColumnNotNull", tableName, columnName, column)
	ret0, _ := ret[0].(string)
	return ret0
}

// MakeColumnNotNull indicates an expected call of MakeColumnNotNull.
func (mr *MockDeparserMockRecorder) MakeColumnNotNull(tableName, columnName, column interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeColumnNotNull", reflect.TypeOf((*MockDeparser)(nil).MakeColumnNotNull), tableName, columnName, column)
}

// MakeColumnNullable mocks base method.
func (m *MockDeparser) MakeColumnNullable(tableName, columnName string, column *domain.Column) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeColumnNullable", tableName, columnName, column)
	ret0, _ := ret[0].(string)
	return ret0
}

// MakeColumnNullable indicates an expected call of MakeColumnNullable.
func (mr *MockDeparserMockRecorder) MakeColumnNullable(tableName, columnName, column interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeColumnNullable", reflect.TypeOf((*MockDeparser)(nil).MakeColumnNullable), tableName, columnName, column)
}

// WriteScript mocks base method.
func (m *MockDeparser) WriteScript(statements []string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteScript", statements)
	ret0, _ := ret[0].(string)
	return ret0
}

// WriteScript indicates an expected call of WriteScript.
func (mr *MockDeparserMockRecorder) WriteScript(statements interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteScript", reflect.TypeOf((*MockDeparser)(nil).WriteScript), statements)
}
