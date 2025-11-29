import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { AuthService } from './Authser.vice';
import { TenantService } from './TenantService';
import { Observable, switchMap, take } from 'rxjs';
import { ShippingResponse } from '../Models/ShippingResponse';


@Injectable({
  providedIn: 'root',
})
export class ShippingService {
  private apiUrlShippingCart = 'http://localhost:80/';

  constructor(
    private http: HttpClient,
    private authService: AuthService,
    private tenantService: TenantService
  ) {}

  private getTenant(): Observable<string> {
    return this.tenantService.getTenantId().pipe(take(1));
  }

  /**
   * 
   * @returns lista de envíos
   */
  getShippingList(): Observable<ShippingResponse[]> {
    return this.getTenant().pipe(
      switchMap((tenantId) => {
        const url = `${this.apiUrlShippingCart}${tenantId}/api/shippings`;
        return this.http.get<ShippingResponse[]>(url);
      })
    );
  }
  /**
   * funcion para obtener un envío por su id
   * @param id id del envío
   * @returns información del envío
   */
  getShippingById(id: string): Observable<ShippingResponse> {
    return this.getTenant().pipe(
      switchMap((tenantId) => {
        const url = `${this.apiUrlShippingCart}${tenantId}/api/shippings/${id}`;
        return this.http.get<ShippingResponse>(url);
      })
    );
  }
  /**
   * Funcion para obtener el envío asociado a una orden
   * @param idOrder id de la orden
   * @returns información del envío asociado a la orden
   */
  getShippingByOrderId(idOrder: string): Observable<ShippingResponse> {
    return this.getTenant().pipe(
      switchMap((tenantId) => {
        const url = `${this.apiUrlShippingCart}${tenantId}/api/shippings/order/${idOrder}`;
        return this.http.get<ShippingResponse>(url);
      })
    );
  }

  updateStatusShipping(id: string, status: string): Observable<boolean> {
    return this.getTenant().pipe(
      switchMap((tenantId) => {
        const url = `${this.apiUrlShippingCart}${tenantId}/api/shippings/${id}/status`;
        return this.http.put<boolean>(url, { status });
      })
    );
  }
}
