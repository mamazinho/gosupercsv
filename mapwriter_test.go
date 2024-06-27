package csv_test

import (
	"testing"

	csv "github.com/mamazinho/gosupercsv"
	"github.com/stretchr/testify/assert"
)

func Test_GetData_Success(t *testing.T) {
	mapWriter := csv.NewMapWriter()

	buf, err := mapWriter.GetData()
	assert.NoError(t, err)
	assert.Equal(t, "", buf.String())
}

func Test_WriteHeaders_Success(t *testing.T) {
	var (
		given = []string{"header_1", "header_2"}
		want  = "header_1,header_2\n"
	)
	mapWriter := csv.NewMapWriter()

	err := mapWriter.WriteHeaders(given)
	assert.NoError(t, err)

	buf, err := mapWriter.GetData()
	assert.NoError(t, err)
	assert.Equal(t, want, buf.String())
}

func Test_WriteLine_Success(t *testing.T) {
	t.Run("without header", func(t *testing.T) {
		var (
			given = map[string]string{"header_1": "value_1", "header_2": "value_2"}
			want  = "header_1,header_2\nvalue_1,value_2\n"
		)
		mapWriter := csv.NewMapWriter()

		err := mapWriter.WriteLine(given)
		assert.NoError(t, err)

		buf, err := mapWriter.GetData()
		assert.NoError(t, err)
		assert.Equal(t, want, buf.String())
	})
	t.Run("without header and write two times", func(t *testing.T) {
		var (
			given = map[string]string{"header_1": "value_1", "header_2": "value_2"}
			want  = "header_1,header_2\nvalue_1,value_2\nvalue_1,value_2\n"
		)
		mapWriter := csv.NewMapWriter()

		err := mapWriter.WriteLine(given)
		assert.NoError(t, err)

		err = mapWriter.WriteLine(given)
		assert.NoError(t, err)

		buf, err := mapWriter.GetData()
		assert.NoError(t, err)
		assert.Equal(t, want, buf.String())
	})
	t.Run("with header and testing the sort too", func(t *testing.T) {
		var (
			givenHeaders = []string{"header_1", "header_2", "first_header"}
			given        = map[string]string{"header_1": "value_1", "first_header": "value", "header_2": "value_2"}
			want         = "first_header,header_1,header_2\nvalue,value_1,value_2\n"
		)
		mapWriter := csv.NewMapWriter()

		err := mapWriter.WriteHeaders(givenHeaders)
		assert.NoError(t, err)
		err = mapWriter.WriteLine(given)
		assert.NoError(t, err)

		buf, err := mapWriter.GetData()
		assert.NoError(t, err)
		assert.Equal(t, want, buf.String())
	})
}

func Test_WriteLine_Errors(t *testing.T) {
	t.Run("without value", func(t *testing.T) {
		wantErr := csv.ErrContract
		mapWriter := csv.NewMapWriter()

		err := mapWriter.WriteLine(nil)
		assert.ErrorIs(t, err, wantErr)
	})
	t.Run("has headers but are different from value keys", func(t *testing.T) {
		var (
			wantErr      = csv.ErrContract
			givenHeaders = []string{"header_1", "header_2", "one_more_header"}
			given        = map[string]string{"header_1": "value_1", "header_2": "value_2"}
		)
		mapWriter := csv.NewMapWriter()

		err := mapWriter.WriteHeaders(givenHeaders)
		assert.NoError(t, err)

		err = mapWriter.WriteLine(given)
		assert.ErrorIs(t, err, wantErr)
	})
}

