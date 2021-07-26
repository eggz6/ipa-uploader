package plist

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/suite"
	"howett.net/plist"
)

type ReadTestSuite struct {
	suite.Suite
	mockErr error
}

func (s *ReadTestSuite) SetupTest() {
	s.mockErr = errors.New("mock error")
}

func TestReadTestSuite(t *testing.T) {
	suite.Run(t, new(ReadTestSuite))
}

func (s *ReadTestSuite) TestReadFromFile() {
	cases := []struct {
		name  string
		err   error
		path  string
		want  PList
		patch func() *gomonkey.Patches
	}{
		{
			name: "success",
			want: PList{"id": "id"},
			patch: func() *gomonkey.Patches {
				res := gomonkey.NewPatches()
				res.ApplyFunc(os.Open, func(_ string) (*os.File, error) { return &os.File{}, nil })

				res.ApplyFunc(ioutil.ReadAll, func(_ io.Reader) ([]byte, error) { return []byte{}, nil })

				res.ApplyFunc(plist.Unmarshal, func(_ []byte, ref interface{}) (int, error) {
					pl := make(PList)
					pl["id"] = "id"
					val := reflect.ValueOf(ref)
					val.Elem().Set(reflect.ValueOf(pl))
					return 0, nil
				})

				return res
			},
		},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			if c.patch != nil {
				p := c.patch()

				defer p.Reset()
			}

			pl, err := ReadFromFile(c.path)

			s.Equal(c.err, err)
			s.Equal(c.want, pl)
		})
	}
}
