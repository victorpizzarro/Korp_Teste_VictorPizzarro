import { Routes } from '@angular/router';
export const routes: Routes = [
  {
    path: '',
    redirectTo: 'produtos/cadastro',
    pathMatch: 'full',
  },
  {
    path: 'produtos/cadastro',
    loadComponent: () =>
      import('./features/produtos/cadastro-produto/cadastro-produto.component')
        .then(m => m.CadastroProdutoComponent),
    title: 'Cadastro de Produtos — Korp',
  },
  {
    path: 'produtos/lista',
    loadComponent: () =>
      import('./features/produtos/lista-produtos/lista-produtos.component')
        .then(m => m.ListaProdutosComponent),
    title: 'Lista de Produtos — Korp',
  },
  {
    path: 'notas-fiscais/cadastro',
    loadComponent: () =>
      import('./features/notas-fiscais/cadastro-nota-fiscal/cadastro-nota-fiscal.component')
        .then(m => m.CadastroNotaFiscalComponent),
    title: 'Cadastro de Nota Fiscal — Korp',
  },
  {
    path: 'notas-fiscais/lista',
    loadComponent: () =>
      import('./features/notas-fiscais/lista-notas-fiscais/lista-notas-fiscais.component')
        .then(m => m.ListaNotasFiscaisComponent),
    title: 'Notas Fiscais — Korp',
  },
];
