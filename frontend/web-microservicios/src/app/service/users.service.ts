import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, switchMap, take } from 'rxjs';
import { AuthService } from './Authser.vice';
import { API_BASE } from '../config';
import { TenantService } from './TenantService';
import { UserResponseBack } from '../Models/UserReponseBack';
import { userPeticion } from '../Models/UserPeticion';

@Injectable({
  providedIn: 'root',
})
export class UsersService {
  private apiUrlShippingCart = `${API_BASE}/`;

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
   * @param userId id del usuario a buscar
   * @returns deuelve un observable con la respuesta del usuario
   */
  public getUserByID(userId: string): Observable<UserResponseBack>{
    return this.getTenant().pipe(
      take(1),
      switchMap((tenantId: string) => {
        const url = `${this.apiUrlShippingCart}${tenantId}/api/users/producer/${userId}`;
        return this.http.get<UserResponseBack>(url);
      })
    );
  }
  /**
   * 
   * @returns devuelve un observable con la respuesta de todos los usuarios
   */
  public getAllUsers(): Observable<UserResponseBack[]> {
    return this.getTenant().pipe(
      take(1),
      switchMap((tenantId: string) => {
        const url = `${this.apiUrlShippingCart}${tenantId}/api/users`;
        return this.http.get<UserResponseBack[]>(url);
      })
    );
  }
  /** Actualiza el rol de un usuario específico.
  * @param userId El ID del usuario cuyo rol se va a actualizar.
  * @param newRole El nuevo rol que se asignará al usuario.
  * @returns Un Observable que emite la respuesta del servidor después de la actualización.
  */
  public updateRoleUser(userId: string, newRole: string): Observable<UserResponseBack> {
    return this.getTenant().pipe(
      take(1),
      switchMap((tenantId: string) => {
        const url = `${this.apiUrlShippingCart}${tenantId}/api/users/${userId}/rol`;
        return this.http.patch<UserResponseBack>(url, { rol: newRole });
      })
    );
  }

  /**
   * edita la informacion de un usuario
   * @param userId id del usuario a editar
   * @param userData datos del usuario a editar
   * @returns devuelve un observable con la respuesta del usuario editado
   */
  public updateUser(userId: string, userData: userPeticion): Observable<UserResponseBack> {
    return this.getTenant().pipe(
      take(1),
      switchMap((tenantId: string) => {
        const url = `${this.apiUrlShippingCart}${tenantId}/api/users/${userId}`;
        return this.http.put<UserResponseBack>(url, userData);
      })
    );
  }
}
