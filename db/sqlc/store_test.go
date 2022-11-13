package db

import "testing"

func TestTransferTx(t *testing){
	//NewStore require a db object so
	//Make a Db object global
	 store := NewStore(testDb)
	 

	 account1 := createRandomTransfer()
}