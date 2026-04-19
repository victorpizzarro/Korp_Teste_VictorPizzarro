import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, catchError, tap, throwError } from 'rxjs';
export interface Produto {
  codigo: string;
  descricao: string;
  saldo: number;
  mensagem?: string;
}
export interface CriarProdutoPayload {
  codigo: string;
  descricao: string;
  saldo: number;
}
@Injectable({
  providedIn: 'root' 
})
export class ProdutoService {
  private readonly API_URL = 'http://localhost:8081/api/v1/produtos';
  constructor(private http: HttpClient) {}
  listar(): Observable<Produto[]> {
    return this.http.get<Produto[]>(this.API_URL).pipe(
      tap(produtos => console.log(`[ProdutoService] ${produtos.length} produtos carregados`)),
      catchError(err => {
        console.error('[ProdutoService] Erro ao listar produtos:', err);
        return throwError(() => err);
      })
    );
  }
  cadastrar(payload: CriarProdutoPayload): Observable<Produto> {
    return this.http.post<Produto>(this.API_URL, payload).pipe(
      tap(res => console.log('[ProdutoService] Produto cadastrado:', res.codigo)),
      catchError(err => {
        console.error('[ProdutoService] Erro ao cadastrar produto:', err);
        return throwError(() => err);
      })
    );
  }
}
