package run

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCondition(t *testing.T) {
	name := "dummyStepName"
	valid := map[string]*condition{
		"always()":    {cond: "always", step: ""},
		"any()":       {cond: "any", step: ""},
		"failed()":    {cond: "failed", step: ""},
		"skipped()":   {cond: "skipped", step: ""},
		"succeeded()": {cond: "succeeded", step: ""},

		fmt.Sprintf("always(%s)", name):    {cond: "always", step: name},
		fmt.Sprintf("any(%s)", name):       {cond: "any", step: name},
		fmt.Sprintf("failed(%s)", name):    {cond: "failed", step: name},
		fmt.Sprintf("skipped(%s)", name):   {cond: "skipped", step: name},
		fmt.Sprintf("succeeded(%s)", name): {cond: "succeeded", step: name},
	}

	for k, v := range valid {
		result, err := parseCondition(k)
		assert.NoError(t, err)
		assert.Equal(t, v, result)
	}

	invalid := map[string]*condition{
		"always":   nil,
		"any(":     nil,
		"failed)(": nil,

		fmt.Sprintf("always %s", name): nil,
		fmt.Sprintf("any(%s", name):    nil,
		fmt.Sprintf("failed%s)", name): nil,
	}

	for k := range invalid {
		result, err := parseCondition(k)
		assert.Error(t, err)
		assert.Nil(t, result)
	}
}

func TestShouldExecuteAlways(t *testing.T) {
	sut := statuses{}
	sut.addStatus("somestep", Failed)

	res, err := sut.shouldExecute("always(somestep)")
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = sut.shouldExecute("always(anotherstep)")
	assert.NoError(t, err)
	assert.True(t, res)
}

func TestShouldExecuteAny(t *testing.T) {
	sut := statuses{}
	sut.addStatus("somestep", Failed)

	res, err := sut.shouldExecute("any(somestep)")
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = sut.shouldExecute("any(anotherstep)")
	assert.Error(t, err)
	assert.False(t, res)
}

func TestShouldExecuteFailed(t *testing.T) {
	sut := statuses{}
	sut.addStatus("somestep", Failed)
	sut.addStatus("anotherstep", Succeeded)

	res, err := sut.shouldExecute("failed(somestep)")
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = sut.shouldExecute("failed(anotherstep)")
	assert.NoError(t, err)
	assert.False(t, res)

	res, err = sut.shouldExecute("failed(thirdstep)")
	assert.Error(t, err)
	assert.False(t, res)
}

func TestShouldExecuteSucceeded(t *testing.T) {
	sut := statuses{}
	sut.addStatus("somestep", Succeeded)
	sut.addStatus("anotherstep", Failed)

	res, err := sut.shouldExecute("succeeded(somestep)")
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = sut.shouldExecute("succeeded(anotherstep)")
	assert.NoError(t, err)
	assert.False(t, res)

	res, err = sut.shouldExecute("succeeded(thirdstep)")
	assert.Error(t, err)
	assert.False(t, res)
}
