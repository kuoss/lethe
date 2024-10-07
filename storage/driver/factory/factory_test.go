package factory

import (
	"fmt"
	"testing"

	"github.com/kuoss/common/tester"
	"github.com/kuoss/lethe/storage/driver"
	"github.com/kuoss/lethe/storage/driver/factory/fake"
	"github.com/stretchr/testify/assert"
)

// fakeDriverFactory here: to avoid imports cycle
type fakeDriverFactory struct{}

func (factory *fakeDriverFactory) Create(parameters map[string]interface{}) (driver.Driver, error) {
	return fake.New(), nil
}

var (
	fakeDriverFactory1 = &fakeDriverFactory{}
)

func unregisterAll() {
	driverFactories = map[string]StorageDriverFactory{}
}

func TestRegister(t *testing.T) {
	assert.NotNil(t, fakeDriverFactory1)
	testCases := []struct {
		name      string
		factory   StorageDriverFactory
		wantError string
		wantNames []string
	}{
		{
			"", nil,
			"factory is nil",
			[]string{},
		},
		{
			"", fakeDriverFactory1,
			"",
			[]string{""},
		},
		{
			"fake1", fakeDriverFactory1,
			"",
			[]string{"fake1"},
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i), func(t *testing.T) {
			unregisterAll()
			err := Register(tc.name, tc.factory)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			names := []string{}
			for k := range driverFactories {
				names = append(names, k)
			}
			assert.Equal(t, tc.wantNames, names)
		})
	}
}

func TestRegister_NameDuplicated(t *testing.T) {
	unregisterAll()
	err := Register("fake1", fakeDriverFactory1)
	assert.NoError(t, err)
	err = Register("fake1", fakeDriverFactory1)
	assert.EqualError(t, err, "factory name duplicated: fake1")
}

func TestGet(t *testing.T) {
	unregisterAll()
	err := Register("fake1", fakeDriverFactory1)
	assert.NoError(t, err)

	testCases := []struct {
		name       string
		parameters map[string]interface{}
		want       string
		wantError  string
	}{
		{
			"", map[string]interface{}{},
			"<nil>",
			"invalid driver name: ",
		},
		{
			"fake1", map[string]interface{}{},
			"&fake.driver{}",
			"",
		},
		{
			"fake2", map[string]interface{}{},
			"<nil>",
			"invalid driver name: fake2",
		},
	}
	for i, tc := range testCases {
		t.Run(tester.CaseName(i), func(t *testing.T) {
			driver, err := Get(tc.name, tc.parameters)
			got := fmt.Sprintf("%#v", driver)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
