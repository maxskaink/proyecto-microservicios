import { Injectable } from '@angular/core';
import { CanActivate, ActivatedRouteSnapshot, Router, UrlTree } from '@angular/router';
import { AuthService } from '../service/Authser.vice';
import { UserResponseBack } from '../Models/UserReponseBack';


@Injectable({
  providedIn: 'root'
})
export class RoleGuard implements CanActivate {

  constructor(
    private authService: AuthService,
    private router: Router
  ) {}

  canActivate(route: ActivatedRouteSnapshot): boolean | UrlTree {
  const allowedRoles: string[] = route.data['roles'];  // roles permitidos

  const stored = localStorage.getItem('user_data');

  if (!stored) {
    console.warn(' No hay user_data almacenado');
    return this.router.parseUrl('/login');
  }
  let user: UserResponseBack | null = null;

  try {
    user = JSON.parse(stored);
  } catch (err) {
    console.error('⚠ Error parseando user_data:', err);
    return this.router.parseUrl('/login');
  }

  // 3️⃣ Extraer el rol tal como viene del backend
  const userRole = user?.rol;

  if (!userRole) {
    console.warn(' No se encontró rol en user_data');
    return this.router.parseUrl('/login');
  }

  if (allowedRoles.includes(userRole)) {
    return true;
  }

  console.warn(`Acceso denegado. Requeridos: ${allowedRoles}, pero el usuario tiene: ${userRole}`);
  return this.router.parseUrl('/forbidden');
}

}
