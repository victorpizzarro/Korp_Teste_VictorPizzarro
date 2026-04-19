import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormGroup, FormControl, Validators } from '@angular/forms';
import { Subscription } from 'rxjs';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { ProdutoService } from '../../../core/services/produto.service';
@Component({
  selector: 'app-cadastro-produto',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatSnackBarModule,
    MatProgressSpinnerModule,
  ],
  template: `
    <div class="page-container">
      <h1 class="page-title">Cadastro de Produtos</h1>
      <p class="page-subtitle">Cadastre novos produtos com código, descrição e saldo inicial de estoque.</p>
      <div class="card">
        <form [formGroup]="produtoForm" (ngSubmit)="onSubmit()" id="form-cadastro-produto">
          <div class="form-grid">
            <mat-form-field appearance="outline">
              <mat-label>Código do Produto</mat-label>
              <input matInput formControlName="codigo" placeholder="Ex: PROD-001" id="input-codigo">
              <mat-icon matPrefix>qr_code</mat-icon>
              @if (produtoForm.get('codigo')?.hasError('required') && produtoForm.get('codigo')?.touched) {
                <mat-error>O código é obrigatório</mat-error>
              }
            </mat-form-field>
            <mat-form-field appearance="outline">
              <mat-label>Descrição</mat-label>
              <input matInput formControlName="descricao" placeholder="Ex: Monitor Gamer 24p" id="input-descricao">
              <mat-icon matPrefix>description</mat-icon>
              @if (produtoForm.get('descricao')?.hasError('required') && produtoForm.get('descricao')?.touched) {
                <mat-error>A descrição é obrigatória</mat-error>
              }
            </mat-form-field>
            <mat-form-field appearance="outline">
              <mat-label>Saldo Inicial</mat-label>
              <input matInput type="number" formControlName="saldo" placeholder="0" id="input-saldo">
              <mat-icon matPrefix>inventory</mat-icon>
              @if (produtoForm.get('saldo')?.hasError('required') && produtoForm.get('saldo')?.touched) {
                <mat-error>O saldo é obrigatório</mat-error>
              }
              @if (produtoForm.get('saldo')?.hasError('min') && produtoForm.get('saldo')?.touched) {
                <mat-error>O saldo não pode ser negativo</mat-error>
              }
            </mat-form-field>
          </div>
          <div class="form-actions">
            <button
              mat-raised-button
              color="primary"
              type="submit"
              [disabled]="produtoForm.invalid || isLoading"
              id="btn-cadastrar-produto">
              <span class="button-content">
                @if (isLoading) {
                  <mat-spinner diameter="20" class="btn-spinner"></mat-spinner>
                } @else {
                  <mat-icon>save</mat-icon>
                }
                <span>Cadastrar Produto</span>
              </span>
            </button>
          </div>
        </form>
      </div>
    </div>
  `,
  styles: [`
    .form-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 16px;
      mat-form-field:last-child {
        grid-column: 1 / -1;
        max-width: 300px;
      }
    }
    .form-actions {
      margin-top: 24px;
      display: flex;
      justify-content: flex-end;
      button {
        padding: 0 24px;
        height: 44px;
        font-weight: 600;
        letter-spacing: 0.3px;
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
export class CadastroProdutoComponent implements OnInit, OnDestroy {
  produtoForm!: FormGroup;
  isLoading = false;
  private subscription?: Subscription;
  constructor(
    private produtoService: ProdutoService,
    private snackBar: MatSnackBar
  ) {}
  ngOnInit(): void {
    this.produtoForm = new FormGroup({
      codigo:    new FormControl('', [Validators.required]),
      descricao: new FormControl('', [Validators.required]),
      saldo:     new FormControl(0,  [Validators.required, Validators.min(0)]),
    });
  }
  ngOnDestroy(): void {
    this.subscription?.unsubscribe();
  }
  onSubmit(): void {
    if (this.produtoForm.invalid) return;
    this.isLoading = true;
    this.subscription = this.produtoService
      .cadastrar(this.produtoForm.value)
      .subscribe({
        next: (res) => {
          this.snackBar.open(`${res.mensagem || 'Produto cadastrado com sucesso!'}`, 'OK', {
            duration: 4000,
            horizontalPosition: 'end',
            verticalPosition: 'top',
            panelClass: ['snack-success'],
          });
          this.produtoForm.reset({ saldo: 0 });
          this.isLoading = false;
        },
        error: () => {
          this.isLoading = false;
        },
      });
  }
}
