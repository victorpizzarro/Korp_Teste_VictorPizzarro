package repository

import (
	"Korp_Teste_VictorPizzarro/service-estoque/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type ProdutoDB struct {
	gorm.Model

	Codigo    string `gorm:"uniqueIndex; not null"`
	Descricao string `gorm:"not null"`
	Saldo     int    `gorm:"not null; default:0"`
	Version   int    `gorm:"not null; default:1"`
}

func (ProdutoDB) TableName() string {
	return "produtos"
}

type ProdutoRepository interface {
	SalvarProduto(produto *domain.Produto) error
	BuscarProdutoPorCodigo(codigo string) (*domain.Produto, error)
	AtualizarSaldo(produto *domain.Produto) error
}

type produtoRepositoryPostgres struct {
	db *gorm.DB
}

func NewProdutoRepository(db *gorm.DB) ProdutoRepository {
	return &produtoRepositoryPostgres{db: db}
}

func (repository *produtoRepositoryPostgres) SalvarProduto(produto *domain.Produto) error {
	produtoDB := ProdutoDB{

		Codigo:    produto.Codigo(),
		Descricao: produto.Descricao(),
		Saldo:     produto.Saldo(),
	}

	return repository.db.Create(&produtoDB).Error
}

func (repository *produtoRepositoryPostgres) BuscarProdutoPorCodigo(codigo string) (*domain.Produto, error) {
	var produtoDB ProdutoDB

	err := repository.db.Where("codigo = ?", codigo).First(&produtoDB).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return domain.NewProduto(produtoDB.Codigo, produtoDB.Descricao, produtoDB.Saldo)
}

func (repository *produtoRepositoryPostgres) AtualizarSaldo(produto *domain.Produto) error {
	err := repository.db.Model(&ProdutoDB{}).
		Where("codigo = ?", produto.Codigo()).
		Update("saldo", produto.Saldo()).Error

	return err
}