func Test_WriteLines_Success(t *testing.T) {
	t.Run("without header", func(t *testing.T) {
		var (
			given = []map[string]string{
				{"header_1": "value_1", "header_2": "value_2"},
				{"header_1": "value_a", "header_2": "value_b"},
			}
			want = "header_1,header_2\nvalue_1,value_2\nvalue_a,value_b\n"
		)
		mapWriter := csv.NewMapWriter()

		err := mapWriter.WriteLines(given)
		assert.NoError(t, err)

		buf, err := mapWriter.GetData()
		assert.NoError(t, err)
		assert.Equal(t, want, buf.String())
	})
	t.Run("with header and testing the sort too", func(t *testing.T) {
		var (
			givenHeaders = []string{"last_header", "header_1", "header_2", "first_header"}
			given        = []map[string]string{
				{"header_1": "value_1", "last_header": "last_value_2", "first_header": "first_value_1", "header_2": "value_2"},
				{"header_1": "value_a", "first_header": "first_value_a", "last_header": "last_value_b", "header_2": "value_b"},
				{"first_header": "first_value_c", "last_header": "last_value_d", "header_1": "value_c", "header_2": "value_d"},
				{"last_header": "last_value_f", "first_header": "first_value_e", "header_1": "value_e", "header_2": "value_f"},
			}
			want = "first_header,header_1,header_2,last_header\n" +
				"first_value_1,value_1,value_2,last_value_2\n" +
				"first_value_a,value_a,value_b,last_value_b\n" +
				"first_value_c,value_c,value_d,last_value_d\n" +
				"first_value_e,value_e,value_f,last_value_f\n"
		)
		mapWriter := csv.NewMapWriter()

		err := mapWriter.WriteHeaders(givenHeaders)
		assert.NoError(t, err)
		err = mapWriter.WriteLines(given)
		assert.NoError(t, err)

		buf, err := mapWriter.GetData()
		assert.NoError(t, err)
		assert.Equal(t, want, buf.String())
	})
	t.Run("with header and ignore extra keys from value that are not in headers", func(t *testing.T) {
		var (
			givenHeaders = []string{"header_1", "header_2"}
			given        = []map[string]string{
				{"header_1": "value_1", "header_2": "value_2"},
				{"header_1": "value_a", "extra_header": "value", "header_2": "value_b"},
			}
			want = "header_1,header_2\nvalue_1,value_2\nvalue_a,value_b\n"
		)
		mapWriter := csv.NewMapWriter()

		err := mapWriter.WriteHeaders(givenHeaders)
		assert.NoError(t, err)
		err = mapWriter.WriteLines(given)
		assert.NoError(t, err)

		buf, err := mapWriter.GetData()
		assert.NoError(t, err)
		assert.Equal(t, want, buf.String())
	})
	t.Run("with header and ignore extra header that are in first value on the other values", func(t *testing.T) {
		var (
			given = []map[string]string{
				{"header_1": "value_1", "extra_header": "extra_value", "header_2": "value_2"},
				{"header_1": "value_a", "header_2": "value_b"},
			}
			want = "extra_header,header_1,header_2\nextra_value,value_1,value_2\n,value_a,value_b\n"
		)
		mapWriter := csv.NewMapWriter()

		err := mapWriter.WriteLines(given)
		assert.NoError(t, err)

		buf, err := mapWriter.GetData()
		assert.NoError(t, err)
		assert.Equal(t, want, buf.String())
	})
}

func Test_WriteLines_Errors(t *testing.T) {
	t.Run("without value", func(t *testing.T) {
		wantErr := csv.ErrContract
		mapWriter := csv.NewMapWriter()

		err := mapWriter.WriteLines(nil)
		assert.ErrorIs(t, err, wantErr)
	})
	t.Run("has headers but are different from value keys", func(t *testing.T) {
		var (
			wantErr      = csv.ErrContract
			givenHeaders = []string{"header_1", "header_2", "one_more_header"}
			given        = []map[string]string{
				{"header_1": "value_1", "header_2": "value_2"},
			}
		)
		mapWriter := csv.NewMapWriter()

		err := mapWriter.WriteHeaders(givenHeaders)
		assert.NoError(t, err)

		err = mapWriter.WriteLines(given)
		assert.ErrorIs(t, err, wantErr)
	})
}
