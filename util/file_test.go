// Copyright 2013 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package util_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gc "gopkg.in/check.v1"

	"github.com/zlxtqbdgdgd/sailor/util"
)

func Test(t *testing.T) { gc.TestingT(t) }

type fileSuite struct {
}

var _ = gc.Suite(&fileSuite{})

func (*fileSuite) TestCopyFile(c *gc.C) {
	dir := c.MkDir()
	f, err := ioutil.TempFile(dir, "source")
	c.Assert(err, gc.IsNil)
	defer f.Close()
	_, err = f.Write([]byte("hello world"))
	c.Assert(err, gc.IsNil)
	dest := filepath.Join(dir, "dest")

	err = util.CopyFile(dest, f.Name())
	c.Assert(err, gc.IsNil)
	data, err := ioutil.ReadFile(dest)
	c.Assert(err, gc.IsNil)
	c.Assert(string(data), gc.Equals, "hello world")
}

var atomicWriteFileTests = []struct {
	summary   string
	change    func(filename string, contents []byte) error
	check     func(c *gc.C, fileInfo os.FileInfo)
	expectErr string
}{{
	summary: "atomic file write and chmod 0644",
	change: func(filename string, contents []byte) error {
		return util.AtomicWriteFile(filename, contents, 0765)
	},
	check: func(c *gc.C, fi os.FileInfo) {
		c.Assert(fi.Mode(), gc.Equals, 0765)
	},
}, {
	summary: "atomic file write and change",
	change: func(filename string, contents []byte) error {
		chmodChange := func(f string) error {
			// FileMod.Chmod() is not implemented on Windows, however, os.Chmod() is
			return os.Chmod(f, 0700)
		}
		return util.AtomicWriteFileAndChange(filename, contents, chmodChange)
	},
	check: func(c *gc.C, fi os.FileInfo) {
		c.Assert(fi.Mode(), gc.Equals, 0700)
	},
}, {
	summary: "atomic file write empty contents",
	change: func(filename string, contents []byte) error {
		nopChange := func(string) error {
			return nil
		}
		return util.AtomicWriteFileAndChange(filename, contents, nopChange)
	},
}, {
	summary: "atomic file write and failing change func",
	change: func(filename string, contents []byte) error {
		errChange := func(string) error {
			return fmt.Errorf("pow!")
		}
		return util.AtomicWriteFileAndChange(filename, contents, errChange)
	},
	expectErr: "pow!",
}}

func (*fileSuite) TestAtomicWriteFile(c *gc.C) {
	dir := c.MkDir()
	name := "test.file"
	path := filepath.Join(dir, name)
	assertDirContents := func(names ...string) {
		fis, err := ioutil.ReadDir(dir)
		c.Assert(err, gc.IsNil)
		c.Assert(fis, gc.HasLen, len(names))
		for i, name := range names {
			c.Assert(fis[i].Name(), gc.Equals, name)
		}
	}
	/*assertNotExist := func(path string) {
		_, err := os.Lstat(path)
		c.Assert(err, jc.Satisfies, os.IsNotExist)
	}*/

	for i, test := range atomicWriteFileTests {
		c.Logf("test %d: %s", i, test.summary)
		// First - test with file not already there.
		assertDirContents()
		//assertNotExist(path)
		contents := []byte("some\ncontents")

		err := test.change(path, contents)
		if test.expectErr == "" {
			c.Assert(err, gc.IsNil)
			data, err := ioutil.ReadFile(path)
			c.Assert(err, gc.IsNil)
			c.Assert(data, gc.DeepEquals, contents)
			assertDirContents(name)
		} else {
			c.Assert(err, gc.ErrorMatches, test.expectErr)
			assertDirContents()
			continue
		}

		// Second - test with a file already there.
		contents = []byte("new\ncontents")
		err = test.change(path, contents)
		c.Assert(err, gc.IsNil)
		data, err := ioutil.ReadFile(path)
		c.Assert(err, gc.IsNil)
		c.Assert(data, gc.DeepEquals, contents)
		assertDirContents(name)

		// Remove the file to reset scenario.
		c.Assert(os.Remove(path), gc.IsNil)
	}
}

/*func (*fileSuite) TestMoveFile(c *gc.C) {
	d := c.MkDir()
	dest := filepath.Join(d, "foo")
	f1Name := filepath.Join(d, ".foo1")
	f2Name := filepath.Join(d, ".foo2")
	err := ioutil.WriteFile(f1Name, []byte("macaroni"), 0644)
	c.Assert(err, gc.IsNil)
	err = ioutil.WriteFile(f2Name, []byte("cheese"), 0644)
	c.Assert(err, gc.IsNil)

	ok, err := util.MoveFile(f1Name, dest)
	c.Assert(ok, gc.Equals, true)
	c.Assert(err, gc.IsNil)

	ok, err = util.MoveFile(f2Name, dest)
	c.Assert(ok, gc.Equals, false)
	c.Assert(err, gc.NotNil)

	contents, err := ioutil.ReadFile(dest)
	c.Assert(err, gc.IsNil)
	c.Assert(contents, gc.DeepEquals, []byte("macaroni"))
}*/
