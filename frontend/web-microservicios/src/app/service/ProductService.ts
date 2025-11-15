import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, from, switchMap, combineLatest, throwError } from 'rxjs';
import { map, filter, take, tap, catchError } from 'rxjs/operators';
import { Product } from '../Models/Product';
import { ProductPeticion } from '../Models/PrdocutPeticion';
import { AuthService } from './Authser.vice';

@Injectable({ providedIn: 'root' })
export class ProductService {
  private apiUrlProduct = 'http://localhost:80/';

  constructor(
    private http: HttpClient, 
    private authService: AuthService
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
   * Combina tenant ID y headers de autenticación
   */
  private getTenantAndHeaders(): Observable<{ tenantId: string; headers: HttpHeaders }> {
    return combineLatest([
      this.authService.idTenant$.pipe(
        filter((tenantId): tenantId is string => !!tenantId), // Filtrar valores null/undefined
        take(1) // Tomar solo el primer valor válido
      ),
      this.getAuthHeaders()
    ]).pipe(
      map(([tenantId, headers]) => ({ tenantId, headers }))
    );
  }

  /**
   * Solicita una URL pre-firmada para subir una imagen
   */
  getUploadUrl(filename: string, contentType: string): Observable<{ upload_url: string; object_key: string }> {
    return this.getTenantAndHeaders().pipe(
      switchMap(({ tenantId, headers }) =>
        this.http.post<{ upload_url: string; object_key: string }>(
          `${this.apiUrlProduct}${tenantId}/api/products/upload-url`,
          { filename, content_type: contentType },
          { headers }
        )
      )
    );
  }

  /**
   * Sube la imagen a la URL pre-firmada
   */
  uploadImageToUrl(uploadUrl: string, file: File): Observable<any> {
    const headers = new HttpHeaders({
      'Content-Type': file.type
    });
    return this.http.put(uploadUrl, file, { headers });
  }

  /**
   * Actualiza la foto del producto con el object_key
   */
  updateProductPhoto(productId: string, objectKey: string): Observable<any> {
    return this.getTenantAndHeaders().pipe(
      switchMap(({ tenantId, headers }) =>
        this.http.put(
          `${this.apiUrlProduct}${tenantId}/api/products/${productId}/photo`,
          { object_key: objectKey },
          { headers }
        )
      )
    );
  }

  /**
   * Crea un nuevo producto - VERSIÓN REACTIVA
   */
   postProduct(product: ProductPeticion): Observable<Product> {    
    // Verificar estado actual del tenant antes de continuar
    console.log('🔍 Verificando estado del tenant...');
    this.authService.idTenant$.pipe(take(1)).subscribe(
      tenantId => console.log('🏢 Estado actual del tenant en AuthService:', tenantId),
      error => console.error('❌ Error al obtener tenant:', error)
    );
    
    return this.getTenantAndHeaders().pipe(
      tap(({ tenantId, headers }) => {
        console.log('🌐 Tenant ID obtenido:', tenantId);
        console.log('🔐 Headers construidos:', headers.keys());
      }),
      switchMap(({ tenantId, headers }) => {
        const url = `${this.apiUrlProduct}${tenantId}/api/products`;
        console.log('🌐 URL completa construida:', url);
        return this.http.post<Product>(url, product, { headers }).pipe(
          tap(response => console.log('✅ Respuesta del servidor:', response)),
          catchError(error => {
            console.error('❌ Error en POST request:', error);
            return throwError(error);
          })
        );
      }),
      catchError(error => {
        console.error('❌ Error general en postProduct:', error);
        return throwError(error);
      })
    );
  }

  /**
   * Obtiene productos con paginación - VERSIÓN REACTIVA
   */
  getProducts(page: number, pageSize: number): Observable<Product[]> {
    return this.getTenantAndHeaders().pipe(
      switchMap(({ tenantId, headers }) => 
        this.http.get<Product[]>(
          `${this.apiUrlProduct}${tenantId}/api/products?page=${page}&pag_size=${pageSize}`,
          { headers }
        )
      )
    );
  }

  /**
   * Obtiene un producto por ID - VERSIÓN REACTIVA
   */
  getProductById(productId: string): Observable<Product> {
    return this.getTenantAndHeaders().pipe(
      switchMap(({ tenantId, headers }) =>
        this.http.get<Product>(
          `${this.apiUrlProduct}${tenantId}/api/products/${productId}`,
          { headers }
        )
      )
    );
  }
  /**
   * Obtiene un producto por ID de un tenant específico
   * @param productId ID del producto
   * @param idTenant ID del tenant específico
   */
  getProductTenantById(productId: string, idTenant: string): Observable<Product> {
    return this.getAuthHeaders().pipe(  
      switchMap((headers) =>
        this.http.get<Product>(
          `${this.apiUrlProduct}${idTenant}/api/products/${productId}`,
          { headers }
        )
      )
    );
  }
  /**
   * Actualiza un producto existente
   */
  updateProduct(productId: string, product: ProductPeticion): Observable<Product> {
    return this.getTenantAndHeaders().pipe(
      switchMap(({ tenantId, headers }) =>
        this.http.put<Product>(
          `${this.apiUrlProduct}${tenantId}/api/products/${productId}`,
          product,
          { headers }
        )
      )
    );
  }

  /**
   * Elimina un producto
   */
  deleteProduct(productId: string): Observable<void> {
    return this.getTenantAndHeaders().pipe(
      switchMap(({ tenantId, headers }) =>
        this.http.delete<void>(
          `${this.apiUrlProduct}${tenantId}/api/products/${productId}`,
          { headers }
        )
      )
    );
  }
}