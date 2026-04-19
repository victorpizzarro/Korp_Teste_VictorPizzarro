import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Subscription, finalize } from 'rxjs';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { LoadingOverlayComponent } from '../../../shared/components/loading-overlay/loading-overlay.component';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { NotaFiscalService, NotaFiscal } from '../../../core/services/nota-fiscal.service';
import { ConfirmDialogComponent } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
@Component({
  selector: 'app-lista-notas-fiscais',
  standalone: true,
  imports: [
    CommonModule,
    MatTableModule,
    MatButtonModule,
    MatIconModule,
    MatChipsModule,
    MatTooltipModule,
    MatSnackBarModule,
    MatProgressSpinnerModule,
    MatDialogModule,
    LoadingOverlayComponent,
  ],
  template: `
    <!-- Overlay de loading — bloqueia a tela durante a impressão -->
    <app-loading-overlay [isLoading]="isLoading" />
    <div class="page-container">
      <div class="page-header">
        <div>
          <h1 class="page-title">Notas Fiscais</h1>
          <p class="page-subtitle">Gerencie e imprima as notas fiscais cadastradas.</p>
        </div>
        <button mat-mini-fab color="primary" (click)="carregarNotas()" id="btn-refresh-nf">
          <mat-icon>refresh</mat-icon>
        </button>
      </div>
      @if (isLoadingList) {
        <div class="loading-container">
          <mat-spinner diameter="40"></mat-spinner>
        </div>
      } @else if (notas.length === 0) {
        <div class="empty-state card">
          <mat-icon class="empty-icon">receipt_long</mat-icon>
          <h3>Nenhuma nota fiscal encontrada</h3>
          <p>Emita uma nota fiscal para começar.</p>
        </div>
      } @else {
        <div class="card table-card">
          <table mat-table [dataSource]="notas" class="nf-table" id="table-notas-fiscais">
            <!-- Nº Sequencial -->
            <ng-container matColumnDef="numero">
              <th mat-header-cell *matHeaderCellDef>Nº</th>
              <td mat-cell *matCellDef="let nota">
                <span class="numero-badge">#{{ nota.numeroSequencial }}</span>
              </td>
            </ng-container>
            <!-- Status -->
            <ng-container matColumnDef="status">
              <th mat-header-cell *matHeaderCellDef>Status</th>
              <td mat-cell *matCellDef="let nota">
                <span class="status-chip"
                  [class.status-aberta]="nota.status === 'Aberta'"
                  [class.status-fechada]="nota.status === 'Fechada'">
                  <mat-icon class="status-icon">{{ nota.status === 'Aberta' ? 'lock_open' : 'lock' }}</mat-icon>
                  {{ nota.status }}
                </span>
              </td>
            </ng-container>
            <!-- Data Criação -->
            <ng-container matColumnDef="dataCriacao">
              <th mat-header-cell *matHeaderCellDef>Data</th>
              <td mat-cell *matCellDef="let nota">
                {{ nota.dataCriacao | date:'dd/MM/yyyy HH:mm' }}
              </td>
            </ng-container>
            <!-- Itens -->
            <ng-container matColumnDef="itens">
              <th mat-header-cell *matHeaderCellDef>Itens</th>
              <td mat-cell *matCellDef="let nota">
                <span class="items-summary" [matTooltip]="getItensTooltip(nota)">
                  {{ nota.itens.length }} {{ nota.itens.length === 1 ? 'produto' : 'produtos' }}
                </span>
              </td>
            </ng-container>
            <!-- Ação: Imprimir -->
            <ng-container matColumnDef="acoes">
              <th mat-header-cell *matHeaderCellDef>Ações</th>
              <td mat-cell *matCellDef="let nota">
                <!--
                  Botão de Imprimir SOMENTE VISÍVEL quando status === 'Aberta'.
                  Requisito do edital: "Não permitir a impressão de notas com status diferente de Aberta"
                -->
                @if (nota.status === 'Aberta') {
                  <button
                    mat-raised-button
                    color="accent"
                    (click)="imprimirNota(nota)"
                    [disabled]="isLoading"
                    id="btn-imprimir-{{nota.numeroSequencial}}">
                    <mat-icon>print</mat-icon>
                    Imprimir
                  </button>
                } @else {
                  <span class="impressa-label">
                    <mat-icon>check_circle</mat-icon>
                    Impressa
                  </span>
                }
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
    .nf-table {
      width: 100%;
      background: transparent !important;
    }
    .numero-badge {
      font-family: 'JetBrains Mono', monospace;
      font-weight: 700;
      font-size: 15px;
      color: #7c8aff;
    }
    .status-chip {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      padding: 4px 14px;
      border-radius: 20px;
      font-size: 13px;
      font-weight: 600;
    }
    .status-icon {
      font-size: 16px;
      width: 16px;
      height: 16px;
    }
    .status-aberta {
      background: rgba(46, 160, 67, 0.15);
      color: #3fb950;
    }
    .status-fechada {
      background: rgba(139, 148, 158, 0.15);
      color: #8b949e;
    }
    .items-summary {
      color: var(--korp-text-secondary);
      cursor: help;
      border-bottom: 1px dashed var(--korp-text-secondary);
    }
    .impressa-label {
      display: inline-flex;
      align-items: center;
      gap: 4px;
      color: var(--korp-text-secondary);
      font-size: 13px;
      mat-icon {
        font-size: 18px;
        width: 18px;
        height: 18px;
        color: #3fb950;
      }
    }
    button[color="accent"] {
      display: flex;
      align-items: center;
      gap: 6px;
      font-weight: 600;
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
export class ListaNotasFiscaisComponent implements OnInit, OnDestroy {
  notas: NotaFiscal[] = [];
  isLoading = false;       
  isLoadingList = false;    
  displayedColumns = ['numero', 'status', 'dataCriacao', 'itens', 'acoes'];
  private subscriptions = new Subscription();
  constructor(
    private notaFiscalService: NotaFiscalService,
    private snackBar: MatSnackBar,
    private dialog: MatDialog
  ) {}
  ngOnInit(): void {
    this.carregarNotas();
  }
  ngOnDestroy(): void {
    this.subscriptions.unsubscribe();
  }
  carregarNotas(): void {
    this.isLoadingList = true;
    const sub = this.notaFiscalService.listar().subscribe({
      next: (dados) => {
        this.notas = dados;
        this.isLoadingList = false;
      },
      error: () => {
        this.isLoadingList = false;
      },
    });
    this.subscriptions.add(sub);
  }
  imprimirNota(nota: NotaFiscal): void {
    this.isLoading = true;

    // First check for anomalies via Gemini AI
    const subAnomalia = this.notaFiscalService.analisarAnomalia(nota.numeroSequencial)
      .subscribe({
        next: (resp) => {
          if (resp.tem_anomalia) {
            this.isLoading = false;
            // Ask for confirmation
            const dialogRef = this.dialog.open(ConfirmDialogComponent, {
              width: '400px',
              data: {
                title: 'Possível Erro de Faturamento',
                message: resp.mensagem || 'A IA detectou uma possível anomalia nos itens desta nota.'
              }
            });

            dialogRef.afterClosed().subscribe(result => {
              if (result === true) {
                this.executarImpressao(nota);
              }
            });
          } else {
            // No anomaly, just print
            this.executarImpressao(nota);
          }
        },
        error: (err) => {
          this.isLoading = false;
          this.snackBar.open(
            'Não foi possível analisar anomalias. Verifique o servidor.',
            'OK',
            { duration: 5000, horizontalPosition: 'end', panelClass: ['snack-error'] }
          );
          console.error(err);
        }
      });
    
    this.subscriptions.add(subAnomalia);
  }

  private executarImpressao(nota: NotaFiscal): void {
    this.isLoading = true;
    const sub = this.notaFiscalService
      .imprimir(nota.numeroSequencial)
      .pipe(
        finalize(() => {
          this.isLoading = false;
        })
      )
      .subscribe({
        next: (notaAtualizada) => {
          const index = this.notas.findIndex(n => n.numeroSequencial === nota.numeroSequencial);
          if (index >= 0) {
            this.notas[index] = notaAtualizada;
            this.notas = [...this.notas];
          }
          this.snackBar.open(
            `Nota #${nota.numeroSequencial} impressa com sucesso! Estoque atualizado.`,
            'OK',
            { duration: 5000, horizontalPosition: 'end', verticalPosition: 'top', panelClass: ['snack-success'] }
          );
        },
        error: () => {
        },
      });
    this.subscriptions.add(sub);
  }
  getItensTooltip(nota: NotaFiscal): string {
    return nota.itens
      .map(item => `${item.codigoProduto}: ${item.quantidade} un.`)
      .join(' | ');
  }
}
