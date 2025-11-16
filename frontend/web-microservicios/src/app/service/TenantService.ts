import { Injectable, Injector } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, from, switchMap, map, combineLatest } from 'rxjs';
import { filter, take } from 'rxjs/operators';
import { Tenant } from '../Models/Tenant';
import { TenantPeticion } from '../Models/TenantPeticion';
import { AuthService } from './Authser.vice';

@Injectable({ providedIn: 'root' })
export class TenantService {

  private apiUrlTenant = 'http://localhost:80/api/'; 

  constructor(
    private http: HttpClient, 
    private authService: AuthService, 
    private injector: Injector
  ) {}

  /**
   * Método privado para obtener headers con autenticación
   */
  private getAuthHeaders(): Observable<HttpHeaders> {
    return from(this.authService.getToken()).pipe(
      map(token => {
        if (!token) {
          throw new Error('No se pudo obtener el token de autenticación');
        }
        return new HttpHeaders({
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        });
      })
    );
  }

  /**
   * Obtiene la lista de tenants desde el backend (sin autenticación requerida)
   * @returns Un observable con un array de tenants
   */
  getTenants(): Observable<Tenant[]> {
    const url = `${this.apiUrlTenant}tenants`;
    return this.http.get<Tenant[]>(url);
  }

  /**
   * Obtiene una lista de IDs de tenants.
   */
  getIdTenants(): Observable<string[]> {
    return this.getTenants().pipe(
      map((tenants: Tenant[]) => {
        return tenants.map(tenant => tenant.tenant_id);
      })
    );
  }

  /**
   * Obtiene un tenant por su ID desde el backend.
   */
  getTenantByID(): Observable<Tenant> {
    return this.getAuthHeaders().pipe(
      switchMap(headers => {
        const url = `${this.apiUrlTenant}tenants/id`;
        return this.http.get<Tenant>(url, { headers });
      })
    );
  }

  /**
   * Crea un nuevo tenant en el backend.
   */
  PostTenant(tenantPost: TenantPeticion): Observable<Tenant> {
    return this.getAuthHeaders().pipe(
      switchMap(headers => {
        const url = `${this.apiUrlTenant}tenants`;
        return this.http.post<Tenant>(url, tenantPost, { headers });
      })
    );
  }

  /**
   * Actualiza un tenant existente
   */
  updateTenant(tenantId: string, tenantData: Partial<TenantPeticion>): Observable<Tenant> {
    return this.getAuthHeaders().pipe(
      switchMap(headers => {
        const url = `${this.apiUrlTenant}tenants/${tenantId}`;
        return this.http.put<Tenant>(url, tenantData, { headers });
      })
    );
  }

  /**
   * Elimina un tenant
   */
  deleteTenant(tenantId: string): Observable<any> {
    return this.getAuthHeaders().pipe(
      switchMap(headers => {
        const url = `${this.apiUrlTenant}tenants/${tenantId}`;
        return this.http.delete(url, { headers });
      })
    );
  }

  /**
   * Obtiene el tenant del usuario actual
   */
  getCurrentUserTenant(): Observable<Tenant | null> {
    return combineLatest([
      this.authService.idTenant$.pipe(
        filter((tenantId): tenantId is string => !!tenantId),
        take(1)
      ),
      this.getAuthHeaders()
    ]).pipe(
      switchMap(([tenantId, headers]) => {
        const url = `${this.apiUrlTenant}tenants/${tenantId}`;
        return this.http.get<Tenant>(url, { headers });
      })
    );
  }

  /**
   * Verifica si un tenant ID ya existe
   */
  checkTenantExists(tenantId: string): Observable<boolean> {
    return this.getIdTenants().pipe(
      map(tenantIds => tenantIds.includes(tenantId))
    );
  }
}