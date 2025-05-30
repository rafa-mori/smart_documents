package smart_documents

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	d "github.com/rafa-mori/smart_documents/data_structures"
)

type TrafficContract struct {
	contractapi.Contract
}

func (t *TrafficContract) ContractName() string {
	return "TrafficContract"
}

func (t *TrafficContract) RegisterTrafficDocument(ctx contractapi.TransactionContextInterface, id string, content string) error {
	existing, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("erro ao verificar existência do documento: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("o documento %s já está registrado", id)
	}
	doc := d.Document{
		ID:      id,
		Content: content,
	}
	docBytes, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("erro ao serializar o documento: %v", err)
	}
	return ctx.GetStub().PutState(id, docBytes)
}

func (t *TrafficContract) GetDocumentHistory(ctx contractapi.TransactionContextInterface, id string) ([]string, error) {
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

func (t *TrafficContract) GetDocumentState(ctx contractapi.TransactionContextInterface, id string) (*d.Document, error) {
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

func (t *TrafficContract) DeleteDocumentState(ctx contractapi.TransactionContextInterface, id string) error {
	return ctx.GetStub().DelState(id)
}
