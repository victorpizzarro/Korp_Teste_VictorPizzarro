import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Subscription } from 'rxjs';
import { MatTableModule } from '@angular/material/table';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { ProdutoService, Produto } from '../../../core/services/produto.service';
@Component({
  selector: 'app-lista-produtos',
  standalone: true,
  imports: [
    CommonModule,
    MatTableModule,
    MatProgressSpinnerModule,
    MatIconModule,
    MatButtonModule,
  ],
  template: `
    <div class="page-container">
      <div class="page-header">
        <div>
          <h1 class="page-title">Saldo de Estoque</h1>
          <p class="page-subtitle">Visualize todos os produtos cadastrados e seus saldos atuais.</p>
        </div>
        <button mat-mini-fab color="primary" (click)="carregarProdutos()" id="btn-refresh-produtos">
          <mat-icon>refresh</mat-icon>
        </button>
      </div>
      @if (isLoading) {
        <div class="loading-container">
          <mat-spinner diameter="40"></mat-spinner>
        </div>
      } @else if (produtos.length === 0) {
        <div class="empty-state card">
          <mat-icon class="empty-icon">inventory_2</mat-icon>
          <h3>Nenhum produto cadastrado</h3>
          <p>Cadastre seu primeiro produto para começar.</p>
        </div>
      } @else {
        <div class="card table-card">
          <table mat-table [dataSource]="produtos" class="produtos-table" id="table-produtos">
            <ng-container matColumnDef="codigo">
              <th mat-header-cell *matHeaderCellDef>Código</th>
              <td mat-cell *matCellDef="let produto">
                <span class="codigo-badge">{{ produto.codigo }}</span>
              </td>
            </ng-container>
            <ng-container matColumnDef="descricao">
              <th mat-header-cell *matHeaderCellDef>Descrição</th>
              <td mat-cell *matCellDef="let produto">{{ produto.descricao }}</td>
            </ng-container>
            <ng-container matColumnDef="saldo">
              <th mat-header-cell *matHeaderCellDef>Saldo</th>
              <td mat-cell *matCellDef="let produto">
                <span class="saldo-chip" [class.saldo-baixo]="produto.saldo <= 5" [class.saldo-zero]="produto.saldo === 0">
                  {{ produto.saldo }} un.
                </span>
              </td>
            </ng-container>
            <tr mat-header-row *matHeaderRowDef="displayedColumns"></tr>
            <tr mat-row *matRowDef="let row; columns: displayedColumns;" class="table-row"></tr>
          </table>
        </div>
      }
    </div>
  `,
  styles: [`
    .page-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 24px;
    }
    .loading-container {
      display: flex;
      justify-content: center;
      padding: 64px 0;
    }
    .empty-state {
      text-align: center;
      padding: 64px 24px !important;
      h3 { margin-top: 16px; font-size: 18px; }
      p { color: var(--korp-text-secondary); margin-top: 8px; }
    }
    .empty-icon {
      font-size: 48px;
      width: 48px;
      height: 48px;
      color: var(--korp-text-secondary);
    }
    .table-card {
      padding: 0 !important;
      overflow: hidden;
    }
    .produtos-table {
      width: 100%;
      background: transparent !important;
    }
    .codigo-badge {
      background: rgba(79, 107, 237, 0.15);
      color: #7c8aff;
      padding: 4px 12px;
      border-radius: 6px;
      font-family: 'JetBrains Mono', monospace;
      font-size: 13px;
      font-weight: 500;
    }
    .saldo-chip {
      background: rgba(46, 160, 67, 0.15);
      color: #3fb950;
      padding: 4px 12px;
      border-radius: 6px;
      font-weight: 600;
      font-size: 13px;
    }
    .saldo-baixo {
      background: rgba(227, 179, 65, 0.15);
      color: #e3b341;
    }
    .saldo-zero {
      background: rgba(248, 81, 73, 0.15);
      color: #f85149;
    }
    .table-row:hover {
      background: rgba(79, 107, 237, 0.06) !important;
    }
    th.mat-mdc-header-cell {
      color: var(--korp-text-secondary) !important;
      font-size: 12px;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.8px;
      border-bottom-color: var(--korp-border) !important;
    }
    td.mat-mdc-cell {
      border-bottom-color: var(--korp-border) !important;
    }
  `]
})
export class ListaProdutosComponent implements OnInit, OnDestroy {
  produtos: Produto[] = [];
  isLoading = false;
  displayedColumns = ['codigo', 'descricao', 'saldo'];
  private subscription?: Subscription;
  constructor(private produtoService: ProdutoService) {}
  ngOnInit(): void {
    this.carregarProdutos();
  }
  ngOnDestroy(): void {
    this.subscription?.unsubscribe();
  }
  carregarProdutos(): void {
    this.isLoading = true;
    this.subscription = this.produtoService.listar().subscribe({
      next: (dados) => {
        this.produtos = dados;
        this.isLoading = false;
      },
      error: () => {
        this.isLoading = false;
      },
    });
  }
}
