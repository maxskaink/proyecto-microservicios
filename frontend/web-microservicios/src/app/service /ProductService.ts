import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, from, switchMap } from 'rxjs';
import { filter, take } from 'rxjs/operators';
import { Product } from '../Models/Product';
import { Auth, authState } from '@angular/fire/auth';

@Injectable({ providedIn: 'root' })
export class ProductService {
  private apiUrl = 'http://localhost:80/api/products?page=1&pag_size=10'; // cambia por tu endpoint

  constructor(private http: HttpClient, private auth: Auth) {}

  getProducts(): Observable<Product[]> {
    return authState(this.auth).pipe(
      filter(user => user !== null), 
      take(1), // Tomar solo el primer valor válido
      switchMap(user => {
        console.log('Usuario autenticado:', user.email);
        return from(user.getIdToken()).pipe(
          switchMap(idToken => {        
            if (!idToken) {
              throw new Error('No se pudo obtener el token de autenticación');
            }
            const headers = new HttpHeaders({
              'Authorization': `Bearer ${idToken}`,
              'Content-Type': 'application/json'
            });
            return this.http.get<Product[]>(this.apiUrl, { headers });
          })
        );
      })
    );
  }
}
