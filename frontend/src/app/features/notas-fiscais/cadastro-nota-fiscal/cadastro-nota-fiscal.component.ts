import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  ReactiveFormsModule,
  FormGroup,
  FormArray,
  FormControl,
  Validators
} from '@angular/forms';
import { Subscription } from 'rxjs';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatDividerModule } from '@angular/material/divider';
import { ProdutoService, Produto } from '../../../core/services/produto.service';
import { NotaFiscalService } from '../../../core/services/nota-fiscal.service';
import { MatOptionModule } from '@angular/material/core';
@Component({
  selector: 'app-cadastro-nota-fiscal',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatOptionModule,
    MatButtonModule,
    MatIconModule,
    MatSnackBarModule,
    MatProgressSpinnerModule,
    MatDividerModule,
  ],
  template: `
    <div class="page-container">
      <h1 class="page-title">Emitir Nota Fiscal</h1>
      <p class="page-subtitle">Crie uma nova nota fiscal adicionando produtos e suas quantidades. A nota será criada com status <strong>Aberta</strong>.</p>
      <div class="card">
        <form [formGroup]="notaForm" (ngSubmit)="onSubmit()" id="form-cadastro-nf">
          <!-- Header dos itens -->
          <div class="items-header">
            <h3>
              <mat-icon>shopping_cart</mat-icon>
              Itens da Nota Fiscal
            </h3>
            <button
              mat-stroked-button
              color="primary"
              type="button"
              (click)="adicionarItem()"
              id="btn-adicionar-item">
              <mat-icon>add</mat-icon>
              Adicionar Item
            </button>
          </div>
          <mat-divider></mat-divider>
          <!-- FormArray — Lista dinâmica de itens -->
          <div formArrayName="itens" class="items-list">
            @for (item of itensFormArray.controls; track $index) {
              <div class="item-row" [formGroupName]="$index">
                <span class="item-number">{{ $index + 1 }}</span>
                <mat-form-field appearance="outline" class="field-produto">
                  <mat-label>Produto</mat-label>
                  <mat-select formControlName="codigoProduto" id="select-produto-{{$index}}">
                    <mat-option *ngFor="let produto of produtosDisponiveis" [value]="produto.codigo">
                      {{ produto.codigo }} — {{ produto.descricao }} (saldo: {{ produto.saldo }})
                    </mat-option>
                  </mat-select>
                  @if (item.get('codigoProduto')?.hasError('required') && item.get('codigoProduto')?.touched) {
                    <mat-error>Selecione um produto</mat-error>
                  }
                </mat-form-field>
                <mat-form-field appearance="outline" class="field-qtd">
                  <mat-label>Quantidade</mat-label>
                  <input matInput type="number" formControlName="quantidade" id="input-quantidade-{{$index}}">
                  @if (item.get('quantidade')?.hasError('required') && item.get('quantidade')?.touched) {
                    <mat-error>Obrigatório</mat-error>
                  }
                  @if (item.get('quantidade')?.hasError('min') && item.get('quantidade')?.touched) {
                    <mat-error>Mín: 1</mat-error>
                  }
                </mat-form-field>
                <button
                  mat-icon-button
                  color="warn"
                  type="button"
                  (click)="removerItem($index)"
                  [disabled]="itensFormArray.length <= 1"
                  id="btn-remover-item-{{$index}}">
                  <mat-icon>delete</mat-icon>
                </button>
              </div>
            }
          </div>
          @if (itensFormArray.length === 0) {
            <div class="empty-items">
              <mat-icon>info</mat-icon>
              <span>Adicione pelo menos um item para criar a nota fiscal.</span>
            </div>
          }
          <div class="form-actions">
            <button
              mat-raised-button
              color="primary"
              type="submit"
              [disabled]="notaForm.invalid || itensFormArray.length === 0 || isLoading"
              id="btn-emitir-nf">
              <span class="button-content">
                @if (isLoading) {
                  <mat-spinner diameter="20" class="btn-spinner"></mat-spinner>
                } @else {
                  <mat-icon>receipt</mat-icon>
                }
                <span>Emitir Nota Fiscal</span>
              </span>
            </button>
          </div>
        </form>
      </div>
    </div>
  `,
  styles: [`
    .items-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;
      h3 {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 16px;
        font-weight: 600;
        color: var(--korp-text-primary);
        mat-icon { color: #7c8aff; }
      }
    }
    .items-list {
      padding: 16px 0;
    }
    .item-row {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 8px;
      padding: 8px 12px;
      border-radius: 8px;
      background: var(--korp-bg-elevated);
      transition: background 0.2s ease;
      &:hover { background: rgba(79, 107, 237, 0.08); }
    }
    .item-number {
      width: 28px;
      height: 28px;
      border-radius: 50%;
      background: rgba(79, 107, 237, 0.2);
      color: #7c8aff;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 12px;
      font-weight: 700;
      flex-shrink: 0;
    }
    .field-produto {
      flex: 2;
    }
    .field-qtd {
      flex: 0.5;
      min-width: 120px;
    }
    .empty-items {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 24px;
      color: var(--korp-text-secondary);
      font-size: 14px;
      justify-content: center;
    }
    .form-actions {
      margin-top: 24px;
      display: flex;
      justify-content: flex-end;
      button {
        padding: 0 24px;
        height: 44px;
        font-weight: 600;
      }
      .button-content {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        width: 100%;
      }
    }
    .btn-spinner {
      display: inline-block;
    }
  `]
})
export class CadastroNotaFiscalComponent implements OnInit, OnDestroy {
  notaForm!: FormGroup;
  produtosDisponiveis: Produto[] = [];
  isLoading = false;
  private subscriptions = new Subscription();
  constructor(
    private notaFiscalService: NotaFiscalService,
    private produtoService: ProdutoService,
    private snackBar: MatSnackBar
  ) {}
  ngOnInit(): void {
    this.notaForm = new FormGroup({
      itens: new FormArray([])
    });
    this.adicionarItem();
    const sub = this.produtoService.listar().subscribe({
      next: (produtos) => this.produtosDisponiveis = produtos,
      error: () => {
        this.snackBar.open('Não foi possível carregar os produtos', 'Fechar', {
          duration: 4000,
          panelClass: ['snack-error'],
        });
      },
    });
    this.subscriptions.add(sub);
  }
  ngOnDestroy(): void {
    this.subscriptions.unsubscribe();
  }
  get itensFormArray(): FormArray {
    return this.notaForm.get('itens') as FormArray;
  }
  adicionarItem(): void {
    const itemGroup = new FormGroup({
      codigoProduto: new FormControl('', [Validators.required]),
      quantidade:    new FormControl(1,  [Validators.required, Validators.min(1)]),
    });
    this.itensFormArray.push(itemGroup);
  }
  removerItem(index: number): void {
    this.itensFormArray.removeAt(index);
  }
  onSubmit(): void {
    if (this.notaForm.invalid || this.itensFormArray.length === 0) return;
    this.isLoading = true;
    const payload = { itens: this.itensFormArray.value };
    const sub = this.notaFiscalService.cadastrar(payload).subscribe({
      next: (res) => {
        this.snackBar.open(
          `Nota Fiscal #${res.numeroSequencial} criada com sucesso! Status: ${res.status}`,
          'OK',
          { duration: 5000, horizontalPosition: 'end', verticalPosition: 'top', panelClass: ['snack-success'] }
        );
        this.itensFormArray.clear();
        this.adicionarItem();
        this.isLoading = false;
      },
      error: () => {
        this.isLoading = false;
      },
    });
    this.subscriptions.add(sub);
  }
}
