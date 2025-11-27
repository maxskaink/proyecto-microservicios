import { Injectable, Injector } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable, map, filter, take, switchMap } from 'rxjs';
import { Tenant } from '../Models/Tenant';
import { TenantPeticion } from '../Models/TenantPeticion';

@Injectable({ providedIn: 'root' })
export class TenantService {
  private tenantSubject = new BehaviorSubject<string | null>(null);
  tenant$ = this.tenantSubject.asObservable();

  private apiUrlTenant = 'http://localhost:80/api/';

  constructor(
    private http: HttpClient,
    private injector: Injector
  ) {
    const saved = localStorage.getItem('currentTenant');
    if (saved) {
      this.tenantSubject.next(saved);
    }
  }

  setTenant(tenantId: string) {
    this.tenantSubject.next(tenantId);
    localStorage.setItem('currentTenant', tenantId);
  }

  clearTenant() {
    this.tenantSubject.next(null);
    localStorage.removeItem('currentTenant');
  }

  public getTenant(): string | null {
    return this.tenantSubject.getValue();
  }

  /** 🔥 Devuelve el tenant SOLO cuando ya existe */
  getTenantId(): Observable<string> {
    return this.tenant$.pipe(
      filter((t): t is string => !!t),
      take(1)
    );
  }

  getTenants(): Observable<Tenant[]> {
    return this.http.get<Tenant[]>(`${this.apiUrlTenant}tenants`);
  }

  getIdTenants(): Observable<string[]> {
    return this.getTenants().pipe(
      map(tenants => tenants.map(t => t.tenant_id))
    );
  }

  getTenantByID(idTenant: string): Observable<Tenant> {
    return this.http.get<Tenant>(`${this.apiUrlTenant}tenants/${idTenant}`);
  }

  PostTenant(tenantPost: TenantPeticion): Observable<Tenant> {
    return this.http.post<Tenant>(`${this.apiUrlTenant}tenants`, tenantPost);
  }

  updateTenant(tenantId: string, tenantData: Partial<TenantPeticion>): Observable<Tenant> {
    return this.http.put<Tenant>(`${this.apiUrlTenant}tenants/${tenantId}`, tenantData);
  }

  deleteTenant(tenantId: string): Observable<any> {
    return this.http.delete(`${this.apiUrlTenant}tenants/${tenantId}`);
  }

  /** Obtiene los detalles del tenant actual */
  getCurrentUserTenant(): Observable<Tenant> {
    return this.getTenantId().pipe(
      switchMap(id => this.getTenantByID(id))
    );
  }
}
