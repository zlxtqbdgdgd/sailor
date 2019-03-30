// Copyright 2018 JXB. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// CopyFile writes the contents of the given source file to dest.
func CopyFile(dest, source string) error {
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	f, err := os.Open(source)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(df, f)
	return err
}

/*// MoveFile atomically moves the source file to the destination, returning
// whether the file was moved successfully. If the destination already exists,
// it returns an error rather than overwrite it.
//
// On unix systems, an error may occur with a successful move, if the source
// file location cannot be unlinked.
func MoveFile(source, destination string) (bool, error) {
	err := os.Link(source, destination)
	if err != nil {
		return false, err
	}
	err = os.Remove(source)
	if err != nil {
		return true, err
	}
	return true, nil
}*/

// ReplaceFile atomically replaces the destination file or directory
// with the source. The errors that are returned are identical to
// those returned by os.Rename.
func ReplaceFile(source, destination string) error {
	return os.Rename(source, destination)
}

// AtomicWriteFileAndChange atomically writes the filename with the
// given contents and calls the given function after the contents were
// written, but before the file is renamed.
func AtomicWriteFileAndChange(filename string, contents []byte, change func(string) error) (err error) {
	dir, file := filepath.Split(filename)
	f, err := ioutil.TempFile(dir, file)
	if err != nil {
		return fmt.Errorf("cannot create temp file: %v", err)
	}
	defer f.Close()
	defer func() {
		if err != nil {
			// Don't leave the temp file lying around on error.
			// Close the file before removing. Trying to remove an open file on
			// Windows will fail.
			f.Close()
			os.Remove(f.Name())
		}
	}()
	if _, err := f.Write(contents); err != nil {
		return fmt.Errorf("cannot write %q contents: %v", filename, err)
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := change(f.Name()); err != nil {
		return err
	}
	if err := ReplaceFile(f.Name(), filename); err != nil {
		return fmt.Errorf("cannot replace %q with %q: %v", f.Name(), filename, err)
	}
	return nil
}

// AtomicWriteFile atomically writes the filename with the given
// contents and permissions, replacing any existing file at the same
// path.
func AtomicWriteFile(filename string, contents []byte, perms os.FileMode) (err error) {
	return AtomicWriteFileAndChange(filename, contents, func(f string) error {
		// FileMod.Chmod() is not implemented on Windows, however, os.Chmod() is
		if err := os.Chmod(f, perms); err != nil {
			return fmt.Errorf("cannot set permissions: %v", err)
		}
		return nil
	})
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func LoadFileToList(fileName string) ([]string, error) {
	var lines []string
	if !CheckFileIsExist(fileName) {
		return lines, NewErrf("do not find file %s", fileName)
	}

	fHandler, err := os.Open(fileName)
	if err != nil {
		return lines, err
	}

	rd := bufio.NewReader(fHandler)

	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if line == "\n" {
			continue
		}
		lines = append(lines, strings.TrimSpace(line))
	}

	return lines, nil
}
