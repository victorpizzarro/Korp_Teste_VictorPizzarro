import { Component } from '@angular/core';
import { RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    RouterLink,
    RouterLinkActive,
    MatToolbarModule,
    MatSidenavModule,
    MatListModule,
    MatIconModule,
    MatButtonModule,
  ],
  template: `
    <mat-sidenav-container class="app-container">
      <!-- Painel de Navegação Lateral -->
      <mat-sidenav #sidenav mode="side" opened class="sidenav">
        <div class="sidenav-header">
          <div class="logo">
            <span class="logo-text">Teste Korp</span>
          </div>
          <span class="logo-subtitle">Sistema de Faturamento</span>
        </div>
        <mat-nav-list>
          <div class="nav-section-label">Produtos</div>
          <a mat-list-item routerLink="/produtos/cadastro" routerLinkActive="active-link" id="nav-cadastro-produto">
            <mat-icon matListItemIcon>add_box</mat-icon>
            <span matListItemTitle>Cadastrar Produto</span>
          </a>
          <a mat-list-item routerLink="/produtos/lista" routerLinkActive="active-link" id="nav-lista-produtos">
            <mat-icon matListItemIcon>inventory_2</mat-icon>
            <span matListItemTitle>Saldo de Estoque</span>
          </a>
          <div class="nav-section-label">Notas Fiscais</div>
          <a mat-list-item routerLink="/notas-fiscais/cadastro" routerLinkActive="active-link" id="nav-cadastro-nf">
            <mat-icon matListItemIcon>note_add</mat-icon>
            <span matListItemTitle>Emitir Nota Fiscal</span>
          </a>
          <a mat-list-item routerLink="/notas-fiscais/lista" routerLinkActive="active-link" id="nav-lista-nf">
            <mat-icon matListItemIcon>list_alt</mat-icon>
            <span matListItemTitle>Listar / Imprimir</span>
          </a>
        </mat-nav-list>
        <div class="sidenav-footer">
          <span>Victor Pizzarro © 2026</span>
        </div>
      </mat-sidenav>
      <!-- Conteúdo Principal -->
      <mat-sidenav-content class="main-content">
        <mat-toolbar class="top-toolbar">
          <button mat-icon-button (click)="sidenav.toggle()" id="btn-toggle-sidenav">
            <mat-icon>menu</mat-icon>
          </button>
          <span class="toolbar-title">Sistema de Emissão de Notas Fiscais</span>
        </mat-toolbar>
        <main class="content-area">
          <router-outlet />
        </main>
      </mat-sidenav-content>
    </mat-sidenav-container>
  `,
  styles: [`
    .app-container {
      height: 100vh;
    }
    .sidenav {
      width: 260px;
      background: var(--korp-bg-card);
      border-right: 1px solid var(--korp-border);
    }
    .sidenav-header {
      padding: 24px 20px 16px;
      border-bottom: 1px solid var(--korp-border);
    }
    .logo {
      display: flex;
      align-items: center;
      gap: 10px;
    }
    .logo-icon {
      font-size: 32px;
      width: 32px;
      height: 32px;
      color: #7c8aff;
    }
    .logo-text {
      font-size: 24px;
      font-weight: 700;
      background: linear-gradient(135deg, #7c8aff, #4f6bed);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }
    .logo-subtitle {
      display: block;
      font-size: 12px;
      color: var(--korp-text-secondary);
      margin-top: 4px;
    }
    .nav-section-label {
      padding: 20px 16px 8px;
      font-size: 11px;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 1.2px;
      color: var(--korp-text-secondary);
    }
    .active-link {
      background: rgba(79, 107, 237, 0.12) !important;
      border-right: 3px solid #4f6bed;
      mat-icon {
        color: #7c8aff !important;
      }
    }
    .sidenav-footer {
      position: absolute;
      bottom: 0;
      width: 100%;
      padding: 16px;
      text-align: center;
      font-size: 11px;
      color: var(--korp-text-secondary);
      border-top: 1px solid var(--korp-border);
    }
    .top-toolbar {
      background: var(--korp-bg-card) !important;
      border-bottom: 1px solid var(--korp-border);
      color: var(--korp-text-primary) !important;
    }
    .toolbar-title {
      margin-left: 12px;
      font-size: 16px;
      font-weight: 500;
    }
    .main-content {
      background: var(--korp-bg-primary);
    }
    .content-area {
      padding: 0;
      min-height: calc(100vh - 64px);
    }
  `],
})
export class AppComponent {
  title = 'korp-frontend';
}
