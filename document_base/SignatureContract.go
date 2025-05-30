package smart_documents

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	d "github.com/rafa-mori/smart_documents/data_structures"
)

type SignatureContract struct {
	contractapi.Contract
}

func (s *SignatureContract) ContractName() string {
	return "SignatureContract"
}

func (s *SignatureContract) SignDocument(ctx contractapi.TransactionContextInterface, id string, signature string) error {
	docBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("erro ao buscar o documento: %v", err)
	}
	if docBytes == nil {
		return fmt.Errorf("documento %s não existe", id)
	}

	var doc d.Document
	if err = json.Unmarshal(docBytes, &doc); err != nil {
		return fmt.Errorf("erro ao deserializar o documento: %v", err)
	}

	// Adiciona a assinatura
	doc.Signatures = append(doc.Signatures, signature)
	updatedDoc, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("erro ao serializar o documento atualizado: %v", err)
	}
	return ctx.GetStub().PutState(id, updatedDoc)
}

func (s *SignatureContract) GetDocumentHistory(ctx contractapi.TransactionContextInterface, id string) ([]string, error) {
	historyIterator, err := ctx.GetStub().GetHistoryForKey(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter histórico do documento: %v", err)
	}
	defer historyIterator.Close()

	history := []string{}
	for historyIterator.HasNext() {
		historyItem, err := historyIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("erro ao iterar histórico: %v", err)
		}
		history = append(history, string(historyItem.Value))
	}
	return history, nil
}

func (s *SignatureContract) GetDocumentState(ctx contractapi.TransactionContextInterface, id string) (*d.Document, error) {
	docBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar estado do documento: %v", err)
	}
	if docBytes == nil {
		return nil, fmt.Errorf("documento %s não existe", id)
	}

	var doc d.Document
	if err = json.Unmarshal(docBytes, &doc); err != nil {
		return nil, fmt.Errorf("erro ao deserializar o documento: %v", err)
	}
	return &doc, nil
}

func (s *SignatureContract) DeleteDocumentState(ctx contractapi.TransactionContextInterface, id string) error {
	return ctx.GetStub().DelState(id)
}
