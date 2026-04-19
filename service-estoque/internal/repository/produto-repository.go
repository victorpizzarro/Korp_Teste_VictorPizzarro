package repository

import (
	"Korp_Teste_VictorPizzarro/service-estoque/internal/domain"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

type ItemDebito struct {
	Codigo     string
	Quantidade int
}

type ProdutoRepository interface {
	SalvarProduto(produto *domain.Produto) error
	BuscarProdutoPorCodigo(codigo string) (*domain.Produto, error)
	AtualizarSaldo(produto *domain.Produto) error
	ListarTodos() ([]*domain.Produto, error)

	DebitarSaldo(codigo string, quantidade int) error
	DebitarSaldoLote(itens []ItemDebito) error
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

func (repository *produtoRepositoryPostgres) ListarTodos() ([]*domain.Produto, error) {
	var produtosDB []ProdutoDB

	if err := repository.db.Find(&produtosDB).Error; err != nil {
		return nil, err
	}

	var produtos []*domain.Produto
	for _, p := range produtosDB {
		produto, err := domain.NewProduto(p.Codigo, p.Descricao, p.Saldo)
		if err != nil {
			return nil, err
		}
		produtos = append(produtos, produto)
	}

	return produtos, nil
}

func (repository *produtoRepositoryPostgres) AtualizarSaldo(produto *domain.Produto) error {
	err := repository.db.Model(&ProdutoDB{}).
		Where("codigo = ?", produto.Codigo()).
		Update("saldo", produto.Saldo()).Error

	return err
}

func (repository *produtoRepositoryPostgres) DebitarSaldo(codigo string, quantidade int) error {
	return repository.db.Transaction(func(tx *gorm.DB) error {
		var produtoDB ProdutoDB

		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("codigo = ?", codigo).
			First(&produtoDB).Error

		if err != nil {
			return err
		}

		if produtoDB.Saldo < quantidade {
			return errors.New("saldo insuficiente para realizar o débito")
		}

		produtoDB.Saldo -= quantidade

		return tx.Save(&produtoDB).Error
	})
}

func (repository *produtoRepositoryPostgres) DebitarSaldoLote(itens []ItemDebito) error {

	return repository.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range itens {
			var produtoDB ProdutoDB

			err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("codigo = ?", item.Codigo).
				First(&produtoDB).Error

			if err != nil {
				return err
			}

			if produtoDB.Saldo < item.Quantidade {
				return errors.New("saldo insuficiente para o produto: " + item.Codigo)
			}

			produtoDB.Saldo -= item.Quantidade

			if err := tx.Save(&produtoDB).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
