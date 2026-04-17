package repository

import (
	"Korp_Teste_VictorPizzarro/service-faturamento/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type NotaFiscalDB struct {
	gorm.Model

	NumeroSequencial int                `gorm:"column:numero_sequencial;autoIncrement;uniqueIndex;not null"`
	Status           string             `gorm:"not null"`
	Itens            []ItemNotaFiscalDB `gorm:"foreignKey:NotaFiscalID"`
}

func (NotaFiscalDB) TableName() string {
	return "notas_fiscais"
}

type ItemNotaFiscalDB struct {
	gorm.Model

	NotaFiscalID  uint   `gorm:"not null"`
	CodigoProduto string `gorm:"not null"`
	Quantidade    int    `gorm:"not null"`
}

func (ItemNotaFiscalDB) TableName() string {
	return "itens_nota_fiscal"
}

type NotaFiscalRepository interface {
	SalvarNotaFiscal(nota *domain.NotaFiscal) error
	BuscarNotaFiscalPorNumero(numero int) (*domain.NotaFiscal, error)
	AtualizarStatus(nota *domain.NotaFiscal) error
}

type notaFiscalRepositoryPostgres struct {
	db *gorm.DB
}

func NewNotaFiscalRepository(db *gorm.DB) NotaFiscalRepository {
	return &notaFiscalRepositoryPostgres{db: db}
}

func (repository *notaFiscalRepositoryPostgres) SalvarNotaFiscal(nota *domain.NotaFiscal) error {

	var itensDB []ItemNotaFiscalDB
	for _, item := range nota.Itens() {
		itensDB = append(itensDB, ItemNotaFiscalDB{
			CodigoProduto: item.CodigoProduto(),
			Quantidade:    item.Quantidade(),
		})
	}

	notaDB := NotaFiscalDB{
		NumeroSequencial: 0,
		Status:           string(nota.Status()),
		Itens:            itensDB,
	}

	err := repository.db.Create(&notaDB).Error
	if err != nil {
		return err
	}

	nota.DefinirNumeroGerado(notaDB.NumeroSequencial)

	return nil
}

func (repository *notaFiscalRepositoryPostgres) BuscarNotaFiscalPorNumero(numero int) (*domain.NotaFiscal, error) {
	var notaDB NotaFiscalDB

	err := repository.db.Preload("Itens").Where("numero_sequencial = ?", numero).First(&notaDB).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var itensVO []domain.ItemNotaFiscal
	for _, itemDB := range notaDB.Itens {
		itemVO, err := domain.NewItemNotaFiscal(itemDB.CodigoProduto, itemDB.Quantidade)
		if err != nil {
			return nil, err
		}
		itensVO = append(itensVO, *itemVO)
	}

	nota, err := domain.ReconstituirNotaFiscal(
		notaDB.NumeroSequencial,
		domain.StatusNotaFiscal(notaDB.Status),
		notaDB.CreatedAt,
		itensVO,
	)

	if err != nil {
		return nil, err
	}

	return nota, nil
}

func (repository *notaFiscalRepositoryPostgres) AtualizarStatus(nota *domain.NotaFiscal) error {
	err := repository.db.Model(&NotaFiscalDB{}).
		Where("numero_sequencial = ?", nota.NumeroSequencial()).
		Update("status", string(nota.Status())).Error

	return err
}
