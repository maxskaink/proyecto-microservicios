import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, from, switchMap } from 'rxjs';
import { filter, take } from 'rxjs/operators';
import { Product } from '../Models/Product';
import { Auth, authState } from '@angular/fire/auth';
import { ProductPeticion } from '../Models/PrdocutPeticion';
import { Tenant } from '../Models/Tenant';
import { TenantPeticion } from '../Models/TenantPeticion';

@Injectable({ providedIn: 'root' })
export class TenantService {
  private apiUrl = 'http://localhost:80/api/'; // cambia por tu endpoint

  constructor(private http: HttpClient, private auth: Auth) {}
    
  /**
   * Obtiene la lista de tenants desde el backend.
   * @returns Un observable con un array de tenants
   */
  getTenants(): Observable<Tenant[]> {
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
                return this.http.get<Tenant[]>(this.apiUrl  + 'tenants', { headers });
            })
            );
        })
        );
    }
    /**
     * Obtiene un tenant por su ID desde el backend.
     * @return Un observable con el tenant solicitado
     */
    getTenatnByID(): Observable<Tenant>{
     return authState(this.auth).pipe(
        filter(user => user !== null), 
        take(1),
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
                return this.http.get<Tenant>(this.apiUrl  + 'tenants/id', { headers });
            })
            );
        })
        );
    }
    /**
     * Crea un nuevo tenant en el backend.
     * @param tenantPost Los datos del tenant a crear
     */
    PostTenant(tenantPost : TenantPeticion): Observable<Tenant>{
            return authState(this.auth).pipe(
        filter(user => user !== null), 
        take(1),
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
                return this.http.post<Tenant>(this.apiUrl  + 'tenants', tenantPost, { headers });
            })
            );
        })
        );
    }

}
