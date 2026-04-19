import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { catchError, throwError } from 'rxjs';
export const errorInterceptor: HttpInterceptorFn = (req, next) => {
  const snackBar = inject(MatSnackBar);
  return next(req).pipe(
    catchError(error => {
      const mensagem = error?.error?.erro
        || error?.message
        || 'Erro inesperado. Tente novamente.';
      snackBar.open(`${mensagem}`, 'Fechar', {
        duration: 6000,
        horizontalPosition: 'end',
        verticalPosition: 'top',
        panelClass: ['snack-error'],
      });
      return throwError(() => error);
    })
  );
};
