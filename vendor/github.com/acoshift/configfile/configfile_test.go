package configfile_test

import (
	"testing"

	"github.com/acoshift/configfile"
	"github.com/stretchr/testify/assert"
)

func TestConfigfile(t *testing.T) {
	c := configfile.NewReader("testdata")

	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()

		t.Run("Bool", func(t *testing.T) {
			t.Parallel()

			assert.False(t, c.Bool("notfound"))
			assert.False(t, c.BoolDefault("notfound", false))
			assert.True(t, c.BoolDefault("notfound", true))
			assert.Panics(t, func() { c.MustBool("notfound") })
		})

		t.Run("Int", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, 0, c.Int("notfound"))
			assert.Equal(t, 0, c.IntDefault("notfound", 0))
			assert.Equal(t, 1, c.IntDefault("notfound", 1))
			assert.Panics(t, func() { c.MustInt("notfound") })
		})

		t.Run("String", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, "", c.String("notfound"))
			assert.Equal(t, "", c.StringDefault("notfound", ""))
			assert.Equal(t, "a string", c.StringDefault("notfound", "a string"))
			assert.Panics(t, func() { c.MustString("notfound") })
		})

		t.Run("Bytes", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, []byte{}, c.Bytes("notfound"))
			assert.Nil(t, c.BytesDefault("notfound", nil))
			assert.Equal(t, []byte{}, c.BytesDefault("notfound", []byte{}))
			assert.Equal(t, []byte("some bytes"), c.BytesDefault("notfound", []byte("some bytes")))
			assert.Panics(t, func() { c.MustBytes("notfound") })
		})
	})

	t.Run("Empty", func(t *testing.T) {
		t.Parallel()

		t.Run("Bool", func(t *testing.T) {
			t.Parallel()

			assert.False(t, c.Bool("empty"))
			assert.False(t, c.BoolDefault("empty", false))
			assert.True(t, c.BoolDefault("empty", true))
			assert.Panics(t, func() { c.MustBool("empty") })
		})

		t.Run("Int", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, 0, c.Int("empty"))
			assert.Equal(t, 0, c.IntDefault("empty", 0))
			assert.Equal(t, 1, c.IntDefault("empty", 1))
			assert.Panics(t, func() { c.MustInt("empty") })
		})

		t.Run("String", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, "", c.String("empty"))
			assert.Equal(t, "", c.StringDefault("empty", ""))
			assert.Equal(t, "", c.StringDefault("empty", "a string"))
			assert.NotPanics(t, func() { c.MustString("empty") })
		})

		t.Run("Bytes", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, []byte{}, c.Bytes("empty"))
			assert.Equal(t, []byte{}, c.BytesDefault("empty", nil))
			assert.Equal(t, []byte{}, c.BytesDefault("empty", []byte{}))
			assert.Equal(t, []byte{}, c.BytesDefault("empty", []byte("some bytes")))
			assert.NotPanics(t, func() { c.MustBytes("empty") })
		})
	})

	t.Run("Data1", func(t *testing.T) {
		t.Parallel()

		t.Run("Bool", func(t *testing.T) {
			t.Parallel()

			assert.True(t, c.Bool("data1"))
			assert.True(t, c.BoolDefault("data1", false))
			assert.True(t, c.BoolDefault("data1", true))
			assert.NotPanics(t, func() { c.MustBool("data1") })
		})

		t.Run("Int", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, 0, c.Int("data1"))
			assert.Equal(t, 0, c.IntDefault("data1", 0))
			assert.Equal(t, 1, c.IntDefault("data1", 1))
			assert.Panics(t, func() { c.MustInt("data1") })
		})

		t.Run("String", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, "true", c.String("data1"))
			assert.Equal(t, "true", c.StringDefault("data1", ""))
			assert.Equal(t, "true", c.StringDefault("data1", "a string"))
			assert.NotPanics(t, func() { c.MustString("data1") })
		})

		t.Run("Bytes", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, []byte("true"), c.Bytes("data1"))
			assert.Equal(t, []byte("true"), c.BytesDefault("data1", nil))
			assert.Equal(t, []byte("true"), c.BytesDefault("data1", []byte{}))
			assert.Equal(t, []byte("true"), c.BytesDefault("data1", []byte("some bytes")))
			assert.NotPanics(t, func() { c.MustBytes("data1") })
		})
	})

	t.Run("Data2", func(t *testing.T) {
		t.Parallel()

		t.Run("Bool", func(t *testing.T) {
			t.Parallel()

			assert.False(t, c.Bool("data2"))
			assert.False(t, c.BoolDefault("data2", false))
			assert.False(t, c.BoolDefault("data2", true))
			assert.NotPanics(t, func() { c.MustBool("data2") })
		})

		t.Run("Int", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, 0, c.Int("data2"))
			assert.Equal(t, 0, c.IntDefault("data2", 0))
			assert.Equal(t, 1, c.IntDefault("data2", 1))
			assert.Panics(t, func() { c.MustInt("data2") })
		})

		t.Run("String", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, "false", c.String("data2"))
			assert.Equal(t, "false", c.StringDefault("data2", ""))
			assert.Equal(t, "false", c.StringDefault("data2", "a string"))
			assert.NotPanics(t, func() { c.MustString("data2") })
		})

		t.Run("Bytes", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, []byte("false"), c.Bytes("data2"))
			assert.Equal(t, []byte("false"), c.BytesDefault("data2", nil))
			assert.Equal(t, []byte("false"), c.BytesDefault("data2", []byte{}))
			assert.Equal(t, []byte("false"), c.BytesDefault("data2", []byte("some bytes")))
			assert.NotPanics(t, func() { c.MustBytes("data2") })
		})
	})

	t.Run("Data3", func(t *testing.T) {
		t.Parallel()

		t.Run("Bool", func(t *testing.T) {
			t.Parallel()

			assert.True(t, c.Bool("data3"))
			assert.True(t, c.BoolDefault("data3", false))
			assert.True(t, c.BoolDefault("data3", true))
			assert.NotPanics(t, func() { c.MustBool("data3") })
		})

		t.Run("Int", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, 9, c.Int("data3"))
			assert.Equal(t, 9, c.IntDefault("data3", 0))
			assert.Equal(t, 9, c.IntDefault("data3", 1))
			assert.NotPanics(t, func() { c.MustInt("data3") })
		})

		t.Run("String", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, "9", c.String("data3"))
			assert.Equal(t, "9", c.StringDefault("data3", ""))
			assert.Equal(t, "9", c.StringDefault("data3", "a string"))
			assert.NotPanics(t, func() { c.MustString("data3") })
		})

		t.Run("Bytes", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, []byte("9"), c.Bytes("data3"))
			assert.Equal(t, []byte("9"), c.BytesDefault("data3", nil))
			assert.Equal(t, []byte("9"), c.BytesDefault("data3", []byte{}))
			assert.Equal(t, []byte("9"), c.BytesDefault("data3", []byte("some bytes")))
			assert.NotPanics(t, func() { c.MustBytes("data3") })
		})
	})

	t.Run("Data4", func(t *testing.T) {
		t.Parallel()

		t.Run("Bool", func(t *testing.T) {
			t.Parallel()

			assert.False(t, c.Bool("data4"))
			assert.False(t, c.BoolDefault("data4", false))
			assert.False(t, c.BoolDefault("data4", true))
			assert.NotPanics(t, func() { c.MustBool("data4") })
		})

		t.Run("Int", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, 0, c.Int("data4"))
			assert.Equal(t, 0, c.IntDefault("data4", 0))
			assert.Equal(t, 0, c.IntDefault("data4", 1))
			assert.NotPanics(t, func() { c.MustInt("data4") })
		})

		t.Run("String", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, "0", c.String("data4"))
			assert.Equal(t, "0", c.StringDefault("data4", ""))
			assert.Equal(t, "0", c.StringDefault("data4", "a string"))
			assert.NotPanics(t, func() { c.MustString("data4") })
		})

		t.Run("Bytes", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, []byte("0"), c.Bytes("data4"))
			assert.Equal(t, []byte("0"), c.BytesDefault("data4", nil))
			assert.Equal(t, []byte("0"), c.BytesDefault("data4", []byte{}))
			assert.Equal(t, []byte("0"), c.BytesDefault("data4", []byte("some bytes")))
			assert.NotPanics(t, func() { c.MustBytes("data4") })
		})
	})
}
