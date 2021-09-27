package wallet

import (
	"fmt"
	"testing"
	"reflect"

	"github.com/Behzod01/wallet/pkg/types"
)

func TestService_FindAccountByID_success(t *testing.T) {
	service := &Service{
		accounts: []*types.Account{
			{
				ID: 1111,
				Phone: "90-000-00-01",
				Balance: 10_000,
			},
			{
				ID: 4444,
				Phone: "90-444-00-01",
				Balance: 1500_00,
			},
			{
				ID: 2222,
				Phone: "93-000-33-33",
				Balance: 50_000,
			},
			{
				ID: 5555,
				Phone: "77-700-07-77",
				Balance: 100_000,
			},
			{
				ID: 3333,
				Phone: "91-111-11-01",
				Balance: 160_000,
			},
		},
	}
	expected := types.Account{
		ID: 5555,
		Phone: "77-700-07-77",
		Balance: 100_000,
	}
	
	result, err := service.FindAccountByID(5555)
	
	if err == nil {
		fmt.Println(err)
		return
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("invalid result, expected: %v, actual: %v", expected, result)	
	}
}

func TestService_FindAccountByID_notFound(t *testing.T) {
	service := &Service{
		accounts: []*types.Account{
			{
				ID: 1111,
				Phone: "90-000-00-01",
				Balance: 10_000,
			},
			{
				ID: 4444,
				Phone: "90-444-00-01",
				Balance: 1500_00,
			},
			{
				ID: 2222,
				Phone: "93-000-33-33",
				Balance: 50_000,
			},
			{
				ID: 5555,
				Phone: "77-700-07-77",
				Balance: 100_000,
			},
			{
				ID: 3333,
				Phone: "91-111-11-01",
				Balance: 160_000,
			},
		},
	}
	expected := ErrAccountNotFound
	
	result, err := service.FindAccountByID(55555)
	
	if err != nil {
		fmt.Println(err)
		return
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("invalid result, expected: %v, actual: %v", expected, result)	
	}
}