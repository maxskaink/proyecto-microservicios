import { Injectable } from '@angular/core';
import { HttpClient, HttpContext, HttpContextToken, HttpHeaders } from '@angular/common/http';
import { Observable, from, switchMap, combineLatest, throwError } from 'rxjs';
import { map, filter, take, tap, catchError } from 'rxjs/operators';
import { Product } from '../Models/Product';
import { ProductPeticion } from '../Models/PrdocutPeticion';
import { AuthService } from './Authser.vice';
import { API_BASE } from '../config';
import { TenantService } from './TenantService';

export const SKIP_INTERCEPTOR = new HttpContextToken<boolean>(() => false);

@Injectable({ providedIn: 'root' })
export class ProductService {
  private apiUrlProduct = `${API_BASE}/`;

  constructor(
    private http: HttpClient,
    private authService: AuthService,
    private tenatnService: TenantService,
  ) {}

  /**
   * Combina tenant ID y headers de autenticación
   */
  private getTenant(): Observable<{ tenantId: string }> {
    return this.tenatnService.getTenantId().pipe(
      filter((tenantId): tenantId is string => !!tenantId),
      take(1),
      map((tenantId) => ({ tenantId })),
    );
  }

  /**
   * Solicita una URL pre-firmada para subir una imagen
   */
  getUploadUrl(
    filename: string,
    contentType: string,
  ): Observable<{ upload_url: string; object_key: string }> {
    return this.getTenant().pipe(
      switchMap(({ tenantId }) =>
        this.http.post<{ upload_url: string; object_key: string }>(
          `${this.apiUrlProduct}${tenantId}/api/products/upload-url`,
          { filename, content_type: contentType },
        ),
      ),
    );
  }

  /**
   * Sube la imagen a la URL pre-firmada
   */
  uploadImageToUrl(uploadUrl: string, file: File): Observable<any> {
    const headers = new HttpHeaders({
      'Content-Type': file.type,
    });

    const context = new HttpContext().set(SKIP_INTERCEPTOR, true);

    return this.http.put(uploadUrl, file, { headers, context });
  }

  /**
   * Actualiza la foto del producto con el object_key
   */
  updateProductPhoto(productId: string, objectKey: string): Observable<any> {
    return this.getTenant().pipe(
      switchMap(({ tenantId }) =>
        this.http.put(`${this.apiUrlProduct}${tenantId}/api/products/${productId}/photo`, {
          object_key: objectKey,
        }),
      ),
    );
  }

  /**
   * Crea un nuevo producto - VERSIÓN REACTIVA
   */
  postProduct(product: ProductPeticion): Observable<Product> {
    // Verificar estado actual del tenant antes de continuar

    return this.getTenant().pipe(
      tap(({ tenantId }) => {
        console.log('🌐 Tenant ID obtenido:', tenantId);
      }),
      switchMap(({ tenantId }) => {
        const url = `${this.apiUrlProduct}${tenantId}/api/products`;
        console.log('🌐 URL completa construida:', url);
        return this.http.post<Product>(url, product).pipe(
          tap((response) => console.log('✅ Respuesta del servidor:', response)),
          catchError((error) => {
            console.error('❌ Error en POST request:', error);
            return throwError(error);
          }),
        );
      }),
      catchError((error) => {
        console.error('❌ Error general en postProduct:', error);
        return throwError(error);
      }),
    );
  }

  /**
   * Obtiene productos con paginación - VERSIÓN REACTIVA
   */
  getProducts(page: number, pageSize: number): Observable<Product[]> {
    return this.getTenant().pipe(
      switchMap(({ tenantId }) =>
        this.http.get<Product[]>(
          `${this.apiUrlProduct}${tenantId}/api/products?page=${page}&page_size=${pageSize}`,
        ),
      ),
    );
  }

  /**
   * Obtiene un producto por ID - VERSIÓN REACTIVA
   */
  getProductById(productId: string): Observable<Product> {
    return this.getTenant().pipe(
      switchMap(({ tenantId }) =>
        this.http.get<Product>(`${this.apiUrlProduct}${tenantId}/api/products/${productId}`),
      ),
    );
  }
  /**
   * Obtiene un producto por ID de un tenant específico
   * @param productId ID del producto
   * @param idTenant ID del tenant específico
   */
  getProductTenantById(productId: string, idTenant: string): Observable<Product> {
    return this.http.get<Product>(`${this.apiUrlProduct}${idTenant}/api/products/${productId}`);
  }

  /**
   * Actualiza un producto existente
   */
  updateProduct(productId: string, product: ProductPeticion): Observable<Product> {
    return this.getTenant().pipe(
      switchMap(({ tenantId }) =>
        this.http.put<Product>(
          `${this.apiUrlProduct}${tenantId}/api/products/${productId}`,
          product,
        ),
      ),
    );
  }

  /**
   * Elimina un producto
   */
  deleteProduct(productId: string): Observable<void> {
    return this.getTenant().pipe(
      switchMap(({ tenantId }) =>
        this.http.delete<void>(`${this.apiUrlProduct}${tenantId}/api/products/${productId}`),
      ),
    );
  }
  getCategories(): Observable<string[]> {
    return this.getTenant().pipe(
      switchMap(({ tenantId }) =>
        this.http.get<string[]>(`${this.apiUrlProduct}${tenantId}/api/products/categories`),
      ),
    );
  }
}
