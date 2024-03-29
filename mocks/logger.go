// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/waltzofpearls/reckon/logs (interfaces: Logger)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	zap "go.uber.org/zap"
	zapcore "go.uber.org/zap/zapcore"
)

// Logger is a mock of Logger interface.
type Logger struct {
	ctrl     *gomock.Controller
	recorder *LoggerMockRecorder
}

// LoggerMockRecorder is the mock recorder for Logger.
type LoggerMockRecorder struct {
	mock *Logger
}

// NewLogger creates a new mock instance.
func NewLogger(ctrl *gomock.Controller) *Logger {
	mock := &Logger{ctrl: ctrl}
	mock.recorder = &LoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Logger) EXPECT() *LoggerMockRecorder {
	return m.recorder
}

// Core mocks base method.
func (m *Logger) Core() zapcore.Core {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Core")
	ret0, _ := ret[0].(zapcore.Core)
	return ret0
}

// Core indicates an expected call of Core.
func (mr *LoggerMockRecorder) Core() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Core", reflect.TypeOf((*Logger)(nil).Core))
}

// Debug mocks base method.
func (m *Logger) Debug(arg0 string, arg1 ...zapcore.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debug", varargs...)
}

// Debug indicates an expected call of Debug.
func (mr *LoggerMockRecorder) Debug(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*Logger)(nil).Debug), varargs...)
}

// Error mocks base method.
func (m *Logger) Error(arg0 string, arg1 ...zapcore.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error.
func (mr *LoggerMockRecorder) Error(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*Logger)(nil).Error), varargs...)
}

// Fatal mocks base method.
func (m *Logger) Fatal(arg0 string, arg1 ...zapcore.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Fatal", varargs...)
}

// Fatal indicates an expected call of Fatal.
func (mr *LoggerMockRecorder) Fatal(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fatal", reflect.TypeOf((*Logger)(nil).Fatal), varargs...)
}

// Info mocks base method.
func (m *Logger) Info(arg0 string, arg1 ...zapcore.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *LoggerMockRecorder) Info(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*Logger)(nil).Info), varargs...)
}

// Named mocks base method.
func (m *Logger) Named(arg0 string) *zap.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Named", arg0)
	ret0, _ := ret[0].(*zap.Logger)
	return ret0
}

// Named indicates an expected call of Named.
func (mr *LoggerMockRecorder) Named(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Named", reflect.TypeOf((*Logger)(nil).Named), arg0)
}

// Panic mocks base method.
func (m *Logger) Panic(arg0 string, arg1 ...zapcore.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Panic", varargs...)
}

// Panic indicates an expected call of Panic.
func (mr *LoggerMockRecorder) Panic(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Panic", reflect.TypeOf((*Logger)(nil).Panic), varargs...)
}

// Sync mocks base method.
func (m *Logger) Sync() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync")
	ret0, _ := ret[0].(error)
	return ret0
}

// Sync indicates an expected call of Sync.
func (mr *LoggerMockRecorder) Sync() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*Logger)(nil).Sync))
}

// Warn mocks base method.
func (m *Logger) Warn(arg0 string, arg1 ...zapcore.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warn", varargs...)
}

// Warn indicates an expected call of Warn.
func (mr *LoggerMockRecorder) Warn(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warn", reflect.TypeOf((*Logger)(nil).Warn), varargs...)
}

// With mocks base method.
func (m *Logger) With(arg0 ...zapcore.Field) *zap.Logger {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "With", varargs...)
	ret0, _ := ret[0].(*zap.Logger)
	return ret0
}

// With indicates an expected call of With.
func (mr *LoggerMockRecorder) With(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "With", reflect.TypeOf((*Logger)(nil).With), arg0...)
}
