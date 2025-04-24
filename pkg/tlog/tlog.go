/*
Copyright 2025 The  Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// tlog Terminal Logger, simple enough for krb.
package tlog

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	timeFormat = "060102 15:04:05.000000"
)

func Print(msg string) {
	fmt.Print(msg)
}

func Println(msg string) {
	fmt.Println(msg)
}

func Printf(format string, args ...any) {
	Println(strings.TrimSuffix(fmt.Sprintf(format, args...), "\n"))
}

func Panic(msg string) {
	Println(msg)
	os.Exit(1)
}

func Panicf(format string, args ...any) {
	Println(fmt.Sprintf(format, args...))
	os.Exit(1)
}

func Debug(msg string) {
	Printf("%s D %s", time.Now().Format(timeFormat), msg)
}

func Debugf(format string, args ...any) {
	Printf("%s D %s", time.Now().Format(timeFormat), fmt.Sprintf(format, args...))
}

func Info(msg string) {
	Printf("%s I %s", time.Now().Format(timeFormat), msg)
}

func Infof(format string, args ...any) {
	Printf("%s I %s", time.Now().Format(timeFormat), fmt.Sprintf(format, args...))
}

func Warn(msg string) {
	Printf("%s W %s", time.Now().Format(timeFormat), msg)
}

func Warnf(format string, args ...any) {
	Printf("%s W %s", time.Now().Format(timeFormat), fmt.Sprintf(format, args...))
}

func Error(msg string) {
	Printf("%s E %s", time.Now().Format(timeFormat), msg)
}

func Errorf(format string, args ...any) {
	Printf("%s E %s", time.Now().Format(timeFormat), fmt.Sprintf(format, args...))
}

func Fatal(msg string) {
	Panicf("%s F %s", time.Now().Format(timeFormat), msg)
	os.Exit(1)
}

func Fatalf(format string, args ...any) {
	Panicf("%s F %s", time.Now().Format(timeFormat), fmt.Sprintf(format, args...))
	os.Exit(1)
}
