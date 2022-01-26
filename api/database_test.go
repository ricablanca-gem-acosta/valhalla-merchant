package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InitDb_CloseDb(t *testing.T) {
	if assert.Nil(t, InitDb(false)) {
		assert.Nil(t, CloseDb(true))
	}
}

func Test_GetMerchant_CreateMerchant(t *testing.T) {
	if !assert.Nil(t, InitDb(false)) {
		return
	}
	merchant, err := createMerchant(nil)
	assert.Nil(t, err)
	if assert.NotNil(t, merchant) {
		gottenMerchant, err := getMerchant(merchant.Code)
		assert.Nil(t, err)
		assert.Equal(t, merchant.Code, gottenMerchant.Code)
	}
	assert.Nil(t, CloseDb(true))
}

func Test_ListMerchants(t *testing.T) {
	if !assert.Nil(t, InitDb(false)) {
		return
	}
	count := 1000
	for i := 0; i < count; i++ {
		_, err := createMerchant(nil)
		assert.Nil(t, err)
	}
	merchants, err := listMerchants()
	assert.Nil(t, err)
	if assert.NotNil(t, merchants) {
		assert.Equal(t, count, len(merchants))
	}
	assert.Nil(t, CloseDb(true))
}

func Test_DeleteMerchant(t *testing.T) {
	if !assert.Nil(t, InitDb(false)) {
		return
	}
	merchant, err := createMerchant(nil)
	assert.Nil(t, err)
	if assert.NotNil(t, merchant) {
		err := deleteMerchant(merchant.Code)
		if assert.Nil(t, err) {
			merchant, err = getMerchant(merchant.Code)
			assert.NotNil(t, err)
			assert.Nil(t, merchant)
		}
	}
	assert.Nil(t, CloseDb(true))
}

func Test_CreateMember(t *testing.T) {
	if !assert.Nil(t, InitDb(false)) {
		return
	}
	var testEmail = "name@email.com"
	merchant, err := createMerchant(nil)
	assert.Nil(t, err)
	if assert.NotNil(t, merchant) {
		err = createMember(merchant.Code, testEmail)
		assert.Nil(t, err)
		gottenMerchant, err := getMerchant(merchant.Code)
		assert.Nil(t, err)
		if assert.NotNil(t, gottenMerchant) {
			assert.Equal(t, gottenMerchant.Members[0].Email, testEmail)
		}
	}
	assert.Nil(t, CloseDb(true))
}

func Test_DeleteMember(t *testing.T) {
	if !assert.Nil(t, InitDb(false)) {
		return
	}
	var testEmail = "name@email.com"
	merchant, err := createMerchant(nil)
	assert.Nil(t, err)
	if assert.NotNil(t, merchant) {
		err = createMember(merchant.Code, testEmail)
		assert.Nil(t, err)
		gottenMerchant, err := getMerchant(merchant.Code)
		assert.Nil(t, err)
		if assert.NotNil(t, gottenMerchant) {
			assert.Len(t, gottenMerchant.Members, 1)
			err = deleteMember(merchant.Code, testEmail)
			assert.Nil(t, err)
			gottenMerchant, err = getMerchant(merchant.Code)
			assert.Nil(t, err)
			if assert.NotNil(t, gottenMerchant) {
				assert.Len(t, gottenMerchant.Members, 0)
			}
		}
	}
	assert.Nil(t, CloseDb(true))
}
