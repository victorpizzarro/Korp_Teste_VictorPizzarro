import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
@Component({
  selector: 'app-loading-overlay',
  standalone: true,
  imports: [CommonModule, MatProgressSpinnerModule],
  template: `
    @if (isLoading) {
      <div class="overlay">
        <div class="spinner-container">
          <mat-spinner diameter="56" color="primary"></mat-spinner>
          <span class="loading-text">Processando...</span>
        </div>
      </div>
    }
  `,
  styles: [`
    .overlay {
      position: fixed;
      top: 0;
      left: 0;
      width: 100vw;
      height: 100vh;
      background: rgba(0, 0, 0, 0.65);
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 9999;
      backdrop-filter: blur(4px);
    }
    .spinner-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 16px;
    }
    .loading-text {
      color: #e6edf3;
      font-size: 14px;
      font-weight: 500;
      letter-spacing: 0.5px;
    }
  `]
})
export class LoadingOverlayComponent {
  @Input() isLoading = false;
}
