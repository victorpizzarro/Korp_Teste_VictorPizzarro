import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, catchError, tap, throwError } from 'rxjs';
export interface ItemNotaFiscal {
  codigoProduto: string;
  quantidade: number;
}
export interface NotaFiscal {
  numeroSequencial: number;
  status: string;         
  dataCriacao: string;    
  itens: ItemNotaFiscal[];
  mensagem?: string;
}
export interface CriarNotaFiscalPayload {
  itens: ItemNotaFiscal[];
}
@Injectable({
  providedIn: 'root' 
})
export class NotaFiscalService {
  private readonly API_URL = 'http://localhost:8082/api/v1/notas-fiscais';
  constructor(private http: HttpClient) {}
  listar(): Observable<NotaFiscal[]> {
    return this.http.get<NotaFiscal[]>(this.API_URL).pipe(
      tap(notas => console.log(`[NotaFiscalService] ${notas.length} notas carregadas`)),
      catchError(err => {
        console.error('[NotaFiscalService] Erro ao listar notas:', err);
        return throwError(() => err);
      })
    );
  }
  cadastrar(payload: CriarNotaFiscalPayload): Observable<NotaFiscal> {
    return this.http.post<NotaFiscal>(this.API_URL, payload).pipe(
      tap(res => console.log('[NotaFiscalService] NF criada:', res.numeroSequencial)),
      catchError(err => {
        console.error('[NotaFiscalService] Erro ao cadastrar NF:', err);
        return throwError(() => err);
      })
    );
  }
  imprimir(numero: number): Observable<NotaFiscal> {
    return this.http.post<NotaFiscal>(`${this.API_URL}/${numero}/imprimir`, {}).pipe(
      tap(res => console.log(`[NotaFiscalService] NF ${numero} impressa. Status: ${res.status}`)),
      catchError(err => {
        console.error(`[NotaFiscalService] Erro ao imprimir NF ${numero}:`, err);
        return throwError(() => err);
      })
    );
  }
}
