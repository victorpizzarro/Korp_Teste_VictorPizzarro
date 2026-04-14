package domain

type Produto struct {
	codigo    Codigo
	descricao Descricao
	saldo     Saldo
}

func NewProduto(cod string, desc string, qtdSaldo int) (*Produto, error) {
	codigoVO, err := NewCodigo(cod)

	if err != nil {
		return nil, err
	}

	descricaoVO, err := NewDescricao(desc)

	if err != nil {
		return nil, err
	}

	saldoVO, err := NewSaldo(qtdSaldo)

	if err != nil {
		return nil, err
	}

	return &Produto{
		codigo:    codigoVO,
		descricao: descricaoVO,
		saldo:     saldoVO,
	}, nil

}

func (p *Produto) DiminuirEstoque(qtd int) error {
	novoSaldo, err := p.saldo.Subtrair(qtd)

	if err != nil {
		return err
	}

	p.saldo = novoSaldo

	return nil
}

func (p *Produto) Codigo() string    { return p.codigo.Valor() }
func (p *Produto) Descricao() string { return p.descricao.Valor() }
func (p *Produto) Saldo() int        { return p.saldo.Valor() }
