package main

import (
    "encoding/json"
    "fmt"
    "log"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
    contractapi.Contract
}

type Account struct {
    DealerID    string  `json:"dealerID"`
    MSISDN      string  `json:"msisdn"`
    MPIN        string  `json:"mpin"`
    Balance     float64 `json:"balance"`
    Status      string  `json:"status"`
    TransAmount float64 `json:"transAmount"`
    TransType   string  `json:"transType"`
    Remarks     string  `json:"remarks"`
}

// CreateAccount creates a new account
func (s *SmartContract) CreateAccount(ctx contractapi.TransactionContextInterface, dealerID string, msisdn string, mpin string, balance float64, status string) error {
    account := Account{
        DealerID:   dealerID,
        MSISDN:     msisdn,
        MPIN:       mpin,
        Balance:    balance,
        Status:     status,
    }

    accountJSON, err := json.Marshal(account)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(dealerID, accountJSON)
}

// ReadAccount retrieves an account
func (s *SmartContract) ReadAccount(ctx contractapi.TransactionContextInterface, dealerID string) (*Account, error) {
    accountJSON, err := ctx.GetStub().GetState(dealerID)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if accountJSON == nil {
        return nil, fmt.Errorf("the account %s does not exist", dealerID)
    }

    var account Account
    err = json.Unmarshal(accountJSON, &account)
    if err != nil {
        return nil, err
    }

    return &account, nil
}

// UpdateAccount updates balance and status of an account
func (s *SmartContract) UpdateAccount(ctx contractapi.TransactionContextInterface, dealerID string, newBalance float64, newStatus string) error {
    account, err := s.ReadAccount(ctx, dealerID)
    if err != nil {
        return err
    }

    account.Balance = newBalance
    account.Status = newStatus

    accountJSON, err := json.Marshal(account)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(dealerID, accountJSON)
}

func main() {
    chaincode, err := contractapi.NewChaincode(new(SmartContract))
    if err != nil {
        log.Panicf("Error creating chaincode: %v", err)
    }

    if err := chaincode.Start(); err != nil {
        log.Panicf("Error starting chaincode: %v", err)
    }
}
