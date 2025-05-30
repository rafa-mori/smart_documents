package smart_documents

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	d "github.com/rafa-mori/smart_documents/data_structures"
)

type ApprovalContract struct {
	contractapi.Contract
}

func (a *ApprovalContract) ContractName() string {
	return "ApprovalContract"
}

func (a *ApprovalContract) RegisterDocument(ctx contractapi.TransactionContextInterface, id string, content string) error {
	// Verifica se o documento já existe
	existing, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("erro ao verificar existência do documento: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("o documento %s já está registrado", id)
	}
	doc := d.Document{
		ID:         id,
		Content:    content,
		Signatures: []string{},
		Approved:   false,
	}
	docBytes, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("erro ao serializar o documento: %v", err)
	}
	return ctx.GetStub().PutState(id, docBytes)
}

func (a *ApprovalContract) ApproveDocument(ctx contractapi.TransactionContextInterface, id string) error {
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

	// Garante que o documento tenha sido assinado antes de aprovar
	if len(doc.Signatures) == 0 {
		return fmt.Errorf("o documento deve ter pelo menos uma assinatura para ser aprovado")
	}

	doc.Approved = true
	updatedBytes, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("erro ao serializar o documento aprovado: %v", err)
	}
	return ctx.GetStub().PutState(id, updatedBytes)
}

func (a *ApprovalContract) GetDocumentHistory(ctx contractapi.TransactionContextInterface, id string) ([]string, error) {
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

func (a *ApprovalContract) GetDocumentState(ctx contractapi.TransactionContextInterface, id string) (*d.Document, error) {
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

func (a *ApprovalContract) DeleteDocumentState(ctx contractapi.TransactionContextInterface, id string) error {
	return ctx.GetStub().DelState(id)
}
