package main

import (
	"log"
	"os"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	d "github.com/rafa-mori/smart_documents/document_base"
)

func ContractsWrapper[T d.SignatureContract | d.ApprovalContract | d.TrafficContract]() *T {
	// Here fits any logic we want to apply to the contract instantiation.
	// From securing the contract to setting up initial values.
	return new(T)
}

type DocumentChaincode map[string]contractapi.ContractInterface

func newDocumentChaincode() DocumentChaincode {
	return DocumentChaincode{
		"SignatureContract": ContractsWrapper[d.SignatureContract](),
		"ApprovalContract":  ContractsWrapper[d.ApprovalContract](),
		"TrafficContract":   ContractsWrapper[d.TrafficContract](),
	}
}

func getNameFromArgs(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

func (cc DocumentChaincode) GetSnapshot(args ...string) {
	// This method is just a placeholder to show how we could implement some logic,
	// like fetching a specific contract by name and sending its json representation
	// to a webhook or logging it. We could monitor the chaincode's state and when a specific
	// contract is requested, we could eventually block it before it is loaded
	// into the chaincode, or log its state for debugging purposes.
	//
	//if len(args) != 0 {
	//	if contract, exists := cc[args[0]]; exists {
	//		contractJSON, _ := json.MarshalIndent(contract, "", "  ")
	//	}
	//}
}
func (cc DocumentChaincode) GetContracts() []contractapi.ContractInterface {
	contracts := make([]contractapi.ContractInterface, 0, len(cc))
	for _, contract := range cc {
		contracts = append(contracts, contract)
	}
	return contracts
}

func main() {
	contracts := newDocumentChaincode()
	if _, ok := contracts[getNameFromArgs(os.Args[1:])]; ok {
		contracts.GetSnapshot(os.Args[1:]...)
	} else {
		chaincode, err := contractapi.NewChaincode(newDocumentChaincode().GetContracts()...)
		if err != nil {
			log.Panicf("Erro ao criar o chaincode: %v", err)
		}
		if err = chaincode.Start(); err != nil {
			log.Panicf("Erro ao iniciar o chaincode: %v", err)
		}
	}
}
